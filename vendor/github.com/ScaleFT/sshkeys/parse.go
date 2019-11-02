// Portions of this file are based on https://github.com/golang/crypto/blob/master/ssh/keys.go
//
// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sshkeys

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/dchest/bcrypt_pbkdf"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ssh"
)

// ErrIncorrectPassword is returned when the supplied passphrase was not correct for an encrypted private key.
var ErrIncorrectPassword = errors.New("sshkeys: Invalid Passphrase")

const keySizeAES256 = 32

// ParseEncryptedPrivateKey returns a Signer from an encrypted private key. It supports
// the same keys as ParseEncryptedRawPrivateKey.
func ParseEncryptedPrivateKey(data []byte, passphrase []byte) (ssh.Signer, error) {
	key, err := ParseEncryptedRawPrivateKey(data, passphrase)
	if err != nil {
		return nil, err
	}

	return ssh.NewSignerFromKey(key)
}

// ParseEncryptedRawPrivateKey returns a private key from an encrypted private key. It
// supports RSA (PKCS#1 or OpenSSH), DSA (OpenSSL), and ECDSA private keys.
//
// ErrIncorrectPassword will be returned if the supplied passphrase is wrong,
// but some formats like RSA in PKCS#1 detecting a wrong passphrase is difficult,
// and other parse errors may be returned.
func ParseEncryptedRawPrivateKey(data []byte, passphrase []byte) (interface{}, error) {
	var err error

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("no PEM block found")
	}

	if x509.IsEncryptedPEMBlock(block) {
		data, err = x509.DecryptPEMBlock(block, passphrase)
		if err == x509.IncorrectPasswordError {
			return nil, ErrIncorrectPassword
		}
		if err != nil {
			return nil, err
		}
	} else {
		data = block.Bytes
	}

	switch block.Type {
	case "RSA PRIVATE KEY":
		pk, err := x509.ParsePKCS1PrivateKey(data)
		if err != nil {
			// The Algos for PEM Encryption do not include strong message authentication,
			// so sometimes DecryptPEMBlock works, but ParsePKCS1PrivateKey fails with an asn1 error.
			// We are just catching the most common prefix here...
			if strings.HasPrefix(err.Error(), "asn1: structure error") {
				return nil, ErrIncorrectPassword
			}
			return nil, err
		}
		return pk, nil
	case "EC PRIVATE KEY":
		return x509.ParseECPrivateKey(data)
	case "DSA PRIVATE KEY":
		return ssh.ParseDSAPrivateKey(data)
	case "OPENSSH PRIVATE KEY":
		return parseOpenSSHPrivateKey(data, passphrase)
	default:
		return nil, fmt.Errorf("sshkeys: unsupported key type %q", block.Type)
	}
}

func parseOpenSSHPrivateKey(data []byte, passphrase []byte) (interface{}, error) {
	magic := append([]byte(opensshv1Magic), 0)
	if !bytes.Equal(magic, data[0:len(magic)]) {
		return nil, errors.New("sshkeys: invalid openssh private key format")
	}
	remaining := data[len(magic):]

	w := opensshHeader{}

	if err := ssh.Unmarshal(remaining, &w); err != nil {
		return nil, err
	}

	if w.NumKeys != 1 {
		return nil, fmt.Errorf("sshkeys: NumKeys must be 1: %d", w.NumKeys)
	}

	var privateKeyBytes []byte
	var encrypted bool

	switch {
	// OpenSSH supports bcrypt KDF w/ AES256-CBC or AES256-CTR mode
	case w.KdfName == "bcrypt" && w.CipherName == "aes256-cbc":
		iv, block, err := extractBcryptIvBlock(passphrase, w)
		if err != nil {
			return nil, err
		}

		cbc := cipher.NewCBCDecrypter(block, iv)
		privateKeyBytes = []byte(w.PrivKeyBlock)
		cbc.CryptBlocks(privateKeyBytes, privateKeyBytes)

		encrypted = true

	case w.KdfName == "bcrypt" && w.CipherName == "aes256-ctr":
		iv, block, err := extractBcryptIvBlock(passphrase, w)
		if err != nil {
			return nil, err
		}

		stream := cipher.NewCTR(block, iv)
		privateKeyBytes = []byte(w.PrivKeyBlock)
		stream.XORKeyStream(privateKeyBytes, privateKeyBytes)

		encrypted = true

	case w.KdfName == "none" && w.CipherName == "none":
		privateKeyBytes = []byte(w.PrivKeyBlock)

	default:
		return nil, fmt.Errorf("sshkeys: unknown Cipher/KDF: %s:%s", w.CipherName, w.KdfName)
	}

	pk1 := opensshKey{}

	if err := ssh.Unmarshal(privateKeyBytes, &pk1); err != nil {
		if encrypted {
			return nil, ErrIncorrectPassword
		}
		return nil, err
	}

	if pk1.Check1 != pk1.Check2 {
		return nil, ErrIncorrectPassword
	}

	// we only handle ed25519 and rsa keys currently
	switch pk1.Keytype {
	case ssh.KeyAlgoRSA:
		// https://github.com/openssh/openssh-portable/blob/V_7_4_P1/sshkey.c#L2760-L2773
		key := opensshRsa{}

		err := ssh.Unmarshal(pk1.Rest, &key)
		if err != nil {
			return nil, err
		}

		for i, b := range key.Pad {
			if int(b) != i+1 {
				return nil, errors.New("sshkeys: padding not as expected")
			}
		}

		pk := &rsa.PrivateKey{
			PublicKey: rsa.PublicKey{
				N: key.N,
				E: int(key.E.Int64()),
			},
			D:      key.D,
			Primes: []*big.Int{key.P, key.Q},
		}

		err = pk.Validate()
		if err != nil {
			return nil, err
		}

		pk.Precompute()

		return pk, nil
	case ssh.KeyAlgoED25519:
		key := opensshED25519{}

		err := ssh.Unmarshal(pk1.Rest, &key)
		if err != nil {
			return nil, err
		}

		if len(key.Priv) != ed25519.PrivateKeySize {
			return nil, errors.New("sshkeys: private key unexpected length")
		}

		for i, b := range key.Pad {
			if int(b) != i+1 {
				return nil, errors.New("sshkeys: padding not as expected")
			}
		}

		pk := ed25519.PrivateKey(make([]byte, ed25519.PrivateKeySize))
		copy(pk, key.Priv)
		return pk, nil
	default:
		return nil, errors.New("sshkeys: unhandled key type")
	}
}

func extractBcryptIvBlock(passphrase []byte, w opensshHeader) ([]byte, cipher.Block, error) {
	cipherKeylen := keySizeAES256
	cipherIvLen := aes.BlockSize

	var opts struct {
		Salt   []byte
		Rounds uint32
	}

	if err := ssh.Unmarshal([]byte(w.KdfOpts), &opts); err != nil {
		return nil, nil, err
	}
	kdfdata, err := bcrypt_pbkdf.Key(passphrase, opts.Salt, int(opts.Rounds), cipherKeylen+cipherIvLen)
	if err != nil {
		return nil, nil, err
	}

	iv := kdfdata[cipherKeylen : cipherIvLen+cipherKeylen]
	aeskey := kdfdata[0:cipherKeylen]
	block, err := aes.NewCipher(aeskey)

	if err != nil {
		return nil, nil, err
	}

	return iv, block, nil
}
