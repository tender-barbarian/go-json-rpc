package rpc

import (
	"bytes"
	"go-json-rpc/internal/api"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockMethod struct{}

func (m *MockMethod) MockMethod() {}

func TestHandler_ServeHTTP(t *testing.T) {
	tests := []struct {
		name       string
		request    string
		method     interface{}
		methodName string
		response   string
	}{
		{
			name:       "Test 1: Test Health Check method",
			request:    "{\"jsonrpc\":\"2.0\",\"id\":100,\"method\":\"HealthCheck\"}",
			response:   "{\"jsonrpc\":\"2.0\",\"result\":\"OK!\",\"id\":100}",
			method:     &api.Health{},
			methodName: "HealthCheck",
		},
		{
			name:       "Test 2: Wrong number of parameters",
			request:    "{\"jsonrpc\":\"2.0\",\"id\":100,\"method\":\"HealthCheck\", \"params\":[\"test\",\"test1\", \"test2\"]}",
			response:   "{\"Code\":-32602,\"Message\":\"Invalid method parameter(s).\",\"Data\":\"Too many parameters. Method takes 0 params, 3 provided\"}",
			method:     &api.Health{},
			methodName: "HealthCheck",
		},
		{
			name:       "Test 3: Method does not exist",
			request:    "{\"jsonrpc\":\"2.0\",\"id\":100,\"method\":\"Health\", \"params\":[\"test\",\"test1\", \"test2\"]}",
			response:   "{\"Code\":-32601,\"Message\":\"The method does not exist / is not available.\",\"Data\":\"\"}",
			method:     &MockMethod{},
			methodName: "MockMethod",
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
		method     interface{}
		methodName string
		wantErr    bool
	}{
		{
			name:       "Test 1",
			method:     &api.Health{},
			methodName: "HealthCheck",
			wantErr:    false,
		},
		{
			name:       "Test 2: Anonymous function",
			method:     func() {},
			methodName: "Anon",
			wantErr:    true,
		},
		{
			name:       "Test 3: Nil function",
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
