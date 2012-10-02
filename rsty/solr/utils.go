package solr

import (
    "bytes"
    "fmt"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "strings"
)


/*
 * Performs a GET request to the given url
 * Returns a []byte containing the response body
 */
func HTTPGet (url string) ([]byte, error) {

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

func HTTPPost (url string, headers [][]string, payload []byte) ([]byte, error) {
    // setup post client
    client := &http.Client{}
    req, err := http.NewRequest("POST", url, bytes.NewReader(payload))

    // add headers
    if len(headers) > 0 {
        for i := range headers {
            req.Header.Add(headers[i][0], headers[i][1])
        }
    }

    // perform request
    resp, err := client.Do(req)
    defer resp.Body.Close()

    if err != nil {
        return nil, fmt.Errorf(fmt.Sprintf("POST request failed: %s", err))
    }

    // read response & return
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    return body, nil
}


/*
 * Returns a URLEncoded version of a Param Map
 * E.g., ParamMap[foo:bar omg:wtf] => "foo=bar&omg=wtf"
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
func SolrSelectString (c *Connection, q string) string {
    return fmt.Sprintf("http://%s:%d/solr/select?wt=json&%s", c.Host, c.Port, q)
}

/*
 * Generates a Solr update query string
 */
func SolrUpdateString (c *Connection, commit bool) string {
    s := fmt.Sprintf("http://%s:%d/solr/update", c.Host, c.Port)
    if commit {
        return fmt.Sprintf("%s?commit=true", s)
    }
    return s
}


/*
 * Decodes a json formatted []byte into an interface{} type
 */
func BytesToJSON (b *[]byte) (*interface{}, error) {

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
func JSONToBytes (m map[string] interface{}) (*[]byte, error) {
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
func BuildResponse (j *interface{}) (*SelectResponse, error) {

    // look for a response element, bail if not present
    response_root := (*j).(map[string] interface{})
    response := response_root["response"]
    if response == nil {
        return nil, fmt.Errorf("Supplied interface appears invalid (missing response)")
    }

    // begin Response creation
    r := SelectResponse{}

    // do status & qtime, if possible
    r_header := (*j).(map[string] interface{})["responseHeader"].(map[string] interface{})
    if r_header != nil {
        r.Status = int(r_header["status"].(float64))
        r.QTime = int(r_header["QTime"].(float64))
    }

    // now do docs, if they exist in the response
    docs := response.(map[string] interface{})["docs"].([]interface{})
    if docs != nil {
        // the total amount of results, irrespective of the amount returned in the response
        num_found := int(response.(map[string] interface{})["numFound"].(float64))

        // and the amount actually returned
        num_results := len(docs)

        coll := DocumentCollection{}
        coll.NumFound = num_found

        ds := []Document{}

        for i := 0; i < num_results; i++ {
            ds = append(ds, Document{docs[i].(map[string] interface{})})
        }

        coll.Collection = ds
        r.Results = &coll
    }

    // facets
    facet_counts := response_root["facet_counts"].(map[string] interface{})
    if facet_counts != nil {
        // do counts if they exist
        facet_fields := facet_counts["facet_fields"].(map[string] interface{})
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

    return &r, nil
}


/*
 * Decodes a HTTP (Solr) response and returns a Response
 */
func SelectResponseFromHTTPResponse (b []byte) (*SelectResponse, error) {
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
 * is an error response or not
 */
func SolrErrorResponse(m map[string] interface{}) (bool, *ErrorResponse) {
    // check for existance of "error" key
    if _, found := m["error"]; found {
        error := m["error"].(map[string] interface{})
        return true, &ErrorResponse{
            Message: error["msg"].(string),
            Status: int(error["code"].(float64)),
        }
    }
    return false, nil
}

/*
 * Similar to python's itertools.izip_longest;
 * takes an array and chunks it according to a given size
 */
func chunk(s []interface{}, sz int) [][]interface{} {
  r := [][]interface{}{}
  j := len(s)
  for i := 0; i < j; i+=sz {
    r = append(r, s[i:i+sz])
  }
  return r
}
