package solr

import (
	"fmt"

	"github.com/vanng822/go-solr/solr"
)

// globals
var SolrConection *solr.SolrInterface

func ConectSolr(url string, core string) bool {
	newConection, err := solr.NewSolrInterface(url, core)

	if err == nil {
		SolrConection = newConection
		return true
	} else {
		fmt.Println(err)
		return false
	}
}
