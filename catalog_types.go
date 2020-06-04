package puppetdb

/*
CatalogWireFormat - Wire format representation of a catalog.

You probably want to take a look at the NewCatalogWireFormat function, as this
is the suggested way to create a new catalog wire format data structure from
scratch.

More details here: http://docs.puppetlabs.com/puppetdb/latest/api/wire_format/catalog_format.html
*/
type CatalogWireFormat struct {
	// Metadata for this catalog
	Metadata CatalogMetadata `json:"metadata"`
	// Data for this catalog
	Data CatalogData `json:"data"`
}

/*
CatalogMetadata struct.

More details here: http://docs.puppetlabs.com/puppetdb/latest/api/wire_format/catalog_format.html#main-data-type-catalog
*/
type CatalogMetadata struct {
	// Catalog data API version
	APIVersion int `json:"api_version"`
}

/*
CatalogData struct

More details here: http://docs.puppetlabs.com/puppetdb/latest/api/wire_format/catalog_format.html#main-data-type-catalog
*/
type CatalogData struct {
	// Certificate name owning the catalog to be replaced
	Name string `json:"name"`
	// Version of the catalog
	Version string `json:"version"`
	// Unique identifier provided by client to marry catalogs with reports and other (future) objects
	TransactionUUID string `json:"transaction-uuid"`
	// Edges represented in this catalog
	Edges []CatalogEdge `json:"edges"`
	// Resources represented in this catalog
	Resources []CatalogResource `json:"resources"`
}

/*
CatalogEdge struct
A representation of an edge inside a catalog.

More details here: http://docs.puppetlabs.com/puppetdb/latest/api/wire_format/catalog_format.html#data-type-edge
*/
type CatalogEdge struct {
	// Source resource spec for this edge
	Source CatalogResourceSpec `json:"source"`
	// Target resource spec for this edge
	Target CatalogResourceSpec `json:"target"`
	// Relationship type
	Relationship string `json:"relationship"`
}

/*
CatalogResourceSpec struct represents a catalog resource reference for use in edges.

More details here: http://docs.puppetlabs.com/puppetdb/latest/api/wire_format/catalog_format.html#data-type-resource-spec
*/
type CatalogResourceSpec struct {
	// The type of a catalog resource
	Type string `json:"type"`
	// The title of a catalog resource
	Title string `json:"title"`
}

/*
CatalogResources - Collection of catalog resources
*/
type CatalogResources []CatalogResource

/*
CatalogResource struct

More details here: http://docs.puppetlabs.com/puppetdb/latest/api/wire_format/catalog_format.html#data-type-resource
*/
type CatalogResource struct {
	// The type of a catalog resource
	Type string `json:"type"`
	// The title of a catalog resource
	Title string `json:"title"`
	// Aliases for this resource
	//Aliases []string `json:"aliases"`
	// Exported status
	Exported bool `json:"exported"`
	// Source file this resource appears in
	File string `json:"file"`
	// Line in the file this resource appears in
	Line int `json:"line"`
	// All tags applied to this resource
	Tags []string `json:"tags"`
	// A map containing a list of key/value pairs for each parameter of this resource
	Parameters map[string]string `json:"parameters"`
}

/*
NewCatalogWireFormat - Create a new catalog
*/
func NewCatalogWireFormat() CatalogWireFormat {
	metadata := CatalogMetadata{0}
	data := CatalogData{"", "", "", nil, nil}
	return CatalogWireFormat{metadata, data}
}
