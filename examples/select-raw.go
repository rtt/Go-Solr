package main

import (
	"fmt"
	"github.com/rtt/Go-Solr"
)

/*
 * README
 * ------
 * This example shows a Query being performed. A query is built up using
 * a Query type object, the query executed and the results are then
 * printed to the console
 */

func main() {

	// init a connection
	s, err := solr.Init("localhost", 8983, "collection-name")

	if err != nil {
		fmt.Println(err)
		return
	}

	// define a solr query string
	q := "q=*:*"

	// perform a query    
	res, err := s.SelectRaw(q)

	if err != nil {
		fmt.Println(err)
		return
	}

	// grab results for ease of use later on
	results := res.Results

	// print a summary and loop over results, priting the "title" and "latlng" fields
	fmt.Println(
		fmt.Sprintf("Query: %#v\nHits: %d\nNum Results: %d\nQtime: %d\nStatus: %d\n\nResults\n-------\n",
			q,
			results.NumFound,
			results.Len(),
			res.QTime,
			res.Status))

	for i := 0; i < results.Len(); i++ {
		fmt.Println("Some field:", results.Get(i).Field("id"))
		fmt.Println("Some other field:", results.Get(i).Field("title"))

		fmt.Println("")
	}

}
