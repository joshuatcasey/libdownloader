package fakes

import (
	"net/http"
	"sync"
)

type SimpleHttpClient struct {
	GetCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Url string
		}
		Returns struct {
			Response *http.Response
			Error    error
		}
		Stub func(string) (*http.Response, error)
	}
}

func (f *SimpleHttpClient) Get(param1 string) (*http.Response, error) {
	f.GetCall.mutex.Lock()
	defer f.GetCall.mutex.Unlock()
	f.GetCall.CallCount++
	f.GetCall.Receives.Url = param1
	if f.GetCall.Stub != nil {
		return f.GetCall.Stub(param1)
	}
	return f.GetCall.Returns.Response, f.GetCall.Returns.Error
}
