package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var jsonrpc = NewJSONRPC("user", "pa55word", "127.0.0.1", 11349)

func main() {
	fmt.Println("Simple Block Explorer")
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/getinfo", getinfoHandler)
	http.HandleFunc("/getallblocks", getallblocksHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func getinfoHandler(	w http.ResponseWriter, r *http.Request) {
	response, err := jsonrpc.Call("getinfo", nil)
	if err != nil {
		fmt.Println("ERROR", err.Error())
	}
	jsonResponse, err := json.MarshalIndent(response, "  ", "")
	if err != nil {
		fmt.Println("ERROR", err.Error())
	}
	fmt.Fprintf(w, string(jsonResponse))
}
func rootHandler(	w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the root")
}
func getallblocksHandler(	w http.ResponseWriter, r *http.Request) {
	response, err := jsonrpc.Call("getblockcount", nil)
	height := uint32(response.(float64))
	if err != nil {
		fmt.Println("ERROR", err.Error())
	}
	var i uint32
	for i = 0; i < height; i++ {
		response, err = jsonrpc.Call("getblockhash", []uint32{i})
		if err != nil {
			fmt.Println("ERROR", err.Error())
		}
		response, err = jsonrpc.Call("getblock", []string{response.(string)})
		if err != nil {
			fmt.Println("ERROR", err.Error())
		}
		jsonResponse, err := json.MarshalIndent(response, "  ", "")
		if err != nil {
			fmt.Println("ERROR", err.Error())
		}
		fmt.Fprintf(w, "%s\n", jsonResponse)
	}
}
