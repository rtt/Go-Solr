/*
 * Go Solr, a Solr library written in Go.
 * Original author Rich Taylor, 2012 - http://rsty.org/, http://github.com/rtt/
 * Released under the "do whatever the fuck you want" license. http://sam.zoy.org/wtfpl/
 */

package solr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

/*
 * Represents a "connection"; actually just a host and port
 * (and probably at some point a Solr Core name)
 */
type Connection struct {
	Host    string
	Port    int
	Core    string
	Version []int
}

/*
 * Represents a Solr document, as returned by Select queries
 */
type Document struct {
	Fields map[string]interface{}
}

/*
 * Represents a FacetCount for a Facet
 */
type FacetCount struct {
	Value string
	Count int
}

/* chunked size of facet solr return format */
var facet_chunk_size int = 2

/*
 * Represents a Facet with a name and count
 */
type Facet struct {
	Name   string       // accepts_4x4s
	Counts []FacetCount // a set of values
}

/*
 * Represents a collection of solr documents
 * and various other metrics
 */
type DocumentCollection struct {
	Facets     []Facet
	Collection []Document
	NumFacets  int // convenience...
	NumFound   int
	Start      int
}

/*
 * Represents a Solr response
 */
type SelectResponse struct {
	Results *DocumentCollection
	Status  int
	QTime   int
	// TODO: Debug info as well?
}

/*
 * Represents an error from Solr
 */
type ErrorResponse struct {
	Message string
	Status  int
}

type UpdateResponse struct {
	Success bool
}

/*
 * Holds URL parameters
 */
type URLParamMap map[string][]string

/*
 * Query represents a query with various params
 */
type Query struct {
	Params     URLParamMap
	Rows       int
	Start      int
	Sort       string
	DefType    string
	Debug      bool
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
func (document Document) Doc() map[string]interface{} {
	return document.Fields
}

func (r SelectResponse) String() string {
	return fmt.Sprintf("SelectResponse: %d Results, Status: %d, QTime: %d", r.Results.Len(), r.Status, r.QTime)
}

func (r ErrorResponse) String() string {
	return fmt.Sprintf("Solr Error: [code: %d, msg: \"%s\"]", r.Status, r.Message)
}

func (r UpdateResponse) String() string {
	if r.Success {
		return fmt.Sprintf("UpdateResponse: OK")
	}
	return fmt.Sprintf("UpdateResponse: FAIL")
}

/*
 * Performs a GET request to the given url
 * Returns a []byte containing the response body
 */
func HTTPGet(url string) ([]byte, error) {

	r, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("GET failed (%s)", url)
	}

	// read the response and check
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("Response read failed")
	}

	return body, nil
}

/*
 * Performs a HTTP Post request. Takes:
 * * A url
 * * Headers, in the format [][]string{} (e.g., [[key, val], [key, val], ...])
 * * A payload (post request body) which can be nil
 * * Returns the body of the response and an error if necessary
 */
func HTTPPost(url string, headers [][]string, payload *[]byte) ([]byte, error) {
	// setup post client
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader(*payload))

	// add headers
	if len(headers) > 0 {
		for i := range headers {
			req.Header.Add(headers[i][0], headers[i][1])
		}
	}

	// perform request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("POST request failed: %s", err))
	}

	// read response, check & return
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

/*
 * Returns a URLEncoded version of a Param Map
 * E.g., ParamMap[foo:bar omg:wtf] => "foo=bar&omg=wtf"
 * TODO: This isn't exactly safe and there's probably a library pkg to do this already...
 */
func EncodeURLParamMap(m *URLParamMap) string {
	r := []string{}

	for k, v := range *m {
		l := len(v)
		for x := 0; x < l; x++ {
			r = append(r, fmt.Sprintf("%s=%s", k, v[x]))
		}
	}

	return strings.Join(r, "&")
}

/*
 * Generates a Solr query string from a connection and a query string
 */
func SolrSelectString(c *Connection, q string) string {
	fmt.Sprintf("http://%s:%d/solr/%s/select?wt=json&%s", c.Host, c.Port,c.Core, q)
	return fmt.Sprintf("http://%s:%d/solr/%s/select?wt=json&%s", c.Host, c.Port,c.Core, q)
}

/*
 * Generates a Solr update query string. Adds ?commit=true
 * if commit arg is true.
 */
func SolrUpdateString(c *Connection, commit bool) string {
	s := fmt.Sprintf("http://%s:%d/solr/update", c.Host, c.Port)
	if commit {
		return fmt.Sprintf("%s?commit=true", s)
	}
	return s
}

/*
 * Decodes a json formatted []byte into an interface{} type
 */
func BytesToJSON(b *[]byte) (*interface{}, error) {
	var container interface{}
	err := json.Unmarshal(*b, &container)

	if err != nil {
		return nil, fmt.Errorf("Response decode error")
	}

	return &container, nil
}

/*
 * Encodes a map[string]interface{} to bytes and returns
 * a pointer to said bytes
 */
func JSONToBytes(m map[string]interface{}) (*[]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("Failed to encode JSON")
	}

	return &b, nil
}

/*
 * Takes a JSON formatted Solr response (interface{}, not []byte)
 * And returns a *Response
 */
func BuildResponse(j *interface{}) (*SelectResponse, error) {

	// look for a response element, bail if not present
	response_root := (*j).(map[string]interface{})
	response := response_root["response"]
	if response == nil {
		return nil, fmt.Errorf("Supplied interface appears invalid (missing response)")
	}

	// begin Response creation
	r := SelectResponse{}

	// do status & qtime, if possible
	r_header := (*j).(map[string]interface{})["responseHeader"].(map[string]interface{})
	if r_header != nil {
		r.Status = int(r_header["status"].(float64))
		r.QTime = int(r_header["QTime"].(float64))
	}

	// now do docs, if they exist in the response
	docs := response.(map[string]interface{})["docs"].([]interface{})
	if docs != nil {
		// the total amount of results, irrespective of the amount returned in the response
		num_found := int(response.(map[string]interface{})["numFound"].(float64))

		// and the amount actually returned
		num_results := len(docs)

		coll := DocumentCollection{}
		coll.NumFound = num_found

		ds := []Document{}

		for i := 0; i < num_results; i++ {
			ds = append(ds, Document{docs[i].(map[string]interface{})})
		}

		coll.Collection = ds
		r.Results = &coll
	}

	// facets
	facet_response, ok := response_root["facet_counts"].(interface{})
	if ok == true {
		facet_counts := facet_response.(map[string]interface{})
		if facet_counts != nil {
			// do counts if they exist
			facet_fields := facet_counts["facet_fields"].(map[string]interface{})
			facets := []Facet{}
			if facet_fields != nil {
				// iterate over each facet field, create facet & counts for each field
				for k, v := range facet_fields {
					f := Facet{Name: k}
					chunked := chunk(v.([]interface{}), facet_chunk_size)
					lc := len(chunked)
					for i := 0; i < lc; i++ {
						f.Counts = append(f.Counts, FacetCount{
							Value: chunked[i][0].(string),
							Count: int(chunked[i][1].(float64)),
						})
					}
					facets = append(facets, f)
				}
			}

			// add Facets to collection
			r.Results.Facets = facets
			r.Results.NumFacets = len(facets)

		}
	}

	return &r, nil
}

/*
 * Decodes a HTTP (Solr) response and returns a Response
 */
func SelectResponseFromHTTPResponse(b []byte) (*SelectResponse, error) {
	j, err := BytesToJSON(&b)
	if err != nil {
		return nil, fmt.Errorf("Unable to decode")
	}

	resp, err := BuildResponse(j)
	if err != nil {
		return nil, fmt.Errorf("Error building response")
	}

	return resp, nil
}

/*
 * Determines whether a decoded response from Solr
 * is an error response or not. Returns a bool (true if error)
 * and an ErrorResponse (if the response is an error response)
 * otherwise nil
 */
func SolrErrorResponse(m map[string]interface{}) (bool, *ErrorResponse) {
	// check for existance of "error" key
	if _, found := m["error"]; found {
		error := m["error"].(map[string]interface{})
		return true, &ErrorResponse{
			Message: error["msg"].(string),
			Status:  int(error["code"].(float64)),
		}
	}
	return false, nil
}

/*
 * Similar to python's itertools.izip_longest;
 * takes an array and chunks it according to a given splice size
 * eg: chnunk([1,2,3,4,5,6], 2) == [[1,2], [3,4], [5,6]]
 */
func chunk(s []interface{}, sz int) [][]interface{} {
	r := [][]interface{}{}
	j := len(s)
	for i := 0; i < j; i += sz {
		r = append(r, s[i:i+sz])
	}
	return r
}

/*
 * Inits a new Connection to a Solr instance
 * Note: this doesn't actually hold a connection, its just
 *       a container for holding a hostname & port
 */
func Init(host string, port int,core string) (*Connection, error) {

	if len(host) == 0 {
		return nil, fmt.Errorf("Invalid hostname (must be length >= 1)")
	}

	if port <= 0 || port > 65535 {
		return nil, fmt.Errorf("Invalid port (must be 1..65535")
	}

	return &Connection{Host: host, Port: port, Core: core}, nil
}

/*
 * Performs a Select query given a Query
 */
func (c *Connection) Select(q *Query) (*SelectResponse, error) {
	body, err := HTTPGet(SolrSelectString(c, q.String()))

	if err != nil {
		return nil, fmt.Errorf("Some sort of http failure") // TODO: investigate how net/http fails
	}

	r, err := SelectResponseFromHTTPResponse(body)

	if err != nil {
		return nil, err
	}

	return r, nil
}

/*
 * Performs a raw Select query given a raw query string
 */
func (c *Connection) SelectRaw(q string) (*SelectResponse, error) {
	body, err := HTTPGet(SolrSelectString(c, q))

	if err != nil {
		return nil, fmt.Errorf("Some sort of http failure") // TODO: investigate how net/http fails
	}

	r, err := SelectResponseFromHTTPResponse(body)

	if err != nil {
		return nil, err
	}

	return r, nil
}

/*
 * Performs a Solr Update query against a given update document
 * specified in a map[string]interface{} type
 * NOTE: Requires JSON updates to be enabled, see;
 * http://wiki.apache.org/solr/UpdateJSON
 * FUTURE: Will ask for solr version details in Connection and
 * act appropriately
 */
func (c *Connection) Update(m map[string]interface{}, commit bool) (*UpdateResponse, error) {

	// encode "json" to a byte array & check
	payload, err := JSONToBytes(m)
	if err != nil {
		return nil, err
	}

	// perform request
	resp, err := HTTPPost(
		SolrUpdateString(c, commit),
		[][]string{{"Content-Type", "application/json"}},
		payload)

	if err != nil {
		return nil, err
	}

	// decode the response & check
	decoded, err := BytesToJSON(&resp)
	if err != nil {
		return nil, err
	}

	error, report := SolrErrorResponse((*decoded).(map[string]interface{}))
	if error {
		return nil, fmt.Errorf(fmt.Sprintf("%s", *report))
	}

	return &UpdateResponse{true}, nil
}

// func (c *Connection) Commit() (*UpdateResponse, error) {

// }
