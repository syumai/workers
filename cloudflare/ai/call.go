//go:build !js

package ai

import (
	"fmt"
	"io"
)

type AI struct {
	instance string
}

func NewNamespace(varName string) (*AI, error) {
	fmt.Println("NewNamespace called with varName:", varName)
	if varName == "" {
		return nil, fmt.Errorf("%s is undefined", varName)
	}
	return &AI{instance: varName}, nil
}

func (ns *AI) WaitUntil(task func()) {
	fmt.Println("WaitUntil called")
	go func() {
		task()
	}()
	fmt.Println("Task completed")
}

func (ns *AI) Run(key string, opts map[string]interface{}) (string, error) {
	fmt.Println("Run called with key:", key, "and opts:", opts)
	return "mocked response", nil
}

func (ns *AI) RunReader(key string, opts map[string]interface{}) (io.Reader, error) {
	fmt.Println("RunReader called with key:", key, "and opts:", opts)
	return io.NopCloser(nil), nil
}
