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
		panic("Error creating client")
	}

	a := adder{client: client}

	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
	})

	router.HandleFunc("/get-sum", func(w http.ResponseWriter, req *http.Request) {
		logger := glogger.Get(req.Context())

		total, err := a.sum(req.Context())
		if err != nil {
			logger.WithError(err).Error("adder error")
			http.Error(w, "Error", http.StatusInternalServerError)
		}

		returnBody := getSumOrder{total}
		returnBodyInBytes, err := json.Marshal(&returnBody)
		if err != nil {
			logger.WithError(err).Error("Error in marshalling")
			http.Error(w, "Error", http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(returnBodyInBytes)
	})
}

type getSumOrder struct {
	total int
}
