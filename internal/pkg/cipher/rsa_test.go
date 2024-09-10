package cipher

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

type rsaKey struct {
	privateKey     *rsa.PrivateKey
	privateKeyPath string
	publicKey      *rsa.PublicKey
	publicKeyPath  string
}

const (
	_rsaPKCS1PrivateKeyPath = "rsa-pkcs1-private-key.pem"
	_rsaPKCS1PublicKeyPath  = "rsa-pkcs1-public-key.pem"
	_rsaPKCS8PrivateKeyPath = "rsa-pkcs8-private-key.pem"
	_rsaPKIXPublicKeyPath   = "rsa-pkix-public-key.pem"
)

func TestDecryptRSA(t *testing.T) {
	valid, _ := genPKCS1RsaKeys(_rsaPKCS1PrivateKeyPath, _rsaPKCS1PublicKeyPath)
	invalid, _ := genPKCS8PKIXRsaKeys(_rsaPKCS8PrivateKeyPath, _rsaPKIXPublicKeyPath)

	defer func() {
		assert.NoError(t, valid.clean())
		assert.NoError(t, invalid.clean())
	}()

	tests := []struct {
		name       string
		privateKey *rsa.PrivateKey
		publicKey  *rsa.PublicKey
		rawData    []byte
		wantErr    bool
	}{
		{
			name:       "valid",
			privateKey: valid.privateKey,
			publicKey:  valid.publicKey,
			rawData:    []byte("test data"),
			wantErr:    false,
		},
		{
			name:       "encrypted data with invalid algorithm public key",
			privateKey: valid.privateKey,
			publicKey:  invalid.publicKey,
			rawData:    []byte("test data"),
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encryptedData, _ := EncryptRSA(tt.publicKey, tt.rawData)

			got, err := DecryptRSA(tt.privateKey, encryptedData)
			if tt.wantErr {
				require.Error(t, err)
				t.Log(err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.rawData, got)
		})
	}
}

func TestEncryptRSA(t *testing.T) {
	valid, _ := genPKCS1RsaKeys(_rsaPKCS1PrivateKeyPath, _rsaPKCS1PublicKeyPath)
	invalid, _ := genPKCS8PKIXRsaKeys(_rsaPKCS8PrivateKeyPath, _rsaPKIXPublicKeyPath)

	defer func() {
		assert.NoError(t, valid.clean())
		assert.NoError(t, invalid.clean())
	}()

	tests := []struct {
		name       string
		publicKey  *rsa.PublicKey
		privateKey *rsa.PrivateKey
		rawData    []byte
		want       []byte
		wantErr    bool
	}{
		{
			name:       "valid",
			publicKey:  valid.publicKey,
			privateKey: valid.privateKey,
			rawData:    []byte("test data"),
			wantErr:    false,
		},
		{
			name:       "too large data to encrypt",
			privateKey: valid.privateKey,
			publicKey:  valid.publicKey,
			rawData:    bytes.Repeat([]byte("x"), valid.publicKey.Size()),
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncryptRSA(tt.publicKey, tt.rawData)
			if tt.wantErr {
				require.Error(t, err)
				t.Log(err)
				return
			}

			require.NoError(t, err)

			decryptedData, _ := DecryptRSA(tt.privateKey, got)

			require.Equal(t, decryptedData, tt.rawData)
		})
	}
}

func TestUploadRSAPrivateKey(t *testing.T) {
	keys1, err := genPKCS1RsaKeys(_rsaPKCS1PrivateKeyPath, _rsaPKCS1PublicKeyPath)
	require.NoError(t, err)

	keys2, err := genPKCS8PKIXRsaKeys(_rsaPKCS8PrivateKeyPath, _rsaPKIXPublicKeyPath)
	require.NoError(t, err)

	tmp, err := os.CreateTemp("", "invalid-key-*.pem")
	require.NoError(t, err)

	defer func() {
		assert.NoError(t, tmp.Close())
		assert.NoError(t, os.Remove(tmp.Name()))
		assert.NoError(t, keys1.clean())
		assert.NoError(t, keys2.clean())
	}()

	tests := []struct {
		name           string
		filename       string
		wantPrivateKey *rsa.PrivateKey
		wantErr        bool
	}{
		{
			name:           "valid crypto key",
			filename:       keys1.privateKeyPath,
			wantPrivateKey: keys1.privateKey,
			wantErr:        false,
		},
		{
			name:           "crypto key with invalid pem block",
			filename:       tmp.Name(),
			wantPrivateKey: nil,
			wantErr:        true,
		},
		{
			name:           "empty path",
			filename:       "",
			wantPrivateKey: nil,
			wantErr:        true,
		},
		{
			name:           "crypto key with invalid algorithm",
			filename:       keys2.privateKeyPath,
			wantPrivateKey: nil,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPrivateKey, err := UploadRSAPrivateKey(tt.filename)
			if tt.wantErr {
				require.Error(t, err)
				t.Log(err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantPrivateKey, gotPrivateKey)
		})
	}
}

func TestUploadRSAPublicKey(t *testing.T) {
	keys1, err := genPKCS1RsaKeys(_rsaPKCS1PrivateKeyPath, _rsaPKCS1PublicKeyPath)
	require.NoError(t, err)

	keys2, err := genPKCS8PKIXRsaKeys(_rsaPKCS8PrivateKeyPath, _rsaPKIXPublicKeyPath)
	require.NoError(t, err)

	tmp, err := os.CreateTemp("", "invalid-key-*.pem")
	require.NoError(t, err)

	defer func() {
		assert.NoError(t, tmp.Close())
		assert.NoError(t, os.Remove(tmp.Name()))
		assert.NoError(t, keys1.clean())
		assert.NoError(t, keys2.clean())
	}()

	tests := []struct {
		name          string
		filename      string
		wantPublicKey *rsa.PublicKey
		wantErr       bool
	}{
		{
			name:          "valid crypto key",
			filename:      keys1.publicKeyPath,
			wantPublicKey: keys1.publicKey,
			wantErr:       false,
		},
		{
			name:          "crypto key with invalid pem block",
			filename:      tmp.Name(),
			wantPublicKey: nil,
			wantErr:       true,
		},
		{
			name:          "empty path",
			filename:      "",
			wantPublicKey: nil,
			wantErr:       false,
		},
		{
			name:          "crypto key with invalid algorithm",
			filename:      keys2.publicKeyPath,
			wantPublicKey: nil,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPublicKey, err := UploadRSAPublicKey(tt.filename)
			if tt.wantErr {
				require.Error(t, err)
				t.Log(err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantPublicKey, gotPublicKey)
		})
	}
}

func TestGenerateRSAKey(t *testing.T) {
	keys1, err := genPKCS1RsaKeys(_rsaPKCS1PrivateKeyPath, _rsaPKCS1PublicKeyPath)
	require.NoError(t, err)
	require.NotNil(t, keys1)
	require.NotNil(t, keys1.privateKey)
	require.Contains(t, keys1.privateKeyPath, ".pem")
	require.NotNil(t, keys1.publicKey)
	require.Contains(t, keys1.publicKeyPath, ".pem")

	keys2, err := genPKCS8PKIXRsaKeys(_rsaPKCS8PrivateKeyPath, _rsaPKIXPublicKeyPath)
	require.NoError(t, err)
	require.NotNil(t, keys2)
	require.NotNil(t, keys2.privateKey)
	require.Contains(t, keys2.privateKeyPath, ".pem")
	require.NotNil(t, keys2.publicKey)
	require.Contains(t, keys2.publicKeyPath, ".pem")

	defer func() {
		require.NoError(t, keys1.clean())
		require.NoError(t, keys2.clean())
	}()
}

func genPKCS1RsaKeys(privateKeyPath, publicKeyPath string) (*rsaKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	privateKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)

	publicKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
		},
	)

	privKeyFile, err := os.OpenFile(privateKeyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = privKeyFile.Close()
	}()

	_, err = privKeyFile.Write(privateKeyPEM)
	if err != nil {
		return nil, err
	}

	pubKeyFile, err := os.OpenFile(publicKeyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = pubKeyFile.Close()
	}()

	_, err = pubKeyFile.Write(publicKeyPEM)
	if err != nil {
		return nil, err
	}

	return &rsaKey{
		privateKey:     privateKey,
		privateKeyPath: privKeyFile.Name(),
		publicKey:      &privateKey.PublicKey,
		publicKeyPath:  pubKeyFile.Name(),
	}, nil
}

func genPKCS8PKIXRsaKeys(privateKeyPath, publicKeyPath string) (*rsaKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	privKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	privateKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privKeyBytes,
		},
	)

	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}

	publicKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubKeyBytes,
		},
	)

	privKeyFile, err := os.OpenFile(privateKeyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = privKeyFile.Close()
	}()

	_, err = privKeyFile.Write(privateKeyPEM)
	if err != nil {
		return nil, err
	}

	pubKeyFile, err := os.OpenFile(publicKeyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = pubKeyFile.Close()
	}()

	_, err = pubKeyFile.Write(publicKeyPEM)
	if err != nil {
		return nil, err
	}

	return &rsaKey{
		privateKey:     privateKey,
		privateKeyPath: privKeyFile.Name(),
		publicKey:      &privateKey.PublicKey,
		publicKeyPath:  pubKeyFile.Name(),
	}, nil
}

func (k *rsaKey) clean() error {
	if err := os.Remove(k.privateKeyPath); err != nil {
		return err
	}

	if err := os.Remove(k.publicKeyPath); err != nil {
		return err
	}

	return nil
}
