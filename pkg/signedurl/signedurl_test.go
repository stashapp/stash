package signedurl

import (
	"errors"
	"net/url"
	"strconv"
	"testing"
	"time"
)

func TestDerivePrefix(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		// Scene stream variants
		{"/scene/1/stream", "/scene/1/stream"},
		{"/scene/1/stream.mp4", "/scene/1/stream"},
		{"/scene/1/stream.webm", "/scene/1/stream"},
		{"/scene/1/stream.mkv", "/scene/1/stream"},
		{"/scene/1/stream.m3u8", "/scene/1/stream"},
		{"/scene/1/stream.mpd", "/scene/1/stream"},

		// HLS segments
		{"/scene/1/stream.m3u8/0.ts", "/scene/1/stream"},
		{"/scene/1/stream.m3u8/99.ts", "/scene/1/stream"},

		// DASH segments
		{"/scene/1/stream.mpd/5_v.webm", "/scene/1/stream"},
		{"/scene/1/stream.mpd/5_a.webm", "/scene/1/stream"},
		{"/scene/1/stream.mpd/init_v.webm", "/scene/1/stream"},
		{"/scene/1/stream.mpd/init_a.webm", "/scene/1/stream"},

		// Caption
		{"/scene/1/caption", "/scene/1/caption"},

		// Image paths
		{"/image/5/thumbnail", "/image/5/thumbnail"},
		{"/image/5/image", "/image/5/image"},

		// Gallery paths
		{"/gallery/3/cover", "/gallery/3/cover"},
		{"/gallery/3/preview", "/gallery/3/preview"},

		// Short paths
		{"/scene/1", "/scene/1"},
		{"/scene", "/scene"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := DerivePrefix(tt.path)
			if got != tt.expected {
				t.Errorf("DerivePrefix(%q) = %q, want %q", tt.path, got, tt.expected)
			}
		})
	}
}

func TestGenerateCredentialID(t *testing.T) {
	secret := []byte("test-secret-key")

	t.Run("deterministic", func(t *testing.T) {
		cid1 := GenerateCredentialID(secret, "alice")
		cid2 := GenerateCredentialID(secret, "alice")
		if cid1 != cid2 {
			t.Errorf("expected deterministic output, got %q and %q", cid1, cid2)
		}
	})

	t.Run("different usernames produce different cids", func(t *testing.T) {
		cid1 := GenerateCredentialID(secret, "alice")
		cid2 := GenerateCredentialID(secret, "bob")
		if cid1 == cid2 {
			t.Errorf("expected different cids for different usernames, both got %q", cid1)
		}
	})

	t.Run("different secrets produce different cids", func(t *testing.T) {
		cid1 := GenerateCredentialID([]byte("secret-1"), "alice")
		cid2 := GenerateCredentialID([]byte("secret-2"), "alice")
		if cid1 == cid2 {
			t.Errorf("expected different cids for different secrets, both got %q", cid1)
		}
	})

	t.Run("length is 16 hex chars", func(t *testing.T) {
		cid := GenerateCredentialID(secret, "alice")
		if len(cid) != 16 {
			t.Errorf("expected length 16, got %d (%q)", len(cid), cid)
		}
	})
}

func TestSignAndVerifyRoundtrip(t *testing.T) {
	secret := []byte("test-secret-key")
	cid := GenerateCredentialID(secret, "alice")
	expires := time.Now().Add(1 * time.Hour)

	params := SignPrefix("/scene/1/stream", secret, cid, expires)

	gotCID, err := VerifyURL("/scene/1/stream", params, secret)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotCID != cid {
		t.Errorf("expected cid %q, got %q", cid, gotCID)
	}
}

func TestDifferentPathSamePrefixVerifies(t *testing.T) {
	secret := []byte("test-secret-key")
	cid := GenerateCredentialID(secret, "alice")
	expires := time.Now().Add(1 * time.Hour)

	// Sign for the stream prefix
	params := SignPrefix("/scene/1/stream", secret, cid, expires)

	// Verify with different paths that share the same prefix
	paths := []string{
		"/scene/1/stream",
		"/scene/1/stream.mp4",
		"/scene/1/stream.m3u8",
		"/scene/1/stream.m3u8/0.ts",
		"/scene/1/stream.mpd",
		"/scene/1/stream.mpd/5_v.webm",
		"/scene/1/stream.mpd/init_a.webm",
	}

	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			gotCID, err := VerifyURL(path, params, secret)
			if err != nil {
				t.Fatalf("unexpected error for path %q: %v", path, err)
			}
			if gotCID != cid {
				t.Errorf("expected cid %q, got %q", cid, gotCID)
			}
		})
	}
}

func TestDifferentPrefixFails(t *testing.T) {
	secret := []byte("test-secret-key")
	cid := GenerateCredentialID(secret, "alice")
	expires := time.Now().Add(1 * time.Hour)

	params := SignPrefix("/scene/1/stream", secret, cid, expires)

	// Different scene ID
	_, err := VerifyURL("/scene/2/stream", params, secret)
	if !errors.Is(err, ErrInvalidSignature) {
		t.Errorf("expected ErrInvalidSignature, got %v", err)
	}

	// Different resource type
	_, err = VerifyURL("/scene/1/caption", params, secret)
	if !errors.Is(err, ErrInvalidSignature) {
		t.Errorf("expected ErrInvalidSignature, got %v", err)
	}

	// Different entity type
	_, err = VerifyURL("/image/1/stream", params, secret)
	if !errors.Is(err, ErrInvalidSignature) {
		t.Errorf("expected ErrInvalidSignature, got %v", err)
	}
}

func TestExpiredURLFails(t *testing.T) {
	secret := []byte("test-secret-key")
	cid := GenerateCredentialID(secret, "alice")
	expires := time.Now().Add(-1 * time.Hour) // expired 1 hour ago

	params := SignPrefix("/scene/1/stream", secret, cid, expires)

	_, err := VerifyURL("/scene/1/stream", params, secret)
	if !errors.Is(err, ErrExpiredURL) {
		t.Errorf("expected ErrExpiredURL, got %v", err)
	}
}

func TestTamperedSignatureFails(t *testing.T) {
	secret := []byte("test-secret-key")
	cid := GenerateCredentialID(secret, "alice")
	expires := time.Now().Add(1 * time.Hour)

	params := SignPrefix("/scene/1/stream", secret, cid, expires)
	params.Set(SigParam, "tampered0000000000000000000000000000000000000000000000000000000")

	_, err := VerifyURL("/scene/1/stream", params, secret)
	if !errors.Is(err, ErrInvalidSignature) {
		t.Errorf("expected ErrInvalidSignature, got %v", err)
	}
}

func TestTamperedCIDFails(t *testing.T) {
	secret := []byte("test-secret-key")
	cid := GenerateCredentialID(secret, "alice")
	expires := time.Now().Add(1 * time.Hour)

	params := SignPrefix("/scene/1/stream", secret, cid, expires)
	params.Set(CIDParam, "tamperedcid12345")

	_, err := VerifyURL("/scene/1/stream", params, secret)
	if !errors.Is(err, ErrInvalidSignature) {
		t.Errorf("expected ErrInvalidSignature, got %v", err)
	}
}

func TestMissingParamsFails(t *testing.T) {
	secret := []byte("test-secret-key")
	cid := GenerateCredentialID(secret, "alice")
	expires := time.Now().Add(1 * time.Hour)

	full := SignPrefix("/scene/1/stream", secret, cid, expires)

	tests := []struct {
		name    string
		missing string
	}{
		{"missing cid", CIDParam},
		{"missing expires", ExpiresParam},
		{"missing signature", SigParam},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := make(url.Values)
			for k, v := range full {
				params[k] = v
			}
			params.Del(tt.missing)

			_, err := VerifyURL("/scene/1/stream", params, secret)
			if !errors.Is(err, ErrMissingParams) {
				t.Errorf("expected ErrMissingParams, got %v", err)
			}
		})
	}
}

func TestTamperedExpiresFails(t *testing.T) {
	secret := []byte("test-secret-key")
	cid := GenerateCredentialID(secret, "alice")
	expires := time.Now().Add(1 * time.Hour)

	params := SignPrefix("/scene/1/stream", secret, cid, expires)

	// Attacker extends the expiry by 24 hours
	tampered := time.Now().Add(25 * time.Hour)
	params.Set(ExpiresParam, strconv.FormatInt(tampered.Unix(), 10))

	_, err := VerifyURL("/scene/1/stream", params, secret)
	if !errors.Is(err, ErrInvalidSignature) {
		t.Errorf("expected ErrInvalidSignature, got %v", err)
	}
}

func TestWrongSecretFails(t *testing.T) {
	secret := []byte("test-secret-key")
	wrongSecret := []byte("wrong-secret-key")
	cid := GenerateCredentialID(secret, "alice")
	expires := time.Now().Add(1 * time.Hour)

	params := SignPrefix("/scene/1/stream", secret, cid, expires)

	_, err := VerifyURL("/scene/1/stream", params, wrongSecret)
	if !errors.Is(err, ErrInvalidSignature) {
		t.Errorf("expected ErrInvalidSignature, got %v", err)
	}
}
