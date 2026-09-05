package authn

import (
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestPasswordHashAndVerify(t *testing.T) {
	hash, err := HashPassword("correct horse battery staple")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(hash, "$argon2id$v=19$") || strings.Contains(hash, "correct horse") {
		t.Fatal("password hash is not a non-leaking versioned Argon2id string")
	}
	if ok, err := VerifyPassword(hash, "correct horse battery staple"); err != nil || !ok {
		t.Fatalf("correct password ok=%v err=%v", ok, err)
	}
	if ok, err := VerifyPassword(hash, "wrong"); err != nil || ok {
		t.Fatalf("wrong password ok=%v err=%v", ok, err)
	}
}

func TestPasswordHashUsesRandomSalt(t *testing.T) {
	a, _ := HashPassword("same password")
	b, _ := HashPassword("same password")
	if a == b {
		t.Fatal("password hashes reused a salt")
	}
}

func TestLegacyBcryptRequestsUpgrade(t *testing.T) {
	legacy, err := bcrypt.GenerateFromPassword([]byte("legacy password"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	valid, upgrade, err := VerifyPasswordWithUpgrade(string(legacy), "legacy password")
	if err != nil || !valid || !upgrade {
		t.Fatalf("valid=%v upgrade=%v err=%v", valid, upgrade, err)
	}
	valid, upgrade, err = VerifyPasswordWithUpgrade(string(legacy), "wrong")
	if err != nil || valid || upgrade {
		t.Fatalf("wrong valid=%v upgrade=%v err=%v", valid, upgrade, err)
	}
}
