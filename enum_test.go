package enum

import (
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

func acceptsRoleOnly(t *testing.T, role Enum[Role]) {
	t.Log(role)
}

func acceptsRoleIDOnly(t *testing.T, id Role) {
	t.Log(id)
}

func acceptsPermissionOnly(t *testing.T, permission Enum[Permission]) {
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
