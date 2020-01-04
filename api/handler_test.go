// handlers_test.go
package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)


func checkFile(f os.FileInfo, client *http.Client, path string, t *testing.T) {
	file, err := os.Open(fmt.Sprintf("%s/%s", path, f.Name()))

	if err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadAll(file)

	req, err := http.NewRequest("PUT", "/", bytes.NewBuffer(b))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateInvoice)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"alive": true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
	os.Stderr.Close()
	file.Close()
}

func TestHandler(t *testing.T) {

	files, _ := ioutil.ReadDir("../test/jsons/badTest")

	client := &http.Client{}

	for _, f := range files {
		checkFile(f, client, "../test/jsons/badTest", t)
	}

	// files, _ = ioutil.ReadDir("./jsons/big")
	// for _, f := range files {
	// 	for i := 0; i < 100; i++ {
	// 		TestHealthCheckHandler(f, client, "./jsons/big")
	// 	}
	// }

	files, _ = ioutil.ReadDir("../test/jsons/correct")
	for _, f := range files {
		checkFile(f, client, "../test/jsons/correct", t)
	}
}
