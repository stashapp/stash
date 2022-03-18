package session

import (
	"errors"
	"net/http"
	"testing"
)

type config struct {
	username                                   string
	password                                   string
	dangerousAllowPublicWithoutAuth            bool
	securityTripwireAccessedFromPublicInternet string
}

func (c *config) HasCredentials() bool {
	return c.username != "" && c.password != ""
}

func (c *config) GetDangerousAllowPublicWithoutAuth() bool {
	return c.dangerousAllowPublicWithoutAuth
}

func (c *config) GetSecurityTripwireAccessedFromPublicInternet() string {
	return c.securityTripwireAccessedFromPublicInternet
}

func (c *config) IsNewSystem() bool {
	return false
}

func TestCheckAllowPublicWithoutAuth(t *testing.T) {
	c := &config{}

	doTest := func(caseIndex int, r *http.Request, expectedErr interface{}) {
		t.Helper()
		err := CheckAllowPublicWithoutAuth(c, r)

		if expectedErr == nil && err == nil {
			return
		}

		if expectedErr == nil {
			t.Errorf("[%d]: unexpected error: %v", caseIndex, err)
			return
		}

		if !errors.As(err, expectedErr) {
			t.Errorf("[%d]: expected %T, got %v (%T)", caseIndex, expectedErr, err, err)
			return
		}
	}

	{
		// direct connection tests
		testCases := []struct {
			address string
			err     error
		}{
			{"192.168.1.1:8080", nil},
			{"192.168.1.1:8080", nil},
			{"100.64.0.1:8080", nil},
			{"127.0.0.1:8080", nil},
			{"[::1]:8080", nil},
			{"[fe80::c081:1c1a:ae39:d3cd%Ethernet 5]:9999", nil},
			{"193.168.1.1:8080", &ExternalAccessError{}},
			{"[2002:9fc4:ed97:e472:5170:5766:520c:c901]:9999", &ExternalAccessError{}},
		}

		// try with no X-FORWARDED-FOR and valid one
		xFwdVals := []string{"", "192.168.1.1"}

		for i, xFwdVal := range xFwdVals {
			header := make(http.Header)
			header.Set("X-FORWARDED-FOR", xFwdVal)

			for ii, tc := range testCases {
				r := &http.Request{
					RemoteAddr: tc.address,
					Header:     header,
				}

				doTest((i*len(testCases) + ii), r, tc.err)
			}
		}
	}

	{
		// X-FORWARDED-FOR
		testCases := []struct {
			proxyChain string
			err        error
		}{
			{"192.168.1.1, 192.168.1.2, 100.64.0.1, 127.0.0.1", nil},
			{"192.168.1.1, 193.168.1.1", &ExternalAccessError{}},
			{"193.168.1.1, 192.168.1.1", &ExternalAccessError{}},
		}

		const remoteAddr = "192.168.1.1:8080"

		header := make(http.Header)

		for i, tc := range testCases {
			header.Set("X-FORWARDED-FOR", tc.proxyChain)
			r := &http.Request{
				RemoteAddr: remoteAddr,
				Header:     header,
			}

			doTest(i, r, tc.err)
		}
	}

	{
		// test invalid request IPs
		invalidIPs := []string{"192.168.1.a:9999", "192.168.1.1"}

		for _, remoteAddr := range invalidIPs {
			r := &http.Request{
				RemoteAddr: remoteAddr,
			}

			err := CheckAllowPublicWithoutAuth(c, r)
			if err == nil {
				t.Errorf("[%s]: expected error", remoteAddr)
				continue
			}
		}
	}

	{
		// test overrides
		r := &http.Request{
			RemoteAddr: "193.168.1.1:8080",
		}

		c.username = "admin"
		c.password = "admin"

		if err := CheckAllowPublicWithoutAuth(c, r); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		c.username = ""
		c.password = ""

		c.dangerousAllowPublicWithoutAuth = true

		if err := CheckAllowPublicWithoutAuth(c, r); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}
}

func TestCheckExternalAccessTripwire(t *testing.T) {
	c := &config{}
	c.securityTripwireAccessedFromPublicInternet = "4.4.4.4"

	// always return nil if authentication configured or dangerous key set
	c.username = "admin"
	c.password = "admin"

	if err := CheckExternalAccessTripwire(c); err != nil {
		t.Errorf("unexpected error %v", err)
	}

	c.username = ""
	c.password = ""

	// HACK - this key isn't publically exposed
	c.dangerousAllowPublicWithoutAuth = true

	if err := CheckExternalAccessTripwire(c); err != nil {
		t.Errorf("unexpected error %v", err)
	}

	c.dangerousAllowPublicWithoutAuth = false

	if err := CheckExternalAccessTripwire(c); err == nil {
		t.Errorf("expected error %v", ExternalAccessError("4.4.4.4"))
	}

	c.securityTripwireAccessedFromPublicInternet = ""

	if err := CheckExternalAccessTripwire(c); err != nil {
		t.Errorf("unexpected error %v", err)
	}
}
