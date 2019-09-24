package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/parallelcointeam/parallelcoin/pkg/util/cl"
)

// JSONRPC is a handler for sending queries and receiving responses from a JSONRPC endpoint
type JSONRPC struct {
	User     string
	Password string
	Host     string
	Port     int64
}

// Call the RPC server and send a query
func (c *JSONRPC) Call(method string, params interface{}) (interface{}, error) {
	baseURL := fmt.Sprintf("http://%s:%d", c.Host, c.Port)
	client := new(http.Client)
	req, err := http.NewRequest("POST", baseURL, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(c.User, c.Password)
	req.Header.Add("Content-Type", "text/plain")
	args := make(map[string]interface{})
	args["jsonrpc"] = "1.0"
	args["id"] = "BitNodes"
	args["method"] = method
	args["netparams"] = params
	j, err := json.Marshal(args)
	if err != nil {
		fmt.Println(err, cl.Ine())
		return nil, err
	}
	req.Body = ioutil.NopCloser(strings.NewReader(string(j)))
	req.ContentLength = int64(len(string(j)))
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bytes, _ := ioutil.ReadAll(resp.Body)
	var data map[string]interface{}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		fmt.Println(err, cl.Ine())
	}
	if err, found := data["error"]; found && err != nil {
		str, _ := json.Marshal(err)
		return nil, errors.New(string(str))
	}
	if result, found := data["result"]; found {
		return result, nil
	}
	return nil, errors.New("no result")
}

// NewJSONRPC creates a new structure for a JSONRPC connection
func NewJSONRPC(user string, password string, host string, port int64) *JSONRPC {
	c := JSONRPC{user, password, host, port}
	return &c
}
