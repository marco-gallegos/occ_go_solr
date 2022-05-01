package solr_test

import (
	"testing"

	"github.com/vanng822/go-solr/solr"

	solrpackage "occ_go_solr/solr"
)

// ===================================================================================
// basic testing on solar conection and libs
func TestSolrConnection(t *testing.T) {
	var conected bool = solrpackage.ConectSolr("localhost:8983/solr", "core")
	if conected == false {
		t.Error("Solr connection failed")
	}
}

func TestSolrQuery(t *testing.T) {
	solrpackage.ConectSolr("http://192.168.0.111:8983/solr", "jcg_example_core")

	// creating query
	query := solr.NewQuery()

	query.Q("*.*")

	s := solrpackage.SolrConection.Search(query)

	_, error := s.Result(nil)

	if error != nil {
		t.Error("Solr query failed")
	}
}

func TestAExpectedSolrQuery() {

}

// ===================================================================================
// TODO:unit testing on solr
// figure out how to mock data in go to make more complex testings on solr
// test a good result on mocked search
// test a bad result on mocked search
