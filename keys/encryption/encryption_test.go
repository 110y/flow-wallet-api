package encryption

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewCrypter(t *testing.T) {
	key := []byte("test123test123test123test123test")
	symmetricCrypter := NewCrypter(key)

	t.Run("key is saved", func(t *testing.T) {
		if symmetricCrypter.key == nil {
			t.Fatal("key was not set")
		}

		if !bytes.Equal(symmetricCrypter.key, key) {
			t.Fatal("keys do not match")
		}
	})
}

func TestSymmetricCrypter(t *testing.T) {
	key := []byte("testkeytestkeytestkeytestkeytest")
	original := []byte("some-secret-key")

	t.Run("fails with invalid key size", func(t *testing.T) {
		invalidKey := []byte("nope")
		crypter := NewCrypter(invalidKey)

		encValue, err := crypter.Encrypt([]byte("should-not-encrypt"))
		if err == nil {
			t.Fatal("expected error is missing")
		}

		expectedErrorMessagePrefix := "crypto/aes: invalid key size"
		if !strings.HasPrefix(err.Error(), expectedErrorMessagePrefix) {
			t.Errorf("got unexpected error prefix: %v - want: %v", err, expectedErrorMessagePrefix)
		}

		if len(encValue) != 0 {
			t.Errorf("expected encrypted value to be empty, got %v", encValue)
		}

	})

	t.Run("encrypts and decrypts a value", func(t *testing.T) {
		crypter := NewCrypter(key)
		encValue, err := crypter.Encrypt(original)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if len(encValue) == 0 {
			t.Error("encrypted value is empty")
		}

		if bytes.Equal(original, encValue) {
			t.Errorf("value was not encrypted: %v => %v", original, encValue)
		}

		decValue, err := crypter.Decrypt(encValue)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if !bytes.Equal(decValue, original) {
			t.Errorf("decrypted value does not match original: %v vs. %v", decValue, original)
		}
	})

	t.Run("decrypt fails with wrong key", func(t *testing.T) {
		crypter := NewCrypter(key)
		secondKey := []byte("failkeyfailkeyfailkeyfailkeyfail")
		secondCrypter := NewCrypter(secondKey)

		if bytes.Equal(key, secondKey) {
			t.Fatal("keys are equal")
		}

		encValue, err := crypter.Encrypt(original)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		expectedErrorMessage := "cipher: message authentication failed"
		wrongKeyDecValue, err := secondCrypter.Decrypt(encValue)

		if len(wrongKeyDecValue) != 0 {
			t.Fatalf("expected empty value, got: %v", wrongKeyDecValue)
		}

		if err == nil {
			t.Fatal("expected error is missing", err)
		}

		if !strings.Contains(err.Error(), expectedErrorMessage) {
			t.Errorf("unexpected error: %v - expected: %v", err, expectedErrorMessage)
		}

	})
}
