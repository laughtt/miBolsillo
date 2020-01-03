package api

import (
	"encoding/json"
	"strconv"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"
	"unicode"
	"github.com/golang/gddo/httputil/header"
)

func (mr *MalformedRequest) Error() string {
	return mr.Msg
}
 


func CreateInvoice(w http.ResponseWriter, r *http.Request) {
	var p Message


	err := decodeJSONBody(w, r, &p)
	if err != nil {
		var mr *MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {

	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			Msg := "Content-Type header is not application/json"
			return &MalformedRequest{Status: http.StatusUnsupportedMediaType, Msg: Msg}
		}
	}
	//Max file size
	r.Body = http.MaxBytesReader(w, r.Body, 8 * 1024 * 1024 * 1024)
	
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	dec.Token()
	idMaps := make(map[string]*Response)
	for dec.More() {
		newdst := reflect.New(reflect.ValueOf(dst).Elem().Type()).Interface().(*Message)
		if err := dec.Decode(newdst); err != nil {
			return errorHandling(err)
		}
		if err := transformRespond(newdst, idMaps); err != nil {
			return errorHandling(err)
		}
	}

	arrayResponse := make([]*Response , 0)
	
	for _ , res := range idMaps{
		arrayResponse = append(arrayResponse, res)
	}
	//responses := Responses{Responses: arrayResponse}
	json := json.NewEncoder(w)

	json.Encode(arrayResponse)
	dec.Token()
	return nil
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
	return 0, fmt.Errorf("Error parsing value tiene un valor invalido : %v", number)
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
		return fmt.Errorf("Error parsing type es incorrecto: %s", obj.Type)
	}
	if err != nil {
		return err
	}

	re.Expenses += expenses
	re.Revenue += income
	re.Transactions = append(re.Transactions, obj)
	return nil
}

func errorHandling(err error) error {
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError

	switch {
	case errors.As(err, &syntaxError):
		Msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
		return &MalformedRequest{Status: http.StatusBadRequest, Msg: Msg}

	case errors.Is(err, io.ErrUnexpectedEOF):
		Msg := fmt.Sprintf("Request body contains badly-formed JSON")
		return &MalformedRequest{Status: http.StatusBadRequest, Msg: Msg}

	case errors.As(err, &unmarshalTypeError):
		Msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
		return &MalformedRequest{Status: http.StatusBadRequest, Msg: Msg}

	case strings.HasPrefix(err.Error(), "json: unknown field "):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		Msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
		return &MalformedRequest{Status: http.StatusBadRequest, Msg: Msg}

	case errors.Is(err, io.EOF):
		Msg := "Request body must not be empty"
		return &MalformedRequest{Status: http.StatusBadRequest, Msg: Msg}

	case err.Error() == "http: request body too large":
		Msg := "Request body must not be larger than 1MB"
		return &MalformedRequest{Status: http.StatusRequestEntityTooLarge, Msg: Msg}
		
	case strings.HasPrefix(err.Error(), "Error parsing"):
		Msg := fmt.Sprintf(err.Error())
		return &MalformedRequest{Status: http.StatusBadRequest , Msg: Msg}
	default:
		return err
	}
}
