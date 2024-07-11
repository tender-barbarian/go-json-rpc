package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
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
	name    string
	methods map[string]reflect.Value
}

func (h *Handler) Register(name string, method interface{}) error {
	if h.Services == nil {
		h.Services = make(map[string]service)
	}

	if method == nil {
		return fmt.Errorf("method cannot be nil!")
	}

	// get methods from service
	var methods = make(map[string]reflect.Value)
	t := reflect.TypeOf(method)
	for i := 0; i < t.NumMethod(); i++ {
		val := reflect.ValueOf(method).MethodByName(t.Method(i).Name)
		if !val.IsValid() {
			return fmt.Errorf("methods cannot be anonymous!")
		}
		methods[t.Method(i).Name] = val
	}

	h.Services[name] = service{
		name:    name,
		methods: methods,
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

	// find service
	m := strings.Split(req.Method, ".")
	if len(m) < 2 {
		rpcErr, _ := json.Marshal(&Error{Code: MethodNotFound, Message: "The method does not exist / is not available.", Data: "Provided method name needs to have format service.method"})
		fmt.Fprint(w, string(rpcErr))
		return
	}
	serviceName := m[0]
	methodName := m[1]
	service, serviceFound := h.Services[serviceName]
	if !serviceFound {
		rpcErr, _ := json.Marshal(&Error{Code: MethodNotFound, Message: "The method does not exist / is not available.", Data: fmt.Sprintf("Provided service %s does not exist", serviceName)})
		fmt.Fprint(w, string(rpcErr))
		return
	}

	// find method
	method, methodFound := service.methods[methodName]
	if !methodFound {
		rpcErr, _ := json.Marshal(&Error{Code: MethodNotFound, Message: "The method does not exist / is not available.", Data: fmt.Sprintf("Provided method %s does not exist", methodName)})
		fmt.Fprint(w, string(rpcErr))
		return
	}

	// make sure provided params are valid
	if method.Type().NumIn() != len(req.Params) {
		rpcErr, _ := json.Marshal(&Error{Code: InvalidParams, Message: "Invalid method parameter(s).", Data: fmt.Sprintf("Too many parameters. Method takes %d params, %d provided", method.Type().NumIn(), len(req.Params))})
		fmt.Fprint(w, string(rpcErr))
		return
	}

	// compile params into call input
	inputs := make([]reflect.Value, len(req.Params))
	for i := range req.Params {
		inputs[i] = reflect.ValueOf(req.Params[i])
	}
	ret := method.Call(inputs)

	response, _ := json.Marshal(&Response{Jsonrpc: req.Jsonrpc, Result: ret[0].Interface().(string), Id: req.Id})
	fmt.Fprint(w, string(response))

}
