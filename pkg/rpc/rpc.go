package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type request struct {
	Jsonrpc string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params,omitempty"`
	Id      int                    `json:"id"`
}

type response struct {
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
	var req request
	if r.Body == nil {
		rpcError(w, invalidRequest, fmt.Errorf("Request body cannot be nil"))
		return
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		rpcError(w, parseError, err)
		return
	}

	// find registered method
	method, methodFound := h.Methods[req.Method]
	if !methodFound {
		rpcError(w, methodNotFound, nil)
		return
	}

	// execute callback
	m := *method
	res, err := m(req.Params)
	if err != nil {
		rpcError(w, methodError, err)
		return
	}

	response, _ := json.Marshal(&response{Jsonrpc: req.Jsonrpc, Result: res, Id: req.Id})
	fmt.Fprint(w, string(response))
}
