package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type ResponseApicep struct {
	Code  string `json:code`
	State string `json:state`
	City  string `json:city`
}

type ResponseViaCep struct {
	Cep        string `json:cep`
	Uf         string `json:uf`
	Localidade string `json:localidade`
}

func requestApiCep(cep string, ch chan ResponseApicep) {
	req, err := http.Get("https://cdn.apicep.com/file/apicep/" + cep + ".json")
	if err != nil {
		log.Fatalln("Request Apicep Error : ", req.Status)

	}
	defer req.Body.Close()
	response, err := io.ReadAll(req.Body)
	if err != nil {
		log.Fatalln("Request Apicep Error : ", err)
	}
	var data ResponseApicep

	json.Unmarshal(response, &data)

	ch <- data
}

func requestViaCep(cep string, ch chan ResponseViaCep) {
	req, err := http.Get("http://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		log.Fatalln("Request viaCep Error : ", req.Status)
	}
	defer req.Body.Close()
	response, err := io.ReadAll(req.Body)
	if err != nil {
		log.Fatalln("Request viaCep Error : ", err)
	}
	var data ResponseViaCep
	json.Unmarshal(response, &data)

	ch <- data
}

func main() {
	apiCepChannel := make(chan ResponseApicep)
	viaCepChannel := make(chan ResponseViaCep)

	go requestApiCep("59695-000", apiCepChannel)
	go requestViaCep("59695000", viaCepChannel)

	select {
	case response := <-viaCepChannel:
		fmt.Printf("Via cep é mais rapida: %s - %s\n", response.Localidade, response.Uf)
	case response := <-apiCepChannel:
		fmt.Printf("Api cep é mais rapida: %s - %s\n", response.City, response.State)
	case <-time.After(time.Second):
		fmt.Println("timeout")
	}
}
