package enum

import (
	"encoding/json"
	"fmt"
	"testing"
)

type Role int

var (
	UnknownRole = New[Role]("Unknown") // 0
	Admin       = New[Role]("Admin")   // 1
	User        = New[Role]("User")    // 2
	Guest       = New[Role]("Guest")   // 3
)

type Permission int

var (
	UnknownPermission = New[Permission]("Unknown") // 0
	Read              = New[Permission]("Read")    // 1
	Write             = New[Permission]("Write")   // 2
)

func acceptsRoleOnly(t *testing.T, role *Enum[Role]) {
	t.Log(role)
}

func acceptsRoleIDOnly(t *testing.T, id Role) {
	t.Log(id)
}

func acceptsPermissionOnly(t *testing.T, permission *Enum[Permission]) {
	t.Log(permission)
}

func acceptsPermissionIDOnly(t *testing.T, id Permission) {
	t.Log(id)
}

func TestEnum(t *testing.T) {
	acceptsRoleOnly(t, UnknownRole)
	acceptsRoleOnly(t, Admin)
	acceptsRoleOnly(t, User)
	acceptsRoleOnly(t, Guest)

	// acceptsRoleOnly(t, UnknownPermission) // compile error

	acceptsRoleIDOnly(t, UnknownRole.ID())
	acceptsRoleIDOnly(t, Admin.ID())
	acceptsRoleIDOnly(t, User.ID())
	acceptsRoleIDOnly(t, Guest.ID())

	// acceptsRoleIDOnly(t, UnknownPermission.ID()) // compile error

	acceptsPermissionOnly(t, UnknownPermission)
	acceptsPermissionOnly(t, Read)
	acceptsPermissionOnly(t, Write)

	// acceptsPermissionOnly(t, UnknownRole) // compile error

	acceptsPermissionIDOnly(t, UnknownPermission.ID())
	acceptsPermissionIDOnly(t, Read.ID())
	acceptsPermissionIDOnly(t, Write.ID())

	// acceptsPermissionIDOnly(t, UnknownRole.ID()) // compile error
}

func TestEnum_Overflow(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, got normal execution")
		}
	}()

	type int8Enum int8

	// We can only have 128 int8 enums.
	for i := 0; i <= 128; i++ {
		New[int8Enum](fmt.Sprintf("Enum%d", i))
	}
}

func TestEnum_MarshalUnmarshal(t *testing.T) {
	data, err := json.Marshal(Guest)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	var newGuest Enum[Role]
	err = json.Unmarshal(data, &newGuest)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if newGuest.Name() != Guest.Name() {
		t.Errorf("expected %s, got %s", Guest.Name(), newGuest.Name())
	}

	if newGuest.ID() != Guest.ID() {
		t.Errorf("expected %d, got %d", Guest.ID(), newGuest.ID())
	}
}

func TestEnum_Switch(t *testing.T) {
	// Unsing role pointers, which should be the common case.
	role := Admin

	switch role {
	case UnknownRole:
		t.Errorf("expected %s, got %s", role, UnknownRole)
	case Admin:
		// Just do not error out. This is what we want.
	case User:
		t.Errorf("expected %s, got %s", role, User)
	case Guest:
		t.Errorf("expected %s, got %s", role, Guest)
	default:
		t.Errorf("expected %s, got something else", role)
	}

	// Using IDs.

	switch roleID := role.ID(); roleID {
	case UnknownRole.ID():
		t.Errorf("expected %d, got %d", roleID, UnknownRole.ID())
	case Admin.ID():
		// Just do not error out. This is what we want.
	case User.ID():
		t.Errorf("expected %d, got %d", roleID, User.ID())
	case Guest.ID():
		t.Errorf("expected %d, got %d", roleID, Guest.ID())
	default:
		t.Errorf("expected %d, got something else", roleID)
	}
}

func TestEnum_EnumsForType(t *testing.T) {
	enums := EnumsForType[Role]()
	if len(enums) != 4 {
		t.Errorf("expected 4, got %d", len(enums))
	}
}
