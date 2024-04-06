package pp_brands

type CreateBrandPayload struct {
	BrandName string `json:"brandName" validate:"required"`
}

type BrandDto struct {
	BrandID   string `json:"iD"`
	BrandName string `json:"brandName"`
}

type UpdateBrandPayload struct {
	BrandID   string `json:"iD"`
	BrandName string `json:"brandName"`
}
