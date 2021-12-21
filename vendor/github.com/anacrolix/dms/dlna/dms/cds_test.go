package dms

import (
	"strings"
	"testing"
)

func TestEscapeObjectID(t *testing.T) {
	o := object{
		Path: "/some/file",
	}
	id := o.ID()
	if strings.ContainsAny(id, "/") {
		t.Skip("may not work with some players: object IDs contain '/'")
	}
}

func TestRootObjectID(t *testing.T) {
	if (object{Path: "/"}).ID() != "0" {
		t.FailNow()
	}
}

func TestRootParentObjectID(t *testing.T) {
	if (object{Path: "/"}).ParentID() != "-1" {
		t.FailNow()
	}
}
