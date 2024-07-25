package rpc

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testService struct{}

func (m *testService) Hello(args map[string]interface{}) (string, error) {
	return "Hello from Test method one!", nil
}

func (m *testService) HelloParams(args map[string]interface{}) (string, error) {
	var k string
	var v string
	for key, value := range args {
		k = key
		v = value.(string)
	}
	return fmt.Sprintf("Hello from test method two! Key: %s, value: %s", k, v), nil
}

func TestHandler_ServeHTTP(t *testing.T) {
	tests := []struct {
		name       string
		request    string
		method     func(map[string]interface{}) (string, error)
		methodName string
		response   string
	}{
		{
			name:       "Test 1: Test RPC",
			request:    "{\"jsonrpc\":\"2.0\",\"id\":100,\"method\":\"TestService.Hello\"}",
			response:   "{\"jsonrpc\":\"2.0\",\"result\":\"Hello from Test method one!\",\"id\":100}",
			method:     (&testService{}).Hello,
			methodName: "TestService.Hello",
		},
		{
			name:       "Test 2: Test RPC - method with params",
			request:    "{\"jsonrpc\":\"2.0\",\"id\":100,\"method\":\"TestService.HelloParams\", \"params\":{\"test\":\"test2\"}}",
			response:   "{\"jsonrpc\":\"2.0\",\"result\":\"Hello from test method two! Key: test, value: test2\",\"id\":100}",
			method:     (&testService{}).HelloParams,
			methodName: "TestService.HelloParams",
		},
		{
			name:       "Test 3: Invalid parameters",
			request:    "{\"jsonrpc\":\"2.0\",\"id\":100,\"method\":\"TestService.Hello\", \"params\":[\"test\",\"test2\"]}",
			response:   "{\"Code\":-32700,\"Message\":\"Invalid JSON was received by the server. An error occurred on the server while parsing the JSON text.\",\"Data\":\"json: cannot unmarshal array into Go struct field request.params of type map[string]interface {}\"}",
			method:     (&testService{}).Hello,
			methodName: "TestService.Hello",
		},
		{
			name:       "Test 4: Method does not exist",
			request:    "{\"jsonrpc\":\"2.0\",\"id\":100,\"method\":\"Health.Get\", \"params\":{\"test\":\"test1\", \"test2\":\"test2\"}}",
			response:   "{\"Code\":-32601,\"Message\":\"The method does not exist / is not available.\",\"Data\":\"\"}",
			method:     (&testService{}).Hello,
			methodName: "TestService.Hello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{}
			err := h.Register(tt.methodName, tt.method)
			if err != nil {
				t.Fatal(err)
			}

			handler := http.HandlerFunc(h.ServeHTTP)

			rr := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/", bytes.NewReader([]byte(tt.request)))
			if err != nil {
				t.Fatal(err)
			}

			handler.ServeHTTP(rr, req)
			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
			}

			expected := tt.response
			if rr.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
			}
		})
	}
}

func TestHandler_Register(t *testing.T) {
	tests := []struct {
		name       string
		method     func(map[string]interface{}) (string, error)
		methodName string
		wantErr    bool
	}{
		{
			name:       "Test 1",
			method:     (&testService{}).Hello,
			methodName: "TestService.Hello",
			wantErr:    false,
		},
		{
			name:       "Test 2: Nil function",
			method:     nil,
			methodName: "",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{}
			if err := h.Register(tt.methodName, tt.method); (err != nil) != tt.wantErr {
				t.Errorf("Handler.Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
