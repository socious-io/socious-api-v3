package tests_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func contractGroup() {

	BeforeAll(func() {
		// ctx := context.Background()

	})

	It("should create contract", func() {
		for i, data := range contractsData {
			w := httptest.NewRecorder()
			data["client_id"] = usersData[1].ID
			reqBody, _ := json.Marshal(data)
			req, _ := http.NewRequest("POST", "/contracts", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", authTokens[0])
			router.ServeHTTP(w, req)
			body := decodeBody(w.Body)
			Expect(w.Code).To(Equal(http.StatusCreated))
			contractsData[i]["id"] = body["id"]
		}
	})

	It("should update contract", func() {
		for i, data := range contractsData {
			w := httptest.NewRecorder()
			data["client_id"] = usersData[1].ID
			reqBody, _ := json.Marshal(data)
			req, _ := http.NewRequest("PATCH", fmt.Sprintf("/contracts/%s", data["id"]), bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", authTokens[0])
			router.ServeHTTP(w, req)
			body := decodeBody(w.Body)
			Expect(w.Code).To(Equal(http.StatusAccepted))
			contractsData[i]["id"] = body["id"]
		}
	})

	It("should get contracts", func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/contracts", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authTokens[0])
		router.ServeHTTP(w, req)
		body := decodeBody(w.Body)
		Expect(len(body["results"].([]interface{}))).To(Equal(1))
		Expect(w.Code).To(Equal(http.StatusOK))
	})

	It("should get contract", func() {
		for _, data := range contractsData {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/contracts/%s", data["id"]), nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", authTokens[0])
			router.ServeHTTP(w, req)
			body := decodeBody(w.Body)
			Expect(body["id"]).To(Equal(data["id"]))
			Expect(w.Code).To(Equal(http.StatusOK))
		}
	})
}
