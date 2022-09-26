package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	client:=http.Client{}
	for i:=1;i<=1000;i++{
		go func(j int){
			url:=fmt.Sprintf("http://back-app-rpc.tst%d:10000/pst",j)
			log.Printf("Invoking host %s",url)
			jsonBody := []byte(`{"name": "someName","last_name"":"doe"`)
			bodyReader := bytes.NewReader(jsonBody)
			req,err:=http.NewRequest(http.MethodPost,url,bodyReader)
			if err!=nil{
				log.Printf("could not send request %s", err.Error())
				return
			}
			resp, err:=client.Do(req)
			processResponse(resp, err)
		}(i)

	}
}
func processResponse(resp *http.Response, err error){
	if err!=nil{
		log.Printf("could not get response %s", err.Error())
		return
	}
	defer resp.Body.Close()
	b , err:=ioutil.ReadAll(resp.Body)
	if err!=nil{
		log.Printf("could not read response body %s", err.Error())
		return
	}
	log.Printf("got response %s", string(b))
}