# Go-Solr

A solr library written in Go, which is my first Go project! This library will start out simple and effective as I learn the ins and outs of Go. Functionality will grow over time, hopefully, to offer a full featured interface.

## Examples

Import the solr library `rsty/solr` and create a "connection" to your solr instance.

```go
s = solr.Init("localhost", 8983)`
```

`Solr` supports just one method (at the time of writing) which is a raw query string (`q` parameter). Queries are performed as follows:

```go
q := "*:*" // a solr query
res, err := s.RawQuery(q)
if err != nil {
    // handle error
}
...
```

`res`, if the `RawQuery` call was successful, will point to a `DocumentCollection` which holds results (if any) returned from your query.

`DocumentCollection` supports the following methods:

* `Len() int` - returns the length of the documents returned
* `Get(i int) Document` - returns document at position i

`DocumentCollection` has the following properties

* NumFound - the total number of results solr matched to your query (irrespective of the amount returned)
