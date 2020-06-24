package puppetdb

/*
Inventory struct
A representation of an inventory query response.



More details here: https://puppet.com/docs/puppetdb/5.2/api/query/v4/inventory.html#response-format
*/
type Inventory struct {
	Certname    string                 `json:"certname"`
	Timestamp   string                 `json:"timestamp"`
	Environment string                 `json:"environment"`
	Facts       map[string]interface{} `json:"facts"`
	Trusted     map[string]interface{} `json:"trusted"`
}

/*
InventoryFact struct
Facts contained in the inventory
*/
type InventoryFact struct {
	Name  string                   `json:"name"`
	Value map[string]InventoryFact `json:"value"`
}
