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

func userGroup() {

	Describe("User", func() {
		It("should update user wallet", func() {
			wallet := map[string]any{
				"address": "0x123",
				"network": "bsc",
				"testnet": false,
			}
			data, _ := json.Marshal(wallet)
			req, err := http.NewRequest("PUT", "/users/wallets", bytes.NewBuffer(data))
			Expect(err).ToNot(HaveOccurred())
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", authTokens[0])

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			body := decodeBody(w.Body)
			fmt.Println(body)
			Expect(w.Code).To(Equal(http.StatusOK))
		})
	})

	It("should fetch current user", func() {
		req, err := http.NewRequest("GET", "/users", nil)
		Expect(err).ToNot(HaveOccurred())
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authTokens[0])

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		Expect(w.Code).To(Equal(http.StatusOK))
	})

}
