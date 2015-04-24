package solr

import (
	"testing"
)

func TestInitOK(t *testing.T) {

	c, _ := Init("localhost", 8696, "core0")
	if c.URL != "http://localhost:8696/solr/core0" {
		t.Fail()
	}

}

func TestInitInvalidHost(t *testing.T) {

	_, err := Init("", 700000, "core0")
	if err == nil {
		t.Fail()
	}
}

func TestInitInvalidPort(t *testing.T) {

	_, err := Init("localhost", 700000, "core0")
	if err == nil {
		t.Fail()
	}
}

func TestSolrSelectString(t *testing.T) {
	c, _ := Init("localhost", 8696, "core0")
	q := &Query{
		Params: URLParamMap{
			"q": []string{"id:1"},
		},
	}
	s := SolrSelectString(c, q.String(), "select")
	if s != "http://localhost:8696/solr/core0/select?wt=json&q=id%3A1" {
		t.Fail()
	}

}

func TestSolrUpdateString(t *testing.T) {
	c, _ := Init("localhost", 8696, "core0")
	s := SolrUpdateString(c, true)
	if s != "http://localhost:8696/solr/core0/update?commit=true" {
		t.Fail()
	}

}
