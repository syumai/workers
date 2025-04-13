package mock

import (
	"fmt"
	"io"
)

type AIMock struct {
	instance string
}

func NewNamespace(varName string) (*AIMock, error) {
	fmt.Println("NewNamespace called with varName:", varName)
	if varName == "" {
		return nil, fmt.Errorf("%s is undefined", varName)
	}
	return &AIMock{instance: varName}, nil
}

func (ns *AIMock) WaitUntil(task func()) {
	fmt.Println("WaitUntil called")
	go func() {
		task()
	}()
	fmt.Println("Task completed")
}

func (ns *AIMock) Run(key string, opts map[string]interface{}) (string, error) {
	fmt.Println("Run called with key:", key, "and opts:", opts)
	return "mocked response", nil
}

func (ns *AIMock) RunReader(key string, opts map[string]interface{}) (io.Reader, error) {
	fmt.Println("RunReader called with key:", key, "and opts:", opts)
	return io.NopCloser(nil), nil
}
