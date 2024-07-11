package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
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
	Jsonrpc string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []string `json:"params,omitempty"`
	Id      int      `json:"id"`
}

type Response struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  string `json:"result"`
	Id      int    `json:"id"`
}

type Handler struct {
	Services map[string]service
}

type service struct {
	name   string
	method reflect.Value
}

func (h *Handler) Register(name string, method interface{}) error {
	if h.Services == nil {
		h.Services = make(map[string]service)
	}

	if method == nil {
		return fmt.Errorf("method cannot be nil!")
	}

	m := reflect.ValueOf(method)

	if !m.MethodByName(name).IsValid() {
		return fmt.Errorf("method cannot be anonymous!")
	}

	h.Services[name] = service{
		name:   name,
		method: m,
	}

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

	service, ok := h.Services[req.Method]
	if ok {
		method := service.method.MethodByName(service.name)
		if method.Type().NumIn() != len(req.Params) {
			rpcErr, _ := json.Marshal(&Error{Code: InvalidParams, Message: "Invalid method parameter(s).", Data: fmt.Sprintf("Too many parameters. Method takes %d params, %d provided", method.Type().NumIn(), len(req.Params))})
			fmt.Fprint(w, string(rpcErr))
			return
		}

		inputs := make([]reflect.Value, len(req.Params))
		for i := range req.Params {
			inputs[i] = reflect.ValueOf(req.Params[i])
		}

		ret := method.Call(inputs)

		response, _ := json.Marshal(&Response{Jsonrpc: req.Jsonrpc, Result: ret[0].Interface().(string), Id: req.Id})
		fmt.Fprint(w, string(response))
	} else {
		rpcErr, _ := json.Marshal(&Error{Code: MethodNotFound, Message: "The method does not exist / is not available.", Data: ""})
		fmt.Fprint(w, string(rpcErr))
	}
}
