package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/davidebianchi/go-jsonclient"
)

type adder struct {
	client   *jsonclient.Client
	basePath string
}

var (
	errGenericError = errors.New("Generic error")
	errCrudError    = errors.New("Crud error")
)

func (a *adder) sum(cxt context.Context) (int, error) {
	urlPath := fmt.Sprintf("%s/?status=delivered", a.basePath)
	req, err := a.client.NewRequestWithContext(cxt, http.MethodGet, urlPath, nil)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", errGenericError, err.Error())
	}
	orderBody := []orderResponseBody{}
	_, err = a.client.Do(req, &orderBody)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", errCrudError, err.Error())
	}

	sum := 0
	for i := 0; i < len(orderBody); i++ {
		sum += orderBody[i].TotalPrice
	}

	return sum, nil
}

type orderResponseBody struct {
	TotalPrice int `json:"totalPrice"`
}
