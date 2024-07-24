package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	ParseError     = -32700
	InvalidRequest = -32600
	MethodNotFound = -32601
	InvalidParams  = -32602
	InternalError  = -32603
	ServerError    = -32000
)

type Error struct {
	Code    int
	Message string
	Data    string
}

type Request struct {
	Jsonrpc string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params,omitempty"`
	Id      int                    `json:"id"`
}

type Response struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  string `json:"result"`
	Id      int    `json:"id"`
}

type Handler struct {
	Methods map[string]*func(map[string]interface{}) (string, error)
}

func (h *Handler) Register(name string, method func(map[string]interface{}) (string, error)) error {
	if method == nil {
		return fmt.Errorf("method cannot be nil!")
	}

	if h.Methods == nil {
		h.Methods = make(map[string]*func(map[string]interface{}) (string, error))
	}

	h.Methods[name] = &method
	return nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req Request
	if r.Body == nil {
		rpcErr, _ := json.Marshal(&Error{Code: InvalidRequest, Message: "The JSON sent is not a valid Request object.", Data: "Empty object received"})
		fmt.Fprint(w, string(rpcErr))
		return
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		rpcErr, _ := json.Marshal(&Error{Code: ParseError, Message: "Invalid JSON was received by the server. An error occurred on the server while parsing the JSON text.", Data: err.Error()})
		fmt.Fprint(w, string(rpcErr))
		return
	}

	// find registered method
	method, methodFound := h.Methods[req.Method]
	if !methodFound {
		rpcErr, _ := json.Marshal(&Error{Code: MethodNotFound, Message: "The method does not exist / is not available.", Data: ""})
		fmt.Fprint(w, string(rpcErr))
		return
	}

	// execute callback
	m := *method
	res, err := m(req.Params)
	if err != nil {
		rpcErr, _ := json.Marshal(&Error{Code: ServerError, Message: "Method returned error.", Data: err.Error()})
		fmt.Fprint(w, string(rpcErr))
		return
	}

	response, _ := json.Marshal(&Response{Jsonrpc: req.Jsonrpc, Result: res, Id: req.Id})
	fmt.Fprint(w, string(response))

}
