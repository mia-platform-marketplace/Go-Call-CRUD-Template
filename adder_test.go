package main

import (
	"context"
	"errors"
	"net/url"
	"testing"

	"github.com/davidebianchi/go-jsonclient"
	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gock.v1"
)

func TestAdder(t *testing.T) {
	crudBaseURL := "http://crud.example.org/"
	crudBasePath := "/my-crud"
	client, err := jsonclient.New(jsonclient.Options{
		BaseURL: crudBaseURL,
	})
	require.NoError(t, err)
	require.NotNil(t, client)
	expectedQuery, err := url.ParseQuery("status=delivered")
	require.NoError(t, err)

	t.Run("should correctly calculate the sum", func(t *testing.T) {
		defer gock.Off()

		a := adder{client: client, basePath: crudBasePath}
		cxt := context.Background()

		mockGetOrdersWithQueryParameters(crudBaseURL, crudBasePath, &expectedQuery, 200, []map[string]int{
			{"totalPrice": 0},
			{"totalPrice": 3},
			{"totalPrice": 1},
		})

		total, err := a.sum(cxt)

		require.NoError(t, err)
		require.Equal(t, 4, total)
	})

	t.Run("should return if crud return 500", func(t *testing.T) {
		defer gock.Off()

		a := adder{client: client, basePath: crudBasePath}
		cxt := context.Background()

		mockGetOrdersWithQueryParameters(crudBaseURL, crudBasePath, &expectedQuery, 500, "")

		total, err := a.sum(cxt)

		require.True(t, errors.Is(err, errCrudError))
		require.Equal(t, 0, total)
	})
}

func mockGetOrdersWithQueryParameters(baseURL string, crudBasePath string, query *url.Values, statusCode int, responseBody interface{}) {
	mockRequest := gock.New(baseURL).
		Get("/orders/")
	if query != nil {
		mockRequest.URLStruct.RawQuery = query.Encode()
	}
	mockRequest.
		Reply(statusCode).
		JSON(responseBody)
}
