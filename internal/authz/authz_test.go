package authz

import "testing"

func TestRoleMatrix(t *testing.T) {
	tests := []struct {
		role       Role
		capability Capability
		want       bool
	}{
		{RoleUser, LibraryRead, true}, {RoleUser, MetadataWrite, false}, {RoleUser, MediaDownload, false},
		{RoleModerator, MetadataWrite, true}, {RoleModerator, LibraryDestructive, false}, {RoleModerator, ScraperRun, false},
		{RoleAdmin, MediaDownload, true}, {RoleAdmin, DatabaseSQL, true}, {RoleAdmin, AccountManage, true},
	}
	for _, test := range tests {
		p := Principal{UserID: "u1", Role: test.role, Status: StatusActive}
		if got := p.Allows(test.capability); got != test.want {
			t.Errorf("%s allows %s=%v, want %v", test.role, test.capability, got, test.want)
		}
	}
}

func TestFailClosed(t *testing.T) {
	for _, p := range []Principal{
		{},
		{UserID: "u1", Role: "UNKNOWN", Status: StatusActive},
		{UserID: "u1", Role: RoleAdmin, Status: StatusDisabled},
		{UserID: "u1", Role: RoleAdmin, Status: StatusPasswordChangeRequired},
	} {
		if p.Allows(LibraryRead) {
			t.Fatalf("principal unexpectedly authorized: %#v", p)
		}
	}
	if (Principal{UserID: "u1", Role: RoleAdmin, Status: StatusActive}).Allows(Capability("invented.action")) {
		t.Fatal("unknown capability authorized")
	}
}

func TestTokenScopesOnlyReduce(t *testing.T) {
	p := Principal{UserID: "u1", Role: RoleModerator, Status: StatusActive, TokenScopes: capabilitySet(LibraryRead, DatabaseSQL)}
	if !p.Allows(LibraryRead) {
		t.Fatal("approved reduced scope denied")
	}
	if p.Allows(MetadataWrite) {
		t.Fatal("omitted role capability survived token reduction")
	}
	if p.Allows(DatabaseSQL) {
		t.Fatal("token scope escalated beyond role")
	}
}

func TestEffectiveCapabilitiesHonorTokenReduction(t *testing.T) {
	p := Principal{UserID: "u1", Role: RoleModerator, Status: StatusActive, TokenScopes: capabilitySet(LibraryRead)}
	got := p.EffectiveCapabilities()
	if len(got) != 1 || got[0] != LibraryRead {
		t.Fatalf("effective capabilities=%v, want [%s]", got, LibraryRead)
	}
}

func TestOwnershipSeparateFromCapability(t *testing.T) {
	p := Principal{UserID: "u1", Role: RoleAdmin, Status: StatusActive}
	if err := RequireOwned(p, AccountSelfRead, "u2"); err == nil {
		t.Fatal("admin bypassed private ownership")
	}
	if err := RequireOwned(p, AccountSelfRead, "u1"); err != nil {
		t.Fatal(err)
	}
}

func TestPublicBootstrapNeverGrantedToAccount(t *testing.T) {
	p := Principal{UserID: "u1", Role: RoleAdmin, Status: StatusActive}
	if p.Allows(PublicBootstrap) {
		t.Fatal("bootstrap capability granted to authenticated admin")
	}
}

func TestStructuredAuthenticationAndForbiddenErrors(t *testing.T) {
	if err := Require(Principal{}, LibraryRead); err == nil {
		t.Fatal("anonymous principal was authorized")
	} else if client, ok := err.(ClientError); !ok || client.Code() != CodeUnauthenticated || client.HTTPStatus() != 401 {
		t.Fatalf("anonymous error=%T %v", err, err)
	}
	user := Principal{UserID: "u1", Role: RoleUser, Status: StatusActive}
	if err := Require(user, MetadataWrite); err == nil {
		t.Fatal("forbidden principal was authorized")
	} else if client, ok := err.(ClientError); !ok || client.Code() != CodeForbidden || client.HTTPStatus() != 403 {
		t.Fatalf("forbidden error=%T %v", err, err)
	}
}

func TestPasswordChangeRequiredIsNarrow(t *testing.T) {
	p := Principal{UserID: "u1", Role: RoleAdmin, Status: StatusPasswordChangeRequired}
	if !p.CanChangeRequiredPassword() {
		t.Fatal("password-change-required principal cannot enter password change flow")
	}
	for _, capability := range AllCapabilities() {
		if p.Allows(capability) {
			t.Fatalf("password-change-required principal retained %s", capability)
		}
	}
}
