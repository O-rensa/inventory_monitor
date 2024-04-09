package c_admin

import c_sharedTypes "github.com/o-rensa/iv/internal/core/sharedtypes"

type Admin struct {
	*c_sharedTypes.ModelID
	*c_sharedTypes.BaseModel
	Username       string `gorm:"column:username;unique"`
	HashedPassword string `gorm:"hashedPassword"`
}
