package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/vanng822/go-solr/solr"
)

type jobPosition struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Salary      string `json:"salary"`
	Category    string `json:"category"`
}

func createJobList(data [][]string) []jobPosition {
	var jobList []jobPosition
	for i, line := range data {
		if i > 0 { // omit header line
			var rec jobPosition
			for j, field := range line {
				if j == 0 {
					rec.Title = field
				}
				if j == 1 {
					rec.Description = field
				}
				if j == 2 {
					rec.Location = field
				}
				if j == 3 {
					rec.Salary = field
				}
				if j == 4 {
					rec.Category = field
				}
			}
			jobList = append(jobList, rec)
		}
	}
	return jobList
}

func main() {
	// open file
	f, err := os.Open("./import.csv")
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// convert records to array of structs
	jobList := createJobList(data)

	// print the array to debug
	// fmt.Printf("%+v\n", jobList)

	// store in solr
	solrConection, _ := solr.NewSolrInterface("http://192.168.0.111:8983/solr", "jcg_example_core")
	var updates []solr.Document

	for _, job := range jobList {
		update := solr.Document{
			"titulo":    job.Title,
			"cuerpo":    job.Description,
			"localidad": job.Location,
			"salario":   job.Salary,
			"categoria": job.Category,
		}
		updates = append(updates, update)
	}
	var size int = len(updates)
	fmt.Println(size, updates)
	solrConection.Add(updates, size, nil)

	//TODO: how to avoid duplicate entries?
	solrConection.Commit()
}
