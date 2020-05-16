/*
 * Copyright © 2020-present Mia s.r.l.
 * All rights reserved
 */

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gock.v1"
)

var jsonschemaPath = "./serviceconfiguration/config.schema.json"
var providerTokenPassPhrase = "providerTokenPassPhrase"

func TestSumTotalPrice(t *testing.T) {
	crudBaseURL := "http://crud.example.org"
	env := EnvironmentVariables{
		CrudBaseURL: crudBaseURL,
	}
	expectedQuery, err := url.ParseQuery("status=delivered")
	require.NoError(t, err)

	t.Run("should return the right number", func(t *testing.T) {
		defer gock.Off()

		mockGetOrdersWithQueryParameters(crudBaseURL, &expectedQuery, 200, []map[string]int{
			{"totalPrice": 0},
			{"totalPrice": 3},
			{"totalPrice": 1},
		})
		serviceRouter := mux.NewRouter()
		setupRouter(serviceRouter, &env)

		urlPath := "/get-sum"
		r := httptest.NewRequest(http.MethodGet, urlPath, nil)
		w := httptest.NewRecorder()
		serviceRouter.ServeHTTP(w, r)

		require.Equal(t, http.StatusOK, w.Code, "Unexpected status code.")

		v := getSumOrderResponse{}
		err := json.NewDecoder(w.Body).Decode(&v)
		require.NoError(t, err)
		require.Equal(t, getSumOrderResponse{Total: 4}, v)
	})

	t.Run("should return error if crud call fails", func(t *testing.T) {
		defer gock.Off()

		mockGetOrdersWithQueryParameters(crudBaseURL, &expectedQuery, 500, []map[string]int{})
		serviceRouter := mux.NewRouter()
		setupRouter(serviceRouter, &env)

		urlPath := "/get-sum"
		r := httptest.NewRequest(http.MethodGet, urlPath, nil)
		w := httptest.NewRecorder()
		serviceRouter.ServeHTTP(w, r)

		require.Equal(t, http.StatusInternalServerError, w.Code, "Unexpected status code.")

		v := errorResponse{}
		err := json.NewDecoder(w.Body).Decode(&v)
		require.NoError(t, err)
		require.NotEmpty(t, v.Message)
	})
}
