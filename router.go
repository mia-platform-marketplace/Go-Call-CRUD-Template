/*
 * Copyright 2019 Mia srl
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/davidebianchi/go-jsonclient"
	"github.com/gorilla/mux"
	"github.com/mia-platform/glogger"
)

func setupRouter(router *mux.Router, env *EnvironmentVariables) {
	opts := jsonclient.Options{
		BaseURL: env.CrudBaseURL,
	}
	client, err := jsonclient.New(opts)
	if err != nil {
		panic(fmt.Errorf("%w: error creating client", err))
	}

	adderInstance := adder{
		client:   client,
		basePath: env.CrudBasePath,
	}

	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
	})

	router.HandleFunc("/get-sum", func(w http.ResponseWriter, req *http.Request) {
		logger := glogger.Get(req.Context())

		total, err := adderInstance.sum(req.Context())
		if err != nil {
			logger.WithError(err).Error("adder error")

			returnWithMarshal(w, http.StatusInternalServerError, &errorResponse{
				Message: err.Error(),
			})
			return
		}

		returnBody := getSumOrderResponse{total}
		returnWithMarshal(w, http.StatusOK, &returnBody)
	})
}

type getSumOrderResponse struct {
	Total int `json:"total"`
}
type errorResponse struct {
	Message string `json:"message"`
}

func returnWithMarshal(w http.ResponseWriter, statusCode int, obj interface{}) {
	body, err := json.Marshal(obj)
	if err != nil {
		http.Error(w, "InternalError", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", "application/json")
	w.Write(body)
}
