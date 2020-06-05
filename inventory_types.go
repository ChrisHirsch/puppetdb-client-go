package puppetdb

/*
Inventory struct
A representation of an inventory query response.



More details here: https://puppet.com/docs/puppetdb/5.2/api/query/v4/inventory.html#response-format
*/
type Inventory struct {
	Certname    string            `json:"certname"`
	Timestamp   string            `json:"timestamp"`
	Environment string            `json:"environment"`
	Facts       map[string]string `json:"facts"`
	Trusted     map[string]string `json:"trusted"`
}
