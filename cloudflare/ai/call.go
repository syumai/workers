//go:build !js

package ai

import (
	"fmt"
	"io"
)

type AI struct {
	instance string
}

func New(varName string) (*AI, error) {
	fmt.Println("NewNamespace called with varName:", varName)
	if varName == "" {
		return nil, fmt.Errorf("%s is undefined", varName)
	}
	return &AI{instance: varName}, nil
}

func (ns *AI) Run(key string, opts map[string]any) (string, error) {
	fmt.Println("Run called with key:", key, "and opts:", opts)
	return "mocked response", nil
}

func (ns *AI) RunWithOutJson(key string, opts map[string]any) (string, error) {
	fmt.Println("Run called with key:", key, "and opts:", opts)
	return "mocked response", nil
}

func (ns *AI) RunReader(key string, opts map[string]any) (io.Reader, error) {
	fmt.Println("RunReader called with key:", key, "and opts:", opts)
	return io.NopCloser(nil), nil
}
