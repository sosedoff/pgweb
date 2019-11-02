package sshkeys

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"math/big"
	mrand "math/rand"

	"github.com/dchest/bcrypt_pbkdf"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ssh"
)

// Format of private key to use when Marshaling.
type Format int

const (
	// FormatOpenSSHv1 encodes a private key using OpenSSH's PROTOCOL.key format: https://github.com/openssh/openssh-portable/blob/master/PROTOCOL.key
	FormatOpenSSHv1 Format = iota
	// FormatClassicPEM encodes private keys in PEM, with a key-specific encoding, as used by OpenSSH.
	FormatClassicPEM
)

// MarshalOptions provides the Marshal function format and encryption options.
type MarshalOptions struct {
	// Passphrase to encrypt private key with, if nil, the key will not be encrypted.
	Passphrase []byte
	// Format to encode the private key in.
	Format Format
}

// Marshal converts a private key into an optionally encrypted format.
func Marshal(pk interface{}, opts *MarshalOptions) ([]byte, error) {
	switch opts.Format {
	case FormatOpenSSHv1:
		return marshalOpenssh(pk, opts)
	case FormatClassicPEM:
		return marshalPem(pk, opts)
	default:
		return nil, fmt.Errorf("sshkeys: invalid format %d", opts.Format)
	}
}

func marshalPem(pk interface{}, opts *MarshalOptions) ([]byte, error) {
	var err error
	var plain []byte
	var pemType string

	switch key := pk.(type) {
	case *rsa.PrivateKey:
		pemType = "RSA PRIVATE KEY"
		plain = x509.MarshalPKCS1PrivateKey(key)
	case *ecdsa.PrivateKey:
		pemType = "EC PRIVATE KEY"
		plain, err = x509.MarshalECPrivateKey(key)
		if err != nil {
			return nil, err
		}
	case *dsa.PrivateKey:
		pemType = "DSA PRIVATE KEY"
		plain, err = marshalDSAPrivateKey(key)
		if err != nil {
			return nil, err
		}
	case *ed25519.PrivateKey:
		return nil, fmt.Errorf("sshkeys: ed25519 keys must be marshaled with FormatOpenSSHv1")
	default:
		return nil, fmt.Errorf("sshkeys: unsupported key type %T", pk)
	}

	if len(opts.Passphrase) > 0 {
		block, err := x509.EncryptPEMBlock(rand.Reader, pemType, plain, opts.Passphrase, x509.PEMCipherAES128)
		if err != nil {
			return nil, err
		}
		return pem.EncodeToMemory(block), nil
	}

	return pem.EncodeToMemory(&pem.Block{
		Type:  pemType,
		Bytes: plain,
	}), nil
}

type dsaOpenssl struct {
	Version int
	P       *big.Int
	Q       *big.Int
	G       *big.Int
	Pub     *big.Int
	Priv    *big.Int
}

// https://github.com/golang/crypto/blob/master/ssh/keys.go#L793-L804
func marshalDSAPrivateKey(pk *dsa.PrivateKey) ([]byte, error) {
	k := dsaOpenssl{
		Version: 0,
		P:       pk.P,
		Q:       pk.Q,
		G:       pk.G,
		Pub:     pk.Y,
		Priv:    pk.X,
	}

	return asn1.Marshal(k)
}

const opensshv1Magic = "openssh-key-v1"

type opensshHeader struct {
	CipherName   string
	KdfName      string
	KdfOpts      string
	NumKeys      uint32
	PubKey       string
	PrivKeyBlock string
}

type opensshKey struct {
	Check1  uint32
	Check2  uint32
	Keytype string
	Rest    []byte `ssh:"rest"`
}

type opensshRsa struct {
	N       *big.Int
	E       *big.Int
	D       *big.Int
	Iqmp    *big.Int
	P       *big.Int
	Q       *big.Int
	Comment string
	Pad     []byte `ssh:"rest"`
}

type opensshED25519 struct {
	Pub     []byte
	Priv    []byte
	Comment string
	Pad     []byte `ssh:"rest"`
}

func padBytes(data []byte, blocksize int) []byte {
	if blocksize != 0 {
		var i byte
		for i = byte(1); len(data)%blocksize != 0; i++ {
			data = append(data, i&0xFF)
		}
	}
	return data
}

func marshalOpenssh(pk interface{}, opts *MarshalOptions) ([]byte, error) {
	var blocksize int
	var keylen int

	out := opensshHeader{
		CipherName: "none",
		KdfName:    "none",
		KdfOpts:    "",
		NumKeys:    1,
		PubKey:     "",
	}

	if len(opts.Passphrase) > 0 {
		out.CipherName = "aes256-cbc"
		out.KdfName = "bcrypt"
		keylen = keySizeAES256
		blocksize = aes.BlockSize
	}

	check := mrand.Uint32()
	pk1 := opensshKey{
		Check1: check,
		Check2: check,
	}

	switch key := pk.(type) {
	case *rsa.PrivateKey:
		k := &opensshRsa{
			N:       key.N,
			E:       big.NewInt(int64(key.E)),
			D:       key.D,
			Iqmp:    key.Precomputed.Qinv,
			P:       key.Primes[0],
			Q:       key.Primes[1],
			Comment: "",
		}

		data := ssh.Marshal(k)
		pk1.Keytype = ssh.KeyAlgoRSA
		pk1.Rest = data
		publicKey, err := ssh.NewPublicKey(&key.PublicKey)
		if err != nil {
			return nil, err
		}
		out.PubKey = string(publicKey.Marshal())

	case ed25519.PrivateKey:
		k := opensshED25519{
			Pub:  key.Public().(ed25519.PublicKey),
			Priv: key,
		}
		data := ssh.Marshal(k)
		pk1.Keytype = ssh.KeyAlgoED25519
		pk1.Rest = data

		publicKey, err := ssh.NewPublicKey(key.Public())
		if err != nil {
			return nil, err
		}
		out.PubKey = string(publicKey.Marshal())
	default:
		return nil, fmt.Errorf("sshkeys: unsupported key type %T", pk)
	}

	if len(opts.Passphrase) > 0 {
		rounds := 16
		ivlen := blocksize
		salt := make([]byte, blocksize)
		_, err := rand.Read(salt)
		if err != nil {
			return nil, err
		}

		kdfdata, err := bcrypt_pbkdf.Key(opts.Passphrase, salt, rounds, keylen+ivlen)
		if err != nil {
			return nil, err
		}
		iv := kdfdata[keylen : ivlen+keylen]
		aeskey := kdfdata[0:keylen]

		block, err := aes.NewCipher(aeskey)
		if err != nil {
			return nil, err
		}

		pkblock := padBytes(ssh.Marshal(pk1), blocksize)

		cbc := cipher.NewCBCEncrypter(block, iv)
		cbc.CryptBlocks(pkblock, pkblock)

		out.PrivKeyBlock = string(pkblock)

		var opts struct {
			Salt   []byte
			Rounds uint32
		}

		opts.Salt = salt
		opts.Rounds = uint32(rounds)

		out.KdfOpts = string(ssh.Marshal(&opts))
	} else {
		out.PrivKeyBlock = string(ssh.Marshal(pk1))
	}

	outBytes := []byte(opensshv1Magic)
	outBytes = append(outBytes, 0)
	outBytes = append(outBytes, ssh.Marshal(out)...)
	block := &pem.Block{
		Type:  "OPENSSH PRIVATE KEY",
		Bytes: outBytes,
	}
	return pem.EncodeToMemory(block), nil
}
