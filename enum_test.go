package enum

import (
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

func acceptsRoleValueOnly(t *testing.T, value Role) {
	t.Log(value)
}

func acceptsPermissionOnly(t *testing.T, permission Enum[Permission]) {
	t.Log(permission)
}

func acceptsPermissionValueOnly(t *testing.T, value Permission) {
	t.Log(value)
}

func TestEnum(t *testing.T) {
	acceptsRoleOnly(t, UnknownRole)
	acceptsRoleOnly(t, Admin)
	acceptsRoleOnly(t, User)
	acceptsRoleOnly(t, Guest)

	// acceptsRoleOnly(t, UnknownPermission) // compile error

	acceptsRoleValueOnly(t, UnknownRole.Value())
	acceptsRoleValueOnly(t, Admin.Value())
	acceptsRoleValueOnly(t, User.Value())
	acceptsRoleValueOnly(t, Guest.Value())

	// acceptsRoleValueOnly(t, UnknownPermission.Value()) // compile error

	acceptsPermissionOnly(t, UnknownPermission)
	acceptsPermissionOnly(t, Read)
	acceptsPermissionOnly(t, Write)

	// acceptsPermissionOnly(t, UnknownRole) // compile error

	acceptsPermissionValueOnly(t, UnknownPermission.Value())
	acceptsPermissionValueOnly(t, Read.Value())
	acceptsPermissionValueOnly(t, Write.Value())

	// acceptsPermissionValueOnly(t, UnknownRole.Value()) // compile error
}
