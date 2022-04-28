package main

import (
	"encoding/json"
	"fmt"

	//"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vanng822/go-solr/solr"
)

// globals
var solrConection *solr.SolrInterface

type jobPosition struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Salary      string `json:"salary"`
	Category    string `json:"category"`
}

func addParametertoQ(currentQuery string, value string) string {
	currentQuery += "titulo:" + value + "\n"
	currentQuery += "cuerpo:" + value + "\n"
	currentQuery += "localidad:" + value + "\n"
	currentQuery += "salario:" + value + "\n"
	currentQuery += "categoria:" + value + "\n"
	return currentQuery
}

func search(w http.ResponseWriter, request *http.Request) {
	// creating query
	query := solr.NewQuery()

	titulo := request.PostFormValue("titulo")
	cuerpo := request.PostFormValue("cuerpo")
	localidad := request.PostFormValue("localidad")
	salario := request.PostFormValue("salario")
	categoria := request.PostFormValue("categoria")

	//TODO: filter query
	var queryString string = ""

	if titulo != "" {
		queryString = addParametertoQ(queryString, titulo)
	}

	if cuerpo != "" {
		queryString = addParametertoQ(queryString, cuerpo)
	}

	if localidad != "" {
		queryString = addParametertoQ(queryString, localidad)
	}

	if salario != "" {
		queryString = addParametertoQ(queryString, salario)
	}

	if categoria != "" {
		queryString = addParametertoQ(queryString, categoria)
	}

	if queryString == "" {
		queryString = "*:*"
	}

	//fmt.Println(queryString)

	query.Q(queryString)
	query.AddParam("q.op", "OR")

	// configuring some facets
	query.AddParam("facet", "true")
	query.AddParam("facet.field", "categoria")
	query.AddParam("facet.field", "salario")

	fmt.Println(query)

	s := solrConection.Search(query)
	r, _ := s.Result(nil)

	fmt.Println("==========================")
	fmt.Println(r.FacetCounts)

	var data []solr.Document = r.Results.Docs
	var results []jobPosition

	for _, doc := range data {
		titulo := doc.Get("titulo")
		if titulo == nil {
			titulo = "N/A"
		}

		descripcion := doc.Get("cuerpo")
		if descripcion == nil {
			descripcion = "N/A"
		}

		location := doc.Get("localidad")
		if location == nil {
			location = "N/A"
		}

		salary := doc.Get("salario")
		if salary == nil {
			salary = "N/A"
		}

		category := doc.Get("categoria")
		if category == nil {
			category = "N/A"
		}

		jobPosition := jobPosition{
			Id:          doc.Get("id").(string),
			Title:       fmt.Sprintf("%v", titulo),
			Description: fmt.Sprintf("%v", descripcion),
			Location:    fmt.Sprintf("%v", location),
			Salary:      fmt.Sprintf("%v", salary),
			Category:    fmt.Sprintf("%v", category),
		}
		results = append(results, jobPosition)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(results)
}

func storeJob(w http.ResponseWriter, request *http.Request) {
	titulo := request.PostFormValue("titulo")
	cuerpo := request.PostFormValue("cuerpo")
	localidad := request.PostFormValue("localidad")
	salario := request.PostFormValue("salario")
	categoria := request.PostFormValue("categoria")

	if titulo == "" || cuerpo == "" || localidad == "" || salario == "" || categoria == "" {
		w.WriteHeader(http.StatusBadRequest)
		message := "Missing parameters"
		w.Write([]byte(message))
		return
	}

	update := solr.Document{
		"titulo":    titulo,
		"cuerpo":    cuerpo,
		"localidad": localidad,
		"salario":   salario,
		"categoria": categoria,
	}
	var updates []solr.Document = []solr.Document{update}
	var size int = 1
	solrConection.Add(updates, size, nil)
	solrConection.Commit()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	message := "Job added"
	w.Write([]byte(message))
}

func main() {
	// this code gets the first rows from the solr index
	solrConection, _ = solr.NewSolrInterface("http://192.168.0.111:8983/solr", "jcg_example_core")

	router := mux.NewRouter()
	router.HandleFunc("/search", search).Methods("POST")
	router.HandleFunc("/job", storeJob).Methods("POST")
	http.ListenAndServe(":8080", router)
	fmt.Println("Server started on port 8080")
}
