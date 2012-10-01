package main

import "fmt"

import "rsty/solr"

/*
 * README
 * ------
 * This example shows a Query being performed. A query is built up using
 * a Query type object, the query executed and the results are then
 * printed to the console
 * 
 * Additionally, the returned facets are looped over and displayed.
 */


func main() {

    // init a connection
    s, err := solr.Init("localhost", 8983)

    if err != nil {
        fmt.Println(err)
        return
    }

    // Build a query object
    // Here we are specifying a 'q' param,
    // rows, faceting and facet.fields
    q := solr.Query{
        Params: solr.URLParamMap{
            "q": []string{"id:13"},
            "facet.field": []string{"accepts_4x4s", "accepts_bicycles"},
            "facet": []string{"true"},
        },
        Rows: 10,
    }

    // perform the query, checking for errors
    res, err := s.Select(&q)

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

    // Now display facets...

    fmt.Println("Facets")
    fmt.Println("------")

    // loop over facets returned
    
    for i := 0; i < results.NumFacets; i++ {
        facet := results.Facets[i]
        fmt.Println("Facet:", facet.Name)
        lfc := len(facet.Counts)
        for j := 0; j < lfc; j++ {
            fmt.Println(facet.Counts[j].Count, "=>", facet.Counts[j].Key)
        }
        fmt.Println("")
    }

}
