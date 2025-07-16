package sellProduct

type Command struct {
	ProductUUIDs  []string `json:"product_uuids"`
	LockingEntity string   `json:"locking_entity"`
}
