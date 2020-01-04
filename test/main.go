package main

import (
	"bytes"
	"strings"	
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)


func main() {
	url := "http://localhost:5000"
	fmt.Println("URL:>", url)



	files, _ := ioutil.ReadDir("./jsons")

	client := &http.Client{}

	for _, f := range files {
		file, err := os.Open(fmt.Sprintf("./jsons/%s",f.Name()))
	
		if err != nil {
			log.Fatal(err)
		}
		b, err := ioutil.ReadAll(file)

		req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(b))


		resp , err := client.Do(req)
		fmt.Println(f.Name())
		fmt.Println("response Status:", resp.Status)
		//fmt.Println("response Headers:", resp.Header)
		if strings.HasPrefix(resp.Status,"400"){
			body , _ := ioutil.ReadAll(resp.Body)
			fmt.Println("response Body:", string(body))
		}

	}

}
