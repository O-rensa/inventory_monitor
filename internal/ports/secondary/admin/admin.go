package ps_admin

import c_admin "github.com/o-rensa/iv/internal/core/admin"

type AdminStore interface {
	CreateAdmin(*c_admin.Admin) error
	GetAdminByUsername(string) error
}
