package pp_admin

type CreateAdminPayload struct {
	Username         string `json:"username" validate:"required"`
	UnhashedPassword string `json:"password" validate:"required"`
}

type UpdateAdminPasswordPayload struct {
	ID                  string `json:"adminID" validate:"required"`
	Username            string `json:"username" validate:"required"`
	UnhashedOldPassword string `json:"oldPassword" validate:"required"`
	NewUnhashedPassword string `json:"newPassword" validate:"required"`
}

type LogInAdminPayload struct {
	Username         string `json:"username" validate:"required"`
	UnhashedPassword string `json:"password" validate:"required"`
}
