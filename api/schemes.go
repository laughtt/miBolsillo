package api


const (
	//10 mb
	limitSizeFile = 10485760
)
//Responses array for the output
type Responses struct {
	Responses []*Response `json:""`
}

//Response struct to parse the output
type Response struct {
	Expenses     float64    `json:"expenses"`
	Revenue      float64    `json:"revenue"`
	User         string     `json:"user"`
	Transactions []*Message `json:"transactions"`
}

//Message struct to parse the input
type Message struct {
	Value       interface{} `json:"value"`
	Type        string      `json:"type"`
	Description string      `json:"description"`
	ID          string      `json:"id"`
}