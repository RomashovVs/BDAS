package main

import (
	"main/internal/pkcs"

	"go.uber.org/zap"
)

func main() {
	l, _ := zap.NewProduction()

	firstPair, err := pkcs.GenerateKeyPair()
	if err != nil {
		l.Fatal("Error generating Alice's key pair", zap.Error(err))
	}

	secondPair, err := pkcs.GenerateKeyPair()
	if err != nil {
		l.Fatal("Error generating Bob's key pair", zap.Error(err))
	}

	message := "Secret message for second"
	encryptedMessage, err := pkcs.Encrypt(message, secondPair.GetPublic())
	if err != nil {
		l.Fatal("Error encrypting message", zap.Error(err))
	}

	decryptedMessage, err := pkcs.Decrypt(encryptedMessage, secondPair.GetPrivate())
	if err != nil {
		l.Fatal("Error decrypting message", zap.Error(err))
	}

	l.Info("Original message", zap.String("message", message))
	l.Info("Encrypted message", zap.String("message", encryptedMessage))
	l.Info("Decrypted message", zap.String("message", decryptedMessage))

	signedMessage, err := pkcs.Sign(message, firstPair.GetPrivate())
	if err != nil {
		l.Fatal("Error signing message", zap.Error(err))
	}

	if err = pkcs.Verify(message, signedMessage, firstPair.GetPublic()); err != nil {
		l.Fatal("Signature verification failed", zap.Error(err))
		return
	}

	l.Info("Signature verified successfully!")
}
