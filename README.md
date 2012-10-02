# Go-Solr

An [Apache Solr](http://lucene.apache.org/solr/) library written in Go, which is [my](http://rsty.org) first Go project! Functionality includes:

* Select queries
* Raw queries - useful for more complex queries such as [Function Queries](http://wiki.apache.org/solr/FunctionQuery)
* [Update queries](http://wiki.apache.org/solr/UpdateJSON) (add/replace/delete)
* [Faceting](http://wiki.apache.org/solr/SolrFacetingOverview)

For more information on Solr itself, please refer to [Solr's wiki](http://wiki.apache.org/solr/).

This library is released under the "free as in beer" license. Use it, mess with it, do whatever you like with it... Comments, suggestions, pull requests etc are welcomed.

## Examples / Documentation

Example programs can be found in the `rsty/solr-example` package, [here](https://github.com/rtt/Go-Solr/blob/master/rsty/solr-example/example.go).

### Creating a connection (solr.Init)

Import the `rsty/solr` package (it is assumed you know how to build and install it) and create a "connection" to your solr server.

```go
s := solr.Init("localhost", 8983)
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

Notice that `facet_field` is an array of strings and appears multiple times in the resulting query string (above)

Performing a query using our `solr.Query` is simple and shown below

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
res, err := s.RawQuery(q)
if err != nil {
    // handle error here
}
// ...
```

In other words, under the hood the following query will have been performed:

```
GET http://localhost:8983/solr/select?q={!func}add($v1,$v2)&v1=sqrt(popularity)&v2=100.0
```

## Responses (to solr.Select/solr.SelectRaw queries)

Responses to select queries (`solr.Select()` and `solr.RawSelect()`) come in the form of pointers to `SelectResponse` types. A response wraps up a solr response. The following few paragraphs and sections describe the various parts of a `SelectResponse` object

### SelectResponse object

A `SelectResponse` object and an error indicator is returned from calls to `solr.Select()` and `solr.SelectRaw()`. A `SelectResponse` mimics a Solr response and therefore has the following attributes:

* `Results` - a pointer to a `DocumentCollection` (more on this later) which contains the documents returned by Solr
* `Status` - `status` indicator as returned by Solr
* `QTime` - `QTime` value as returned by Solr

### DocumentCollection object

A `DocumentCollection` wraps up a set of `Document`s providing a convenient interface to them.

`DocumentCollection` supports the following methods:

* `Len() int` - returns the length (int) of the `Document`s returned
* `Get(i int) *Document` - returns a pointer to the document at position `i` within the Collection

`DocumentCollection` has the following properties

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

[TODO]