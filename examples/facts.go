package main

import (
	"fmt"
	"net/url"

	"github.com/ChrisHirsch/puppetdb-client-go"
)

func main() {
	server := puppetdb.NewServer("http://localhost:8080/")

	// A blank query string - just an example
	var values url.Values
	values = map[string][]string{
		"query": []string{""},
	}
	queryString := values.Encode()

	response, _ := server.QueryFacts(queryString)
	fmt.Printf("Fact Names: %v\n", response)
}
