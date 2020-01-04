package api

import (
	"encoding/json"
	"strconv"
	"errors"
	"fmt"
	"log"
	"net/http"
	"unicode"
	"mibolsillo/pkg/tools"
	"github.com/golang/gddo/httputil/header"
)

//CreateInvoice Handler for message
func CreateInvoice(w http.ResponseWriter, r *http.Request) {
	
	//Decoding , it will return a map of IDs and respond
	mapResponse , err := decodeJSONBody(w, r)

	if err != nil {
		var mr *tool.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	arrayResponse := make([]*Response , 0)
	
	//ADD ids to Responses to returne it
	for _ , res := range mapResponse{
		arrayResponse = append(arrayResponse, res)
	}
	//responses := Responses{Responses: arrayResponse}
	json := json.NewEncoder(w)
	//Return response
	json.Encode(arrayResponse)
	//Check last token
}
//
func decodeJSONBody(w http.ResponseWriter, r *http.Request) (map[string]*Response ,error) {


	//Check Header , if its not json return err
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			Msg := "Content-Type header is not application/json"
			return nil , &tool.MalformedRequest{Status: http.StatusUnsupportedMediaType, Msg: Msg}
		}
	}
	//Max file size 10 mb
	r.Body = http.MaxBytesReader(w, r.Body, 8 * 1024 * 1024 * 10)
	
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	
	dec.Token()


	idMaps := make(map[string]*Response)

	for dec.More() {
		mess := &Message{}
		if err := dec.Decode(mess); err != nil {
			return nil , tool.ErrorHandling(err)
		}
		if err := transformRespond(mess, idMaps); err != nil {
			return nil , tool.ErrorHandling(err)
		}
	}

	dec.Token()

	return idMaps , nil
}

func isInt(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

//DecodeJSONBody uncion para filtrar los errores de Json decoding
func parsingValue(number interface{}) (float64, error) {
	switch n := number.(type) {

	case float64:
		return n, nil

	case string:
		if isInt(n) {
			return strconv.ParseFloat(n , 64) 
		}
	}
	return 0, fmt.Errorf("Error parsing value have an undefined type: %v", number)
}

func transformRespond(obj *Message, mapa map[string]*Response) error {
	//value , err := fromInterfactToInt(newObj.Value)
	var income float64
	var expenses float64
	var err error
	
	re , ok := mapa[obj.ID]
	if !ok {
		re = &Response{User: obj.ID}
		mapa[obj.ID] = re
	}

	switch {
	case obj.Type == "income":
		income, err = parsingValue(obj.Value)
	case obj.Type == "expense":
		expenses, err = parsingValue(obj.Value)
	default:
		return fmt.Errorf("Error parsing type es incorrect: %s", obj.Type)
	}
	if err != nil {
		return err
	}

	re.Expenses += expenses
	re.Revenue += income
	re.Transactions = append(re.Transactions, obj)
	return nil
}

