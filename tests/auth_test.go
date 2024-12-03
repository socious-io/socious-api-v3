package tests_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func authGroup() {

	authExecuted = true

	It("should register user", func() {
		w := httptest.NewRecorder()
		reqBody, _ := json.Marshal(usersData[0])
		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		body := decodeBody(w.Body)
		Expect(w.Code).To(Equal(200))
		bodyExpect(body, gin.H{"access_token": "<ANY>", "refresh_token": "<ANY>", "token_type": "Bearer"})
		authTokens = append(authTokens, body["access_token"].(string))
		authRefreshTokens = append(authRefreshTokens, body["refresh_token"].(string))
	})

}
