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

const (
	badTestPath  = "../test/jsons/badTest"
	goodTestPath = "../test/jsons/correct"
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

	if path == badTestPath {
		if status := rr.Code; status == http.StatusAccepted {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)

		}

	}

	if path == goodTestPath {
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	}
	os.Stderr.Close()
	file.Close()
}

func TestHandler(t *testing.T) {

	files, _ := ioutil.ReadDir(badTestPath)

	client := &http.Client{}

	for _, f := range files {
		checkFile(f, client, badTestPath, t)
	}

	// files, _ = ioutil.ReadDir("./jsons/big")
	// for _, f := range files {
	// 	for i := 0; i < 100; i++ {
	// 		TestHealthCheckHandler(f, client, "./jsons/big")
	// 	}
	// }

	files, _ = ioutil.ReadDir(goodTestPath)
	for _, f := range files {
		checkFile(f, client, goodTestPath, t)
	}
}
