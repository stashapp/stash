// Package authz defines the fail-closed principal and capability foundation.
package authz

import "fmt"

const (
	CodeUnauthenticated = "UNAUTHENTICATED"
	CodeForbidden       = "FORBIDDEN"
)

type Role string

const (
	RoleUser      Role = "USER"
	RoleModerator Role = "MODERATOR"
	RoleAdmin     Role = "ADMIN"
)

type AccountStatus string

const (
	StatusActive                 AccountStatus = "ACTIVE"
	StatusDisabled               AccountStatus = "DISABLED"
	StatusPasswordChangeRequired AccountStatus = "PASSWORD_CHANGE_REQUIRED"
)

type Capability string

const (
	PublicBootstrap     Capability = "public.bootstrap"
	AccountSelfRead     Capability = "account.self.read"
	AccountSelfWrite    Capability = "account.self.write"
	PreferenceSelfWrite Capability = "preference.self.write"
	ActivitySelfWrite   Capability = "activity.self.write"
	LibraryRead         Capability = "library.read"
	MediaStream         Capability = "media.stream"
	MediaDownload       Capability = "media.download"
	MetadataWrite       Capability = "metadata.write"
	LibraryDestructive  Capability = "library.destructive"
	FilesystemWrite     Capability = "filesystem.write"
	AutomationRun       Capability = "automation.run"
	ScraperRun          Capability = "scraper.run"
	SystemStatusRead    Capability = "system.status.read"
	SystemConfigure     Capability = "system.configure"
	ExtensionRead       Capability = "extension.read"
	ExtensionManage     Capability = "extension.manage"
	JobRead             Capability = "job.read"
	JobManage           Capability = "job.manage"
	DataAdmin           Capability = "data.admin"
	DatabaseSQL         Capability = "database.sql"
	AccountManage       Capability = "account.manage"
	AuditRead           Capability = "audit.read"
	HashManage          Capability = "hash.manage"
)

var allCapabilities = []Capability{
	PublicBootstrap, AccountSelfRead, AccountSelfWrite, PreferenceSelfWrite, ActivitySelfWrite,
	LibraryRead, MediaStream, MediaDownload, MetadataWrite, LibraryDestructive, FilesystemWrite,
	AutomationRun, ScraperRun, SystemStatusRead, SystemConfigure, ExtensionRead, ExtensionManage,
	JobRead, JobManage, DataAdmin, DatabaseSQL, AccountManage, AuditRead, HashManage,
}

var roleCapabilities = map[Role]map[Capability]struct{}{
	RoleUser:      capabilitySet(AccountSelfRead, AccountSelfWrite, PreferenceSelfWrite, ActivitySelfWrite, LibraryRead, MediaStream),
	RoleModerator: capabilitySet(AccountSelfRead, AccountSelfWrite, PreferenceSelfWrite, ActivitySelfWrite, LibraryRead, MediaStream, MetadataWrite),
	RoleAdmin: capabilitySet(AccountSelfRead, AccountSelfWrite, PreferenceSelfWrite, ActivitySelfWrite, LibraryRead, MediaStream,
		MediaDownload, MetadataWrite, LibraryDestructive, FilesystemWrite, AutomationRun, ScraperRun, SystemStatusRead,
		SystemConfigure, ExtensionRead, ExtensionManage, JobRead, JobManage, DataAdmin, DatabaseSQL, AccountManage, AuditRead, HashManage),
}

type Principal struct {
	UserID      string
	Role        Role
	Status      AccountStatus
	TokenScopes map[Capability]struct{} // nil means an interactive session; non-nil can only reduce role grants.
}

func (p Principal) IsAuthenticated() bool {
	return p.UserID != "" && p.Status == StatusActive && (p.Role == RoleUser || p.Role == RoleModerator || p.Role == RoleAdmin)
}

func (p Principal) Allows(capability Capability) bool {
	if p.UserID == "" || p.Status != StatusActive || !IsKnownCapability(capability) || capability == PublicBootstrap {
		return false
	}
	grants, roleKnown := roleCapabilities[p.Role]
	if !roleKnown {
		return false
	}
	if _, allowed := grants[capability]; !allowed {
		return false
	}
	if p.TokenScopes != nil {
		_, allowed := p.TokenScopes[capability]
		return allowed
	}
	return true
}

func (p Principal) EffectiveCapabilities() []Capability {
	ret := make([]Capability, 0)
	for _, capability := range allCapabilities {
		if p.Allows(capability) {
			ret = append(ret, capability)
		}
	}
	return ret
}

// CanChangeRequiredPassword is narrower than account.self.write and is only
// for the current-password-confirmed password-change flow.
func (p Principal) CanChangeRequiredPassword() bool {
	return p.UserID != "" && p.Status == StatusPasswordChangeRequired && (p.Role == RoleUser || p.Role == RoleModerator || p.Role == RoleAdmin)
}

func (p Principal) Owns(ownerUserID string) bool {
	return p.UserID != "" && ownerUserID != "" && p.UserID == ownerUserID
}

func Require(p Principal, capability Capability) error {
	if p.UserID == "" || p.Status == StatusDisabled {
		return UnauthenticatedError{}
	}
	if !p.Allows(capability) {
		return DeniedError{Capability: capability}
	}
	return nil
}

func RequireOwned(p Principal, capability Capability, ownerUserID string) error {
	if err := Require(p, capability); err != nil {
		return err
	}
	if !p.Owns(ownerUserID) {
		return OwnershipError{}
	}
	return nil
}

func IsKnownCapability(capability Capability) bool {
	for _, known := range allCapabilities {
		if capability == known {
			return true
		}
	}
	return false
}

func AllCapabilities() []Capability { return append([]Capability(nil), allCapabilities...) }

type DeniedError struct{ Capability Capability }

func (e DeniedError) Error() string       { return fmt.Sprintf("capability denied: %s", e.Capability) }
func (DeniedError) Code() string          { return CodeForbidden }
func (DeniedError) HTTPStatus() int       { return 403 }
func (DeniedError) PublicMessage() string { return "forbidden" }

type OwnershipError struct{}

func (OwnershipError) Error() string         { return "resource owner mismatch" }
func (OwnershipError) Code() string          { return CodeForbidden }
func (OwnershipError) HTTPStatus() int       { return 403 }
func (OwnershipError) PublicMessage() string { return "forbidden" }

type UnauthenticatedError struct{}

func (UnauthenticatedError) Error() string         { return "valid authenticated principal required" }
func (UnauthenticatedError) Code() string          { return CodeUnauthenticated }
func (UnauthenticatedError) HTTPStatus() int       { return 401 }
func (UnauthenticatedError) PublicMessage() string { return "authentication required" }

type ClientError interface {
	error
	Code() string
	HTTPStatus() int
	PublicMessage() string
}

func capabilitySet(values ...Capability) map[Capability]struct{} {
	result := make(map[Capability]struct{}, len(values))
	for _, value := range values {
		result[value] = struct{}{}
	}
	return result
}
