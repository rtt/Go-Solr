package solr

import (
    "fmt"
    "strings"
)


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
 * Represents a FacetCount for a Facet
 */
type FacetCount struct {
    Key string
    Count int
}

var facet_chunk_size int = 2
/*
 * Represents a Facet with a name and count
 */
type Facet struct {
    Name string         // accepts_4x4s
    Counts []FacetCount // a set of values
}


/*
 * Represents a collection of solr documents
 * and various other metrics
 */
type DocumentCollection struct {
    Facets []Facet
    Collection []Document
    NumFacets int // convenience...
    NumFound int
    Start int
}


/*
 * Represents a Solr response
 */
type Response struct {
    Results *DocumentCollection
    Status int
    QTime int
    // TODO: Debug info as well?
}

type URLParamMap map[string] []string

/*
 * Query represents a query with various params
 */
type Query struct {
    Params URLParamMap
    Rows int
    Start int
    Sort string
    DefType string
    Debug bool
    OmitHeader bool
}


/*
 * Query.String() returns the Query in solr query string format
 */
func (q *Query) String() string {
    // TODO: this is kinda ugly
    s := []string{}

    if len(q.Params) > 0 {
        s = append(s, EncodeURLParamMap(&q.Params))
    }

    if q.Rows != 0 {
        s = append(s, fmt.Sprintf("rows=%d", q.Rows))
    }

    if q.Start != 0 {
        s = append(s, fmt.Sprintf("start=%d", q.Start))
    }

    if q.Sort != "" {
        s = append(s, fmt.Sprintf("sort=%s", q.Sort))
    }

    if q.DefType != "" {
        s = append(s, fmt.Sprintf("deftype=%s", q.DefType))
    }

    if q.Debug {
        s = append(s, fmt.Sprintf("debugQuery=true"))
    }

    if q.OmitHeader {
        s = append(s, fmt.Sprintf("omitHeader=true"))
    }

    return strings.Join(s, "&")
}


/*
 * DocumentCollection.Get() returns the document in the collection
 * at position i
 */
func (d *DocumentCollection) Get(i int) *Document {
    return &d.Collection[i]
}

/*
 * DocumentCollection.Len() returns the amount of documents
 * in the collection
 */
func (d *DocumentCollection) Len() int {
    return len(d.Collection)
}

/*
 * Document.Field() returns the value of the given field name in the document
 */
func (document Document) Field(field string) interface{} {
    r, _ := document.Fields[field]
    return r
}

/*
 * Document.Doc() returns the raw document (map)
 */
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
 * Performs a Select query given a Query
 */
func (c *Connection) Select (q *Query) (*Response, error) {
    body, err := HTTPGet(SolrString(c, q.String()))

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
 * Performs a raw Select query given a raw query string
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


// func (c *Connection) Update(q Query) (*Response, error) {

// }
