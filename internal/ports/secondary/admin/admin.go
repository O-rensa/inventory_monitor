package ps_admin

import (
	"github.com/google/uuid"
	c_admin "github.com/o-rensa/iv/internal/core/admin"
)

type AdminStore interface {
	CreateAdmin(*c_admin.Admin) error
	GetAdminByUsername(string) error
	GetAdminByID(uuid.UUID) (*c_admin.Admin, error)
	UpdateAdminPassword(uuid.UUID, string) error
}
