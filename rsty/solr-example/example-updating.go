package main

import "fmt"

import "rsty/solr"

/*
 * README
 * ------
 * This example shows an update Query being performed. An update document
 * is built and sent off to Solr.
 */


func main() {

    // init a connection
    s, err := solr.Init("localhost", 8983)

    if err != nil {
        fmt.Println(err)
        return
    }

    // build an update document, in this case adding two documents
    f := map[string]interface {}{
        "add":[]interface {}{
            map[string]interface {}{"id": 22, "title": "abc"},
            map[string]interface {}{"id": 23, "title": "def"},
            map[string]interface {}{"id": 24, "title": "def"},
        },
    }

    // send off the update (2nd parameter indicates we also want to commit the operation)
    resp, err := s.Update(f, true)

    if err != nil {
        fmt.Println("error =>", err)
    } else {
        fmt.Println("resp =>", resp)
    }
}
