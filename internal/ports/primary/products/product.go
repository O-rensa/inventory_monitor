package pp_product

type CreateProductPayload struct {
	ProductCategoryID string `json:"productCategoryID" validate:"required"`
	BrandID           string `json:"brandID" validate:"required"`
	ProductName       string `json:"productName" validate:"required"`
}

type UpdateProductPayload struct {
	ProductID         string `json:"productID" validate:"required"`
	ProductCategoryID string `json:"productCategoryID" validate:"required"`
	BrandID           string `json:"brandID" validate:"required"`
	ProductName       string `json:"productName" validate:"required"`
}

type ProductDto struct {
	ProductID           string `json:"productID"`
	ProductCategoryID   string `json:"productCategoryID"`
	ProductCategoryName string `json:"productCategoryName"`
	BrandID             string `json:"brandID"`
	BrandName           string `json:"brandName"`
	ProductName         string `json:"productName"`
	ItemCount           string `json:"itemCount"`
}
