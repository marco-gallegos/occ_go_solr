package solr

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vanng822/go-solr/solr"
)

func AddParameterToQ(currentQuery string, value string) string {
	currentQuery += "titulo:" + value + "\n"
	currentQuery += "cuerpo:" + value + "\n"
	currentQuery += "localidad:" + value + "\n"
	currentQuery += "salario:" + value + "\n"
	currentQuery += "categoria:" + value + "\n"
	return currentQuery
}

func Search(w http.ResponseWriter, request *http.Request) {
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
		queryString = AddParameterToQ(queryString, titulo)
	}

	if cuerpo != "" {
		queryString = AddParameterToQ(queryString, cuerpo)
	}

	if localidad != "" {
		queryString = AddParameterToQ(queryString, localidad)
	}

	if salario != "" {
		queryString = AddParameterToQ(queryString, salario)
	}

	if categoria != "" {
		queryString = AddParameterToQ(queryString, categoria)
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

	s := SolrConection.Search(query)
	r, _ := s.Result(nil)

	fmt.Println("==========================")
	fmt.Println(r.FacetCounts)

	var data []solr.Document = r.Results.Docs
	var results []JobPosition

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

		jobPosition := JobPosition{
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

func StoreJob(w http.ResponseWriter, request *http.Request) {
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
	SolrConection.Add(updates, size, nil)
	SolrConection.Commit()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	message := "Job added"
	w.Write([]byte(message))
}
