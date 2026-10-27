package user

import (
	"testing"
)

func TestSHA256Hasher(t *testing.T) {
	hasher := &sha256Hasher{}

	password := "mysecretpassword"
	hash, err := hasher.GenerateHash(password)
	if err != nil {
		t.Fatalf("Error generating hash: %v", err)
	}

	match, err := hasher.CompareHash(password, hash)
	if err != nil {
		t.Fatalf("Error comparing hash: %v", err)
	}
	if !match {
		t.Fatal("Expected password to match hash, but it did not")
	}

	wrongPassword := "wrongpassword"
	match, err = hasher.CompareHash(wrongPassword, hash)
	if err != nil {
		t.Fatalf("Error comparing hash with wrong password: %v", err)
	}
	if match {
		t.Fatal("Expected wrong password to not match hash, but it did")
	}
}
