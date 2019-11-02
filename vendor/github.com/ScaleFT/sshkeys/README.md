# sshkeys

[![GoDoc](https://godoc.org/github.com/ScaleFT/sshkeys?status.svg)](https://godoc.org/github.com/ScaleFT/sshkeys)
[![Build Status](https://travis-ci.org/ScaleFT/sshkeys.svg?branch=master)](https://travis-ci.org/ScaleFT/sshkeys)

`sshkeys` provides utilities for parsing and marshalling cryptographic keys used for SSH, in both cleartext and encrypted formats.

[ssh.ParseRawPrivateKey](https://godoc.org/golang.org/x/crypto/ssh#ParseRawPrivateKey) only supports parsing a subset of the formats `sshkeys` supports, does not support parsing encrypted private keys, and does not support marshalling.

## Supported Formats

* OpenSSH's [PROTOCOL.key](https://github.com/openssh/openssh-portable/blob/master/PROTOCOL.key) for RSA and ED25519 keys.
* OpenSSH version >= 7.6 using aes256-ctr encryption
* "Classic" PEM containing RSA (PKCS#1), DSA (OpenSSL), and ECDSA private keys.
