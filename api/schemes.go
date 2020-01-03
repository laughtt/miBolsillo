package api


type Responses struct{
	Responses []*Response `json:""`
}

type Response struct {
	Expenses float64  `json:"expenses"`
	Revenue float64  `json:"revenue"`
	User string `json:"user"`
	Transactions []*Message `json:"transactions"`
}

type Message struct {
	Value       interface{} `json:"value"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
	ID          string  `json:"id"`
}

//MalformedRequest Error requerido
type MalformedRequest struct {
	Status int
	Msg    string
}
