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
    q := "q=*:*"
    
    // perform a query    
    res, err := s.SelectRaw(q)

    if err != nil {
        fmt.Println("hi", err)
        return
    }

    results := res.Results

    // print a summary and loop over results, priting the "title" and "latlng" fields
    fmt.Println(
        fmt.Sprintf("Query: %s\nHits: %d\nNum Results: %d\nQtime: %d\nStatus: %d\n\nResults\n-------\n",
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