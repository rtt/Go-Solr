# Go-Solr

A solr library written in Go, which is my first Go project! This library will start out simple and effective as I learn the ins and outs of Go. Functionality will grow over time, hopefully, to offer a full featured interface. My aim is to be completely transparent as the library grows, so you can see all of my mistakes and learning in the repository history.

This library is released under the "free as in beer" license. Steal it, mess with it, do whatever you like with it...

## Examples / Documentation

Import the solr library `rsty/solr` and create a "connection" to your solr server (multicore support will come later!).

```go
s := solr.Init("localhost", 8983)
```

`Solr` supports just one method (at the moment) which is a raw query string ([solr's `q` parameter](http://wiki.apache.org/solr/CommonQueryParameters#q)). Raw queries are performed as follows:

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

`res`, if the `RawQuery` method call was successful, will point to a `DocumentCollection` type which holds any results (in `Document`s) returned from your query, and a few other metrics (explained shortly).

`DocumentCollection` supports the following methods:

* `Len() int` - returns the length (int) of the documents returned
* `Get(i int)` - returns a pointer to the document (Document) at position i within the Collection

`DocumentCollection` has the following properties

* `NumFound` - the total number of results solr matched to your query (irrespective of the amount returned)

`Document`s implement the following methods:

* `Field(field_name string)` - returns the value of the field, specified by `field_name`

An example program can be found in the solr-example folder, [here](https://github.com/rtt/Go-Solr/blob/master/rsty/solr-example/example.go)
