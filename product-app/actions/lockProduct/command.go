package lockProduct

type Command struct {
	ProductUUID   string `json:"product_uuid"`
	LockingEntity string `json:"locking_entity"`
}
