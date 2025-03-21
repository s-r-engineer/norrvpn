package libraryHttp

import (
	"net/http"
	url2 "net/url"
	"time"
)

var defaultClient *http.Client
var defaultTimeoutInSeconds int64 = 5

const defaultRetries int = 5

func init() {
	defaultClient = &http.Client{
		Timeout: time.Duration(defaultTimeoutInSeconds) * time.Second,
	}
}

func Do(req *http.Request) (*http.Response, error) {
	return do(req)
}

func GetUrl(ref string) (*http.Response, error) {
	req := http.Request{}
	url, err := url2.Parse(ref)
	if err != nil {
		return nil, err
	}
	req.URL = url
	return do(&req)
}

func do(req *http.Request) (response *http.Response, err error) {
	for i := defaultRetries; i > 0; i-- {
		response, err = defaultClient.Do(req)
		if err == nil {
			return response, nil
		}
	}
	return response, err
}
