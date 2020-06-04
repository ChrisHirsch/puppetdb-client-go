package puppetdb

/*
FactsWireFormat struct for submitting the 'replace facts' command to PuppetDB.
More details: https://puppet.com/docs/puppetdb/5.2/api/wire_format/facts_format_v5.html
*/
type FactsWireFormat struct {
	// Certificate name of node to replace facts for
	Name string `json:"name"`
	// A map of fact key/value pairs
	Values map[string]string `json:"values"`
}

/*
Fact response based query end-points.
https://puppet.com/docs/puppetdb/5.2/api/query/v4/facts.html#query-fields
*/
type Fact struct {
	Certname    string `json:"certname"`
	Name        string `json:"name"`
	Value       string `json:"value"`
	Environment string `json:"pupenv"`
}
