# Go-Solr

An [Apache Solr](http://lucene.apache.org/solr/) library written in [Go](http://golang.org/), which is [my](http://rsty.org) first Go project! Functionality includes:

* Select queries
* Raw Select tqueries - useful for more complex queries such as [Function Queries](http://wiki.apache.org/solr/FunctionQuery)
* [Update queries](http://wiki.apache.org/solr/UpdateJSON) (add/replace/delete)
* [Faceting](http://wiki.apache.org/solr/SolrFacetingOverview)

For more information on Solr itself, please refer to [Solr's wiki](http://wiki.apache.org/solr/).

This library is released under the "do whatever you like" license.

## Examples / Documentation

Example programs can be found in the `examples` folder [here](https://github.com/rtt/Go-Solr/tree/master/examples).

### Creating a connection - solr.Init()

Import the `solr` package (it is assumed you know how to build/install it, if not, [see here](http://golang.org/doc/install#install)) and create a "connection" to your solr server by calling the `solr.Init(hostname, port int)` function supplying a hostname and port.

```go
// connect to server running on localhost port 8983
s, err := solr.Init("localhost", 8983)
```

### Performing Select Queries - solr.Select()

Select queries are performed using the `solr.Select(q *Query)` method, passing it a pointer to a `Query` type.

Here's an example:

```go
q := solr.Query{
    Params: solr.URLParamMap{
        "q": []string{"id:31"},
        "facet.field": []string{"some_field", "some_other_field"},
        "facet": "true",
    },
    Rows: 10,
    Sort: "title ASC"
}
```

Here we have defined a set of URL parameters - `q`, `facet.field`, `facet`, `rows` and `sort` using the `solr.Query{}` struct. Under the hood this would work out as the following Solr query string:

```
GET http://localhost:8983/solr/select?q=id:31&facet.field=some_field&facet.field=some_other_field&facet=true
```

Notice that `facet_field`, like `q`, is an array of strings and appears multiple times in the resulting query string (shown above)

Performing a query using our `solr.Query()` is simple and shown below

```go
res, err := s.Select(&q)
if err != nil {
  fmt.Println(err)
}

// ...
```

A pointer to `q` is passed to `s.Select()`, and returned is a pointer to a `SelectResponse` (`res`) and an `error` (`err`) if an error occurred.

Iterating over the results is shown later in this document.

### Performing 'Raw' Select Queries - solr.SelectRaw()

`rsty/solr` supports raw queries where you can specify your exact query in string form. This is useful for specifying complex queries where a `Query` type would be cumbersome. Raw queries are performed as follows:

```go
q := "q={!func}add($v1,$v2)&v1=sqrt(popularity)&v2=100.0" // a solr query
res, err := s.SelectRaw(q)
if err != nil {
    // handle error here
}
// ...
```

In other words, under the hood the following query will have been performed:

```
GET http://localhost:8983/solr/select?q={!func}add($v1,$v2)&v1=sqrt(popularity)&v2=100.0
```
As with `solr.Select()`, `solr.SelectRaw()` returns a pointer to a `SelectResponse` and an error, `err`.

## Responses - SelectResponse type

Responses to select queries (`solr.Select()` and `solr.RawSelect()`) come in the form of pointers to `SelectResponse` types. A `SelectResponse` wraps a Solr response with a convenient interface. The following few paragraphs and sections describe the various parts of a `SelectResponse` object

### SelectResponse type

A pointer to a `SelectResponse` and an error are returned from calls to `solr.Select()` and `solr.SelectRaw()`. A `SelectResponse` mimics a Solr response and therefore has the following attributes:

* `Results` - a pointer to a `DocumentCollection` (more on this later) which contains the documents returned by Solr
* `Status` - query `status` indicator as returned by Solr
* `QTime` - `QTime` value as returned by Solr

More information on `Status` and `QTime` can be found [here](http://wiki.apache.org/solr/SolrTerminology).

### DocumentCollection object

A `DocumentCollection` wraps up a set of `Document`s providing a convenient interface to them.

`DocumentCollection` supports the following methods:

* `Len() int` - returns the length (int) of the `Document`s returned
* `Get(i int) *Document` - returns a pointer to the document at position `i` within the Collection

`DocumentCollection` has the following properties:

* `NumFound` - the total number of results solr matched to your query (irrespective of the amount returned)
* `Facets` - an array of `Facet` objects
* `NumFacets` - the number of facet fields returned (if any)

### Document object

`Document`s implement the following methods:

* `Field(field_name string) interface{}` - returns the value of the field, specified by `field_name`

## Faceting

### Facet object

If your select query specifies facets, facets will be found under `Response.Results.Facets` which is an array of `Facet`s. A `Facet` has the following attributes

* `Name` - the name of the facet (field) as returned by Solr
* `Counts` - an array of `FacetCount`s, the corresponding value counts for the field.

### FacetCount object

A `FacetCount` has the following attributes

* `Value` - the facet field value
* `Count` - the count (int) for the field value

### Faceting example

Below is an example showing an iteration over a collection of `Facet`s

```go
// q is assumed to have been set up

// perform the query
res, err := s.Query(&q)

// handle error, err, here

results := res.Results
for i := 0; i < results.NumFacets; i++ {
    facet := results.Facets[i]
    fmt.Println("Facet:", facet.Name)
    k := len(facet.Counts)
    for j := 0; j < k; j++ {
        fmt.Println(facet.Counts[j].Value, "=>", facet.Counts[j].Count)
    }
    fmt.Println("")
}
```

This might output the following:

```bash
Facets
------
Facet: category
cameras => 1

Facet: type
digital_slr => 10
compact => 2
```

## Update Queries - solr.Update()

Update queries are used to add, replace or delete documents in Solr's index. Please see the [Solr Wiki](http://wiki.apache.org/solr/UpdateJSON) for more information.  Go-Solr uses JSON for update queries, *not* XML. Solr3.1 will need to be [configured](http://wiki.apache.org/solr/UpdateJSON#Requirements) to support JSON for update messages, Solr 4.0+ supports JSON natively via `/update`.

### Creating an Update query - example

`solr.Update(document map[string]interface{}, commit bool)` takes two arguments, an "update document" and a commit flag (boolean) which specifies whether or not a commit should be performed at the same time as the update is performed. An example may look like the following

```go
q, err := solr.Update(document, true);
if err != nil {
    // ...
}
```

An update document must be of type `map[string]interface{}`, and may look like the following:

```go
doc := map[string]interface {}{
    "add":[]interface {}{
        map[string]interface {}{"id": 22, "title": "abc"},
        map[string]interface {}{"id": 23, "title": "def"},
        map[string]interface {}{"id": 24, "title": "def"},
    },
}
```
... which is equivalent to the following JSON:

```json
{"add": [{"id": 22, "title": "abc"}, {"id": 23, "title": "def"}, {"id": 24, "title": "def"}]}
```

... which is an Update which adds (or replaces) 3 documents in a fictional Solr index.

You can define any type of document to send off to Solr in an update. Support will be added later to allow raw JSON strings to be used in Updates.

`solr.Update()` returns an `UpdateResponse` and an `error`. `UpdateResponse` has a `Success` (bool) property.
