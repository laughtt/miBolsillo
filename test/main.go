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

const (
	url = "http://localhost:5000"
)


func Test(f os.FileInfo , client *http.Client, path string){
	file, err := os.Open(fmt.Sprintf("%s/%s",path,f.Name()))
	
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
			resp.Body.Close()
		}
		os.Stderr.Close()

		defer func(){
			file.Close()
		}()
}


func main() {
	fmt.Println("URL:>", url)

	files, _ := ioutil.ReadDir("./jsons/badTest")

	client := &http.Client{}

	for _, f := range files {
		Test(f , client,"./jsons/badTest")
	}

	files , _ = ioutil.ReadDir("./jsons/big")
	for _, f := range files {
		for i := 0; i < 100 ; i++{
			Test(f , client ,"./jsons/big")
		}
	}

	
	files , _ = ioutil.ReadDir("./jsons/correct")
	for _, f := range files {
		Test(f , client, "./jsons/correct")
	}
	defer func(){
		os.Stderr.Close()
	}()


}
