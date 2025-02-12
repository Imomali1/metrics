package cipher

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

func UploadRSAPublicKey(filename string) (publicKey *rsa.PublicKey, err error) {
	if filename == "" {
		return nil, nil
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		return publicKey, err
	}

	block, _ := pem.Decode(content)
	if block == nil {
		return publicKey, errors.New("failed to parse PEM block")
	}

	publicKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return publicKey, err
	}

	return publicKey, err
}

func UploadRSAPrivateKey(filename string) (privateKey *rsa.PrivateKey, err error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return privateKey, err
	}

	block, _ := pem.Decode(content)
	if block == nil {
		return privateKey, errors.New("failed to parse PEM block")
	}

	privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return privateKey, err
	}

	return privateKey, err
}

func EncryptRSA(publicKey *rsa.PublicKey, rawData []byte) ([]byte, error) {
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, rawData)
	if err != nil {
		return nil, err
	}

	return encryptedData, nil
}

func DecryptRSA(privateKey *rsa.PrivateKey, encryptedData []byte) ([]byte, error) {
	decryptedData, err := rsa.DecryptPKCS1v15(nil, privateKey, encryptedData)
	if err != nil {
		return nil, err
	}

	return decryptedData, nil
}
