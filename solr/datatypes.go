package solr

//public
type JobPosition struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Salary      string `json:"salary"`
	Category    string `json:"category"`
}
