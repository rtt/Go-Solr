package solr

import (
    "fmt"
)

type URLParamMap map[string] []string

/*
 * Represents a "connection"; actually just a host and port
 * (and probably at some point a Solr Core name)
 */
type Connection struct {
    Host string
    Port int
    Core string
}

/* 
 * Represents a Solr document, as returned by Select queries
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
    Start int
}

type Response struct {
    Results *DocumentCollection
    Status int
    QTime int
    // TODO: Debug info as well?
}

/*
 * Represents a Query with various params
 */
type Query struct {
    Params map[string] interface{}
    Rows int
    Start int
    Sort string
    DefType string
    Debug bool
    OmitHeader bool
}

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
 * @returns *Connection, error
 */
func Init(host string, port int) (*Connection, error) {
    
    if len(host) == 0 {
        return nil, fmt.Errorf("Invalid hostname (must be length >= 1)")
    }

    if port <= 0 || port > 65535 {
        return nil, fmt.Errorf("Invalid port (must be 1..65535")
    }

    return &Connection{Host: host, Port: port}, nil
}


/*
 * Performs a raw Select query using a given string
 */
func (c *Connection) SelectRaw (q string) (*Response, error) {

    body, err := HTTPGet(SolrString(c, q))

    if err != nil {
        return nil, fmt.Errorf("Some sort of http failure") // TODO: investigate how net/http fails
    }

    r, err := ResponseFromHTTPResponse(body)
    

    if err != nil {
        return nil, err
    }

    return r, nil
}

/*
 * Performs a Select query against the given query
 */
func (c *Connection) Select(q Query) (*Response, error) {
    return nil, nil
}

// func (c *Connection) Update(q Query) (*Response, error) {

// }