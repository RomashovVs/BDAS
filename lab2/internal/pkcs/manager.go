package pkcs

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
)

type IPair interface {
	GetPrivate() *rsa.PrivateKey
	GetPublic() *rsa.PublicKey
}

type Pair struct {
	private *rsa.PrivateKey
	public  *rsa.PublicKey
}

func (p *Pair) GetPrivate() *rsa.PrivateKey {
	return p.private
}

func (p *Pair) GetPublic() *rsa.PublicKey {
	return p.public
}

// GenerateKeyPair generates a new RSA key pair.
func GenerateKeyPair() (IPair, error) {
	private, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return &Pair{
		private: private,
		public:  &private.PublicKey,
	}, nil
}

// Encrypt encrypts a message using the recipient's public key.
func Encrypt(message string, publicKey *rsa.PublicKey) (string, error) {
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(message))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a message using the recipient's private key.
func Decrypt(ciphertext string, privateKey *rsa.PrivateKey) (string, error) {
	decodedCiphertext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, decodedCiphertext)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// Sign signs a message using the sender's private key.
func Sign(message string, privateKey *rsa.PrivateKey) (string, error) {
	hashed := sha256.Sum256([]byte(message))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

// Verify verifies a message signature using the sender's public key.
func Verify(message, signature string, publicKey *rsa.PublicKey) error {
	hashed := sha256.Sum256([]byte(message))
	decodedSignature, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}

	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], decodedSignature)
}
