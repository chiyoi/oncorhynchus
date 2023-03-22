package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"

	"github.com/chiyoi/oncorhynchus"
)

var config = &oncorhynchus.Config

func Encrypt(msg []byte, key *rsa.PublicKey) (cipher []byte) {
	cipher, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, key, msg, nil)
	if err != nil {

	}
	return
}

func Decrypt(cipher []byte, key *rsa.PrivateKey) (msg []byte, err error) {
	return rsa.DecryptOAEP(sha256.New(), nil, key, cipher, nil)
}

func GetKey() (key *rsa.PublicKey, err error) {
	bs, err := base64.StdEncoding.DecodeString(config.PublicKey)
	if err != nil {
		return
	}
	return x509.ParsePKCS1PublicKey(bs)
}

func EncryptPassword(password string, key *rsa.PublicKey) (token string) {
	sha := sha256.Sum256([]byte(password))
	token = base64.StdEncoding.EncodeToString(sha[:])
	cipher := Encrypt([]byte(token), key)
	return base64.StdEncoding.EncodeToString(cipher)
}
