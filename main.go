package main

import (
	"fmt"
	"net/http"
	"occ_go_solr/solr"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func LoadEnvs() (string, string) {
	var solrUrl string
	var solrCore string
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// even if .env file is not found, we can still use the environment variables from the system/docker
	solrUrl = os.Getenv("SOLR_URL")
	solrCore = os.Getenv("SOLR_CORE")

	// if looks empty we use default values
	if solrUrl == "" || solrCore == "" {
		solrUrl = "http://192.168.0.111:8983/solr"
		solrCore = "jcg_example_core"
	}

	return solrUrl, solrCore
}

/**
In this file we only make load logic.
*/
func main() {
	var solrUrl string
	var solrCore string

	solrUrl, solrCore = LoadEnvs()

	fmt.Printf("SolrUrl: %s, Solr Core : %s\n", solrUrl, solrCore)

	solr.ConectSolr(solrUrl, solrCore)

	router := mux.NewRouter()
	router.HandleFunc("/", sayHello).Methods("GET")
	router.HandleFunc("/search", solr.Search).Methods("POST")
	router.HandleFunc("/job", solr.StoreJob).Methods("POST")

	fmt.Println("Server started on port 8080")
	error := http.ListenAndServe(":8080", router)
	if error != nil {
		fmt.Println(error)
	}
}
