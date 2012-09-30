package main

import "fmt"

import "rsty/solr"


func main() {

    // init a connection
    s, err := solr.Init("localhost", 8983)

    if err != nil {
        fmt.Println(err)
        return
    }

    // define a solr query string
    q := "*:*"
    
    // perform a query    
    res, err := s.RawQuery(q)

    if err != nil {
        fmt.Println(err)
        return
    }

    // print a summary and loop over results, priting the "title" and "latlng" fields
    fmt.Println(fmt.Sprintf("Query: %s\nHits: %d\nNum Results: %d\n\nResults\n-------\n", q, res.NumFound, res.Len())) 

    for i := 0; i < res.Len(); i++ {
        fmt.Println("Title:", res.Get(i).Field("title"))
        fmt.Println("Latlng:", res.Get(i).Field("latlng"))

        fmt.Println("")
    }
}