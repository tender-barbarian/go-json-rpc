package rpc

import (
	"bytes"
	"fmt"
	"go-json-rpc/internal/api"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockService struct{}

func (m *MockService) MockMethodOne() string {
	return "Hello from mock method one!"
}

func (m *MockService) MockMethodTwo(argOne string, argTwo string) string {
	return fmt.Sprintf("Hello from mock method two: %s %s", argOne, argTwo)
}

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
			request:    "{\"jsonrpc\":\"2.0\",\"id\":100,\"method\":\"Health.Check\"}",
			response:   "{\"jsonrpc\":\"2.0\",\"result\":\"OK!\",\"id\":100}",
			method:     &api.Health{},
			methodName: "Health",
		},
		{
			name:       "Test 2: Wrong number of parameters",
			request:    "{\"jsonrpc\":\"2.0\",\"id\":100,\"method\":\"Health.Check\", \"params\":[\"test\",\"test1\", \"test2\"]}",
			response:   "{\"Code\":-32602,\"Message\":\"Invalid method parameter(s).\",\"Data\":\"Too many parameters. Method takes 0 params, 3 provided\"}",
			method:     &api.Health{},
			methodName: "Health",
		},
		{
			name:       "Test 3: Wrong method format",
			request:    "{\"jsonrpc\":\"2.0\",\"id\":100,\"method\":\"HealthCheck\", \"params\":[\"test\",\"test1\", \"test2\"]}",
			response:   "{\"Code\":-32601,\"Message\":\"The method does not exist / is not available.\",\"Data\":\"Provided method name needs to have format service.method\"}",
			method:     &MockService{},
			methodName: "MockService",
		},
		{
			name:       "Test 4: Service does not exist",
			request:    "{\"jsonrpc\":\"2.0\",\"id\":100,\"method\":\"Test.Test\", \"params\":[\"test\",\"test1\", \"test2\"]}",
			response:   "{\"Code\":-32601,\"Message\":\"The method does not exist / is not available.\",\"Data\":\"Provided service Test does not exist\"}",
			method:     &MockService{},
			methodName: "MockService",
		},
		{
			name:       "Test 5: Method does not exist",
			request:    "{\"jsonrpc\":\"2.0\",\"id\":100,\"method\":\"Health.Get\", \"params\":[\"test\",\"test1\", \"test2\"]}",
			response:   "{\"Code\":-32601,\"Message\":\"The method does not exist / is not available.\",\"Data\":\"Provided method Get does not exist\"}",
			method:     &api.Health{},
			methodName: "Health",
		},
		{
			name:       "Test 6: Test service with multiple methods",
			request:    "{\"jsonrpc\":\"2.0\",\"id\":100,\"method\":\"MockService.MockMethodOne\"}",
			response:   "{\"jsonrpc\":\"2.0\",\"result\":\"Hello from mock method one!\",\"id\":100}",
			method:     &MockService{},
			methodName: "MockService",
		},
		{
			name:       "Test 7: Test service with multiple methods",
			request:    "{\"jsonrpc\":\"2.0\",\"id\":100,\"method\":\"MockService.MockMethodTwo\",\"params\":[\"test\",\"test1\"]}",
			response:   "{\"jsonrpc\":\"2.0\",\"result\":\"Hello from mock method two: test test1\",\"id\":100}",
			method:     &MockService{},
			methodName: "MockService",
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
