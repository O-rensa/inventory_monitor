package pp_admin

type CreateAdminPayload struct {
	Username         string `json:"username" validate:"required"`
	UnhashedPassword string `json:"password" validate:"required"`
}
