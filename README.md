# Go-Solr

A solr library written in Go, which is [my](http://rsty.org) first Go project! This library will start out simple and effective as I learn the ins and outs of Go. Functionality will grow over time, hopefully, to offer a full featured interface. My aim is to be completely transparent as the library is developed so you can see all of my mistakes and learning in the repository history.

This library is released under the "free as in beer" license. Steal it, mess with it, do whatever you like with it...

## Examples / Documentation

Note: Please see the [examples](https://github.com/rtt/Go-Solr/blob/master/rsty/solr-example) directory for buildable examples.

### Creating a connection (solr.Init)

Import the `rsty/solr` package (it is assumed you know how to build and install it) and create a "connection" to your solr server (multicore support will come later!).

```go
s := solr.Init("localhost", 8983)
```

### Performing Select Queries - solr.Select()

Select queries are performed using the `solr.Select(q *Query)` method, passing it a pointer to a `Query` struct.

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

Here we have defined a set of URL parameters - `q`, `facet.field`, `facet`, `rows` and `sort` using the `solr.Query{}` struct. This would work out as the following Solr query string:

```
GET http://localhost:8983/solr/select?q=id:31&facet.field=some_field&facet.field=some_other_field&facet=true
```

Notice that `facet_field` is a slice of strings and is doubled-up in the resulting query string.

Performing a query using our `solr.Query` is shown below

```go
res, err := s.Select(&q)
if err != nil {
  fmt.Println(err)
}

// ...
```

A pointer to `q` is passed to `s.Select`, and returned is a `Response` (`res`) and an `error` (`err`).


### Performing 'Raw' Select Queries - solr.SelectRaw

`rsty/solr` supports raw queries where you, as an advanced Solr user, can specify your exact query in string form. Raw queries are performed as follows:

```go
q := "*:*" // a solr query
res, err := s.RawQuery(q)
if err != nil {
    // handle error here
}
...
```

In other words, this will have performed the following HTTP query:

```
GET http://localhost:8983/solr/select?q=*:*
```

## Responses (to solr.Select/solr.SelectRaw queries)

### Responses

A `Response` object and an error indicator is returned from calls to `solr.Select` and `solr.SelectRaw`. A `Response` has the following attributes:

* `Results` - a pointer to a `DocumentCollection` (more on this later) which contains the documents returned by Solr
* `Status` - `status` indicator as returned by Solr
* `QTime` - `QTime` value as returned by Solr


### DocumentCollections

A `DocumentCollection` wraps up a set of `Document`s providing a convenient interface to them.

`DocumentCollection` supports the following methods:

* `Len() int` - returns the length (int) of the `Document`s returned
* `Get(i int) *Document` - returns a pointer to the document at position `i` within the Collection

`DocumentCollection` has the following properties

* `NumFound` - the total number of results solr matched to your query (irrespective of the amount returned)

### Documents

`Document`s implement the following methods:

* `Field(field_name string) interface{}` - returns the value of the field, specified by `field_name`

An example program can be found in the solr-example package, [here](https://github.com/rtt/Go-Solr/blob/master/rsty/solr-example/example.go)
