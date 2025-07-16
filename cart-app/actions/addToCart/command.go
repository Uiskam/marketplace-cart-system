package addToCart

type Command struct {
	CartUUID    string `json:"cart_uuid"`
	ProductUUID string `json:"product_uuid"`
}
