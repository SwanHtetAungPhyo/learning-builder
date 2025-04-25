package common

import (
	"net/http"
	"time"
)

type RpcClient struct {
	Url        string
	httpClient *http.Client
}

func NewHttpClientWithTimeout(timeout time.Duration, url string) *RpcClient {
	client := &http.Client{
		Timeout: timeout,
	}

	rpcClient := &RpcClient{
		httpClient: client,
		Url:        url,
	}
	return rpcClient
}

//func (c *RpcClient) SubmitTransaction(request *http.Request) *http.Response {
//	response, err := c.httpClient.Post("http://localhost:8545/submit", "application/json", request.Body)
//	if err != nil {
//		log.Println(response.StatusCode, response.Status, response.Body)
//		panic(err.Error())
//	}
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//			panic(err.Error())
//		}
//	}(response.Body)
//	return response
//}
