package args

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"github.com/newrelic/infra-integrations-sdk/log"
)

// If RSAPrivateKey is present then we assume the password is RSA encrypted and base64 encoded
func decryptPassword(al ArgumentList) error {
	if al.RSAPrivateKey == "" {
		return nil
	}

	pemString, err := ioutil.ReadFile(al.RSAPrivateKey)
	if err != nil {
		return fmt.Errorf("decryptPassword: Unable to read private key: %s err: %s", al.RSAPrivateKey, err)
	}

	block, _ := pem.Decode(pemString)
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("decryptPassword: Unable to parse private key. err: %s", err)
	}

	ciphertext, err := base64.StdEncoding.DecodeString(al.Password)
	if err != nil {
		return fmt.Errorf("decryptPassword: Unable to base64 decode password. %s", err)
	}

	label := []byte("")
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, key, ciphertext, label)
	if err != nil {
		return fmt.Errorf("decryptPassword: Decryption error. err: %s", err)
	}

	al.Password = string(plaintext)
	log.Debug("decryptPassword: Decrypted password: %s", al.Password)
	return nil
}

func encryptPassword(al ArgumentList) string {
	if al.RSAPrivateKey == "" {
		log.Error("encrypt: rsa_private_key is a required parameter")
		return ""
	}

	pemString, err := ioutil.ReadFile(al.RSAPrivateKey)
	if err != nil {
		log.Error("encrypt: Unable to read private key: %s err: %s", al.RSAPrivateKey, err)
		return ""
	}

	block, _ := pem.Decode(pemString)
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Error("encrypt: Unable to parse private key. err: %s", err)
		return ""
	}

	label := []byte("")
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &key.PublicKey, []byte(al.Password), label)
	if err != nil {
		log.Error("encrypt: encryption error. err: %s", err)
		return ""
	}

	al.Password = base64.StdEncoding.EncodeToString([]byte(ciphertext))

	log.Debug("encrypt: encrypted password: %s", al.Password)
	return al.Password
}
