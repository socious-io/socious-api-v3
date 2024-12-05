package tests_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"socious/src/apps/auth"

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
		bodyExpect(body, gin.H{"access_token": "<ANY>", "refresh_token": "<ANY>", "token_type": "Bearer"})
		Expect(w.Code).To(Equal(200))
		authTokens = append(authTokens, body["access_token"].(string))
		authRefreshTokens = append(authRefreshTokens, body["refresh_token"].(string))
		claims, _ := auth.VerifyToken(body["access_token"].(string))
		usersData[0]["id"] = claims.ID
	})

}
