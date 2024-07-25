package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	parseError     = ErrorWrapper{Code: -32700, Message: "Invalid JSON was received by the server. An error occurred on the server while parsing the JSON text."}
	invalidRequest = ErrorWrapper{Code: -32600, Message: "The JSON sent is not a valid Request object."}
	methodNotFound = ErrorWrapper{Code: -32601, Message: "The method does not exist / is not available."}
	// invalidParams  = ErrorWrapper{Code: -32602, Message: "Invalid method parameter(s)."}
	// internalError = ErrorWrapper{Code: -32603, Message: "Internal JSON-RPC error."}
	methodError = ErrorWrapper{Code: -32001, Message: "Method returned error."}
)

type ErrorWrapper struct {
	Code    int
	Message string
	Data    string
}

func rpcError(w http.ResponseWriter, errorWrapper ErrorWrapper, err error) {
	if err != nil {
		errorWrapper.Data = err.Error()
	}

	rpcErr, internalErr := json.Marshal(errorWrapper)
	if internalErr != nil {
		fmt.Fprintf(w, "{\"Code\": -32603, \"Message\": \"Internal JSON-RPC error\", \"Data\": %v\"}", internalErr)
		return
	}

	fmt.Fprint(w, string(rpcErr))
}
