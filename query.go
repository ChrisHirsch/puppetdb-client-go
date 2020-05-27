package puppetdb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

var api_version = ""

/*
Query - Generic query function.
*/
func (server *Server) Query(url string) ([]byte, error) {
	baseURL := server.BaseURL

	fullURL := strings.Join([]string{baseURL, url}, "")

	req, err := http.NewRequest("GET", fullURL, server.Body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	// Set any additional headers such as authentication, proxy, etc
	for key, value := range server.Headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{Transport: server.HTTPTransport, Timeout: server.HTTPTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}

/*
QueryVersion queries the PuppetDB instance version end-point.

More details here: https://puppet.com/docs/puppetdb/5.2/api/meta/v1/version.html
*/
func (server *Server) QueryVersion() (*Version, error) {
	body, err := server.Query("pdb/meta/v1/version")
	if err != nil {
		return nil, err
	}

	var version Version
	json.Unmarshal(body, &version)

	return &version, err
}

/*
QueryServerTime - query the PuppetDB instance server-time end-point.

More details here: https://puppet.com/docs/puppetdb/latest/api/meta/v1/server-time.html
*/
func (server *Server) QueryServerTime() (*ServerTime, error) {
	body, err := server.Query("pdb/meta/v1/server-time")
	if err != nil {
		return nil, err
	}

	var serverTime ServerTime
	json.Unmarshal(body, &serverTime)

	return &serverTime, err
}

/*
QueryFactNames - Query the PuppetDB instance fact-names end-point.

More details here: http://docs.puppetlabs.com/puppetdb/latest/api/query/v3/fact-names.html
*/
func (server *Server) QueryFactNames() ([]string, error) {
	body, err := server.Query("pdb/query/v4/fact-names")
	if err != nil {
		return nil, err
	}

	var factNames []string
	json.Unmarshal(body, &factNames)

	return factNames, err
}

/*
QueryCatalogs - the PuppetDB instance catalogs end-point.

More details here: https://puppet.com/docs/puppetdb/5.2/api/query/v4/catalogs.html
*/
func (server *Server) QueryCatalogs(certname string) (*CatalogWireFormat, error) {
	url := fmt.Sprintf("pdb/query/v4/catalogs/%v", certname)
	body, err := server.Query(url)
	if err != nil {
		return nil, err
	}

	var catalog CatalogWireFormat
	json.Unmarshal(body, &catalog)

	return &catalog, err
}

/*
QueryFact will take in the fact and the query
*/
func (server *Server) QueryFact(fact string, queryElements ...string) (*[]Fact, error) {
	queryString := fact
	for i, query := range queryElements {
		if i == 0 {
			queryString += "?query=["
		} else {
			queryString += `,`
		}
		queryString += `"` + query + `"`

	}
	if len(queryElements) > 0 {
		queryString += "]"
	}

	log.Debugf("queryString=%s\n", queryString)
	return server.QueryFacts(queryString, nil)
}

/*
QueryFacts - Query the PuppetDB instance facts end-point.

More details here: https://puppet.com/docs/puppetdb/5.2/api/query/v4/facts.html
*/
func (server *Server) QueryFacts(queryString string, requestBody body) (*[]Fact, error) {
	//url := fmt.Sprintf("pdb/query/v4/facts?%v", queryString)
	url := fmt.Sprintf("pdb/query/v4/facts/%v", queryString)

	log.Debugf("url=%s\n", url)
	if requestBody != nil {
		server.Body = requestBody
	}

	body, err := server.Query(url)
	if err != nil {
		return nil, err
	}

	var facts []Fact
	json.Unmarshal(body, &facts)

	return &facts, err
}

/*
QueryFactsByName - Query the PuppetDB instance facts end-point.
Return all facts with the given fact name and value for ALL nodes.
ie: Only the certnmae field will differ

More details here: https://puppet.com/docs/puppetdb/5.2/api/query/v4/facts.html#pdbqueryv4factsfact-namevalue
*/
func (server *Server) QueryFactsByName(name string, queryString string) (*[]Fact, error) {
	url := fmt.Sprintf("pdb/query/v4/facts/%v/%v", name, queryString)

	body, err := server.Query(url)
	if err != nil {
		return nil, err
	}

	var facts []Fact
	json.Unmarshal(body, &facts)

	return &facts, err
}

/*
QueryFactsByNameValue - Query the PuppetDB instance facts end-point.

More details here: http://docs.puppetlabs.com/puppetdb/1.6/api/query/v3/facts.html#get-v3factsnamevalue
*/
func (server *Server) QueryFactsByNameValue(name string, value string, queryString string) (*[]Fact, error) {
	url := fmt.Sprintf("pdb/query/v4/facts/%v/%v?%v", name, value, queryString)

	body, err := server.Query(url)
	if err != nil {
		return nil, err
	}

	var facts []Fact
	json.Unmarshal(body, &facts)

	return &facts, err
}

/*
QueryResources - Query the PuppetDB instance resources end-point.

More details here: http://docs.puppetlabs.com/puppetdb/1.6/api/query/v3/resources.html#get-v3resources
*/
func (server *Server) QueryResources(queryString string) (*[]CatalogResource, error) {
	url := fmt.Sprintf("pdb/query/v4/resources/%v", queryString)

	body, err := server.Query(url)
	if err != nil {
		return nil, err
	}

	var resources []CatalogResource
	json.Unmarshal(body, &resources)

	return &resources, err
}

/*
QueryNodes - Query the PuppetDB instance nodes end-point.

More details here: https://puppet.com/docs/puppetdb/5.2/api/query/v4/nodes.html
*/
func (server *Server) QueryNodes(queryString string) (*[]Node, error) {
	url := fmt.Sprintf("pdb/query/v4/nodes/%v", queryString)

	body, err := server.Query(url)
	if err != nil {
		return nil, err
	}

	var nodes []Node
	json.Unmarshal(body, &nodes)

	return &nodes, err
}

/*
QueryReports - Query the PuppetDB instance reports end-point.

More details here: http://docs.puppetlabs.com/puppetdb/1.6/api/query/v3/reports.html#get-v3reports
*/
func (server *Server) QueryReports(queryString string) (*[]Report, error) {
	url := fmt.Sprintf("pdb/query/v4/reports/%v", queryString)

	body, err := server.Query(url)
	if err != nil {
		return nil, err
	}

	var reports []Report
	json.Unmarshal(body, &reports)

	return &reports, err
}

/*
QueryEvents - Query the PuppetDB instance events end-point.

More details here: https://puppet.com/docs/puppetdb/5.2/api/query/v4/events.html
*/
func (server *Server) QueryEvents(queryString string) (*[]Event, error) {
	url := fmt.Sprintf("pdb/query/v4/events/%v", queryString)

	body, err := server.Query(url)
	if err != nil {
		return nil, err
	}

	var event []Event
	json.Unmarshal(body, &event)

	return &event, err
}

/*
QueryEventCounts - Query the PuppetDB instance event-counts end-point.

More details here: https://puppet.com/docs/puppetdb/5.2/api/query/v4/event-counts.html
*/
func (server *Server) QueryEventCounts(queryString string) (*EventCounts, error) {
	url := fmt.Sprintf("pdb/query/v4/event-counts?%v", queryString)

	body, err := server.Query(url)
	if err != nil {
		return nil, err
	}

	var eventCounts EventCounts
	json.Unmarshal(body, &eventCounts)

	return &eventCounts, err
}

/*
QueryAggregateEventCounts - Query the PuppetDB instance aggregate-event-counts end-point.

More details here: https://puppet.com/docs/puppetdb/5.2/api/query/v4/aggregate-event-counts.html
*/
func (server *Server) QueryAggregateEventCounts(queryString string) (*AggregateEventCounts, error) {
	url := fmt.Sprintf("pdb/query/v4/aggregate-event-counts?%v", queryString)

	body, err := server.Query(url)
	if err != nil {
		return nil, err
	}

	var aec AggregateEventCounts
	json.Unmarshal(body, &aec)

	return &aec, err
}
