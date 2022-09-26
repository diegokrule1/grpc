package main

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/",http.FileServer(http.Dir("/Users/alexis/go/src/github.com/diegokrule/gRpc/app/public")))
	log.Fatalln(http.ListenAndServe(":8082",nil))
}