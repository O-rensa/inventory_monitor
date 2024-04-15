package pp_productCategory

type CreateProductCategoryPayload struct {
	ProductCategoryName string `json:"productCategoryName"`
}

type ProductCategoryDto struct {
	ProductCategoryID   string `json:"iD"`
	ProductCategoryName string `json:"productCategoryName"`
}

type UpdateProductCategoryPayload struct {
	ProductCategoryID   string `json:"iD"`
	ProductCategoryName string `json:"productCategoryName"`
}
