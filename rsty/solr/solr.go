package solr

import (
    "fmt"
    "encoding/json"
    "net/http"
    "io/ioutil"
)

/*
 * Represents a "connection"; actually just a host and port
 * (and probably at some point a Solr Core name)
 */
type Connection struct {
    Host string
    Port int
}

/* 
 * Represents a Solr document
 */
type Document struct {
    Fields map[string] interface{}
}

/*
 * Represents a collection of solr documents
 * and various other metrics
 */
type DocumentCollection struct {
    Collection []Document
    NumFound int
}

// type Query struct {

// }

func (d *DocumentCollection) Get(i int) *Document {
    return &d.Collection[i]
}

func (d *DocumentCollection) Len() int {
    return len(d.Collection)
}

func (document Document) Field(field string) interface{} {
    r, _ := document.Fields[field]
    return r
}

func (document Document) Doc() map[string] interface{} {
    return document.Fields
}

/*
 * Inits a new Connection
 * @returns *Connection
 */
func Init(host string, port int) (*Connection, error) {
    if host == "" || port <= 0 || port > 65535 {
        return nil, fmt.Errorf("Invalid host or port")
    }

    c := Connection{host, port}
    return &c, nil
}


func HttpGet (url string) ([]byte, error) {

    r, err := http.Get(url)
    defer r.Body.Close()

    if err != nil {
        return nil, fmt.Errorf("GET failed (%s)", url)
    }

    // read the response
    body, err := ioutil.ReadAll(r.Body)

    if err != nil {
        return nil, fmt.Errorf("Response read failed")
    }

    return body, nil
}

/*
 * Decodes a json []byte array into a populated DocumentCollection
 */
func DecodeJsonToDocCollection (b []byte) (*DocumentCollection, error) {

    var cont interface{}
    err := json.Unmarshal(b, &cont)

    if err != nil {
        return nil, fmt.Errorf("Response decode error")
    }

    response := cont.(map[string] interface{})["response"]
    // the total amount of results, irrespective of the amount returned in the response
    num_found := int(response.(map[string] interface{})["numFound"].(float64))
    // the amount we have here

    docs := response.(map[string] interface{})["docs"].([]interface{})
    num_results := len(docs)

    coll := DocumentCollection{}
    coll.NumFound = num_found

    ds := make([]Document, num_results)
    
    for i := 0; i < num_results; i++ {
        ds[i] = Document{docs[i].(map[string] interface{})}
    }

    coll.Collection = ds

    return &coll, nil
}

/*
 * Performs a Query (select)
 */
func (c *Connection) RawQuery (q string) (*DocumentCollection, error) {

    body, err := HttpGet(fmt.Sprintf("http://%s:%d/solr/select?wt=json&q=%s", c.Host, c.Port, q))

    if err != nil {
        return nil, fmt.Errorf("Some sort of http failure") // TODO: investigate how net/http fails
    }

    dc, err := DecodeJsonToDocCollection(body)
    
    if err != nil {
        return nil, err
    }

    return dc, nil
}

// func (c *Connection) Query(a ...interface{}) *DocumentCollection {

// }