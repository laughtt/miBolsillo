package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	url = "http://localhost:5000"
)

func Test(f os.FileInfo, client *http.Client, path string) {
	file, err := os.Open(fmt.Sprintf("%s/%s", path, f.Name()))

	if err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadAll(file)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(b))

	resp, err := client.Do(req)
	fmt.Println(f.Name())
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	resp.Body.Close()
	os.Stderr.Close()

	defer func() {
		file.Close()
	}()
}

func main() {
	fmt.Println("URL:>", url)

	path := os.Args[0]
	if path == "" {
		log.Fatal("Introducir en nombre de una carpeta o mas")
	}

	client := &http.Client{}

	for _, path := range os.Args {
		dirPath := "./jsons/" + path
		files, _ := ioutil.ReadDir(dirPath)
		for _, f := range files {
			Test(f, client, dirPath)
		}
	}

	defer func() {
		os.Stderr.Close()
	}()
}
