package as_admin

import (
	"database/sql"

	"github.com/google/uuid"
	c_admin "github.com/o-rensa/iv/internal/core/admin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AdminStore struct {
	db *gorm.DB
}

func NewAdminStore(db *sql.DB) (*AdminStore, error) {
	adminStore := &AdminStore{}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err == nil {
		adminStore.db = gormDB
	}

	return adminStore, nil
}

func (as *AdminStore) CreateAdmin(admin *c_admin.Admin) error {
	result := as.db.Create(admin)

	return result.Error
}

func (as *AdminStore) GetAdminByUsername(username string) (*c_admin.Admin, error) {
	// find admin where username = username payload
	var a c_admin.Admin
	gdb := as.db.Where("username = ?", username).First(&a)

	return &a, gdb.Error
}

func (as *AdminStore) GetAdminByID(iD uuid.UUID) (*c_admin.Admin, error) {
	// find admin by ID
	var admin c_admin.Admin
	gdb := as.db.First(&admin, "id = ?", iD.String())
	return &admin, gdb.Error
}

func (as *AdminStore) UpdateAdminPassword(ID uuid.UUID, newHashedPassword string) error {
	// find admin by ID
	var admin c_admin.Admin
	gdb := as.db.First(&admin, "ID = ?", ID.String())
	if gdb.Error != nil {
		return gdb.Error
	}

	// update to new password
	res := as.db.Model(&admin).Update("hashedPassword", newHashedPassword)
	return res.Error
}
