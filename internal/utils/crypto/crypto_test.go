package crypto

import (
	"bytes"
	"testing"
)

func TestEncryptGCM(t *testing.T) {
	// given
	encryptionKey := DerivePbkdf2From([]byte("testpassword"))
	testCases := []struct {
		Name string
		Data []byte
	}{
		{Name: "with empty data", Data: []byte("")},
		{Name: "with non-empty data", Data: []byte("testdata")},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// when
			ciphertext := EncryptGCM(encryptionKey, tc.Data)

			// then
			decrypted, err := DecryptGCM(encryptionKey, ciphertext)
			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
			if !bytes.Equal(decrypted, tc.Data) {
				t.Fatalf("expected %v, got %v", tc.Data, decrypted)
			}
		})
	}
}

func TestDecryptGCM(t *testing.T) {
	encryptionKey := DerivePbkdf2From([]byte("testpassword"))

	testCases := []struct {
		Name string
		Data []byte
	}{
		{Name: "with empty data", Data: []byte("")},
		{Name: "with non-empty data", Data: []byte("testdata")},
	}
	for _, tc := range testCases {
		// given
		ciphertext := EncryptGCM(encryptionKey, tc.Data)

		// when
		decrypted, err := DecryptGCM(encryptionKey, ciphertext)

		// then
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if !bytes.Equal(decrypted, tc.Data) {
			t.Fatalf("expected %v, got %v", tc.Data, decrypted)
		}
	}
}

func TestDerivePbkdf2From(t *testing.T) {
	t.Run("should return the output that differs from input", func(t *testing.T) {
		// given
		password := []byte("testpassword")

		// when
		encryptionKey := DerivePbkdf2From(password)

		// then
		if bytes.Equal(encryptionKey, password) {
			t.Fatal("output the same as input")
		}
	})

	t.Run("should not be reversable", func(t *testing.T) {
		// given
		password := []byte("testpassword")

		// when
		encryptionKey := DerivePbkdf2From(password)
		reversed := DerivePbkdf2From(encryptionKey)

		// then
		if bytes.Equal(reversed, password) {
			t.Fatal("the function is reversable")
		}
	})
}
