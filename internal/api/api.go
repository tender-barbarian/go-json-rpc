package api

import (
	"fmt"
)

type HealthRets struct {
	Result string
}

type HealthParams struct{}

type Health struct{}

func (h Health) Check(params map[string]interface{}) (string, error) {
	return "OK!", nil
}

type RecipeParams struct {
	Title       string   `json:"title"`
	Ingredients []string `json:"ingredients"`
	Method      string   `json:"method"`
}

type Recipe struct{}

func (r Recipe) Create(params map[string]interface{}) error {
	return fmt.Errorf("dupsko")
}

func (r Recipe) Get() string {
	return ""
}
