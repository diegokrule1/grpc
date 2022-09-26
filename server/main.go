package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/diegokrule/gRpc/proto/chat"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	//server()
	//serverHttp()
	gorillaHttp()
}

func gorillaHttp(){
	r:=mux.NewRouter()
	r.HandleFunc("/",handler).Methods("GET")
	r.HandleFunc("/cust/{id}",cust).Methods("GET")
	r.HandleFunc("/pst",pst).Methods("POST")
	r.HandleFunc("/internal",itn).Methods("POST" )
	srv:=http.Server{
		Addr:              "0.0.0.0:10000",
		Handler:           r,
		TLSConfig:         nil,
		ReadTimeout:       15*time.Second,
		ReadHeaderTimeout: 0,
		WriteTimeout:      15*time.Second,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
		BaseContext:       nil,
		ConnContext:       nil,
	}
	go func(){
		if err:=srv.ListenAndServe();err!=nil{
			log.Println(err)
		}
	}()

	c:=make(chan os.Signal,1)
	signal.Notify(c,os.Interrupt)
	<-c

	ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("Shutting down")
	os.Exit(0)
}

func cust(w http.ResponseWriter,r *http.Request){
	params:=mux.Vars(r)
	id:=params["id"]
	fmt.Fprintf(w, "Welcome User %s",id)

}

func itn(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	decoder:=json.NewDecoder(r.Body)


	body:= struct {
		Name string `json:"name"`
		LastName string `json:"last_name"`
	}{}

	decoder.Decode(&body)
	log.Printf("Name %s, lastName %s", body.Name,body.LastName)
	w.Write([]byte(`hello from service`))
}

func pst(w http.ResponseWriter, r *http.Request){

	client:=http.Client{}

	for i:=1;i<=10;i++{
		go func(j int){
			url:=fmt.Sprintf("http://back-app-rpc.tst%d:10000/internal",j)
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

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func serverHttp() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func server(){
	log.Printf("Starting gRpc server")
	lis, err:=net.Listen("tcp", "192.168.0.144:9000")
	if err!=nil{
		log.Fatalf("failed to listen to :%v", err)
	}
	grpcServer:=grpc.NewServer()
	s:=proto.UnimplementedChatServiceServer{}
	proto.RegisterChatServiceServer(grpcServer, &s)
	//go func(){
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %s", err)
		}
	//}()




}
