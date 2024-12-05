package tests_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"socious/src/apps/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func serviceGroup() {

	BeforeAll(func() {
		ctx := context.Background()
		//Creating Job Category
		jcWage := jobCategoryData[0]["hourly_wage_dollars"].(float64)
		jobCategory := models.JobCategory{
			Name:              jobCategoryData[0]["name"].(string),
			HourlyWageDollars: &jcWage,
		}
		jobCategory.Create(ctx)
		jobCategoryData[0]["id"] = jobCategory.ID
		servicesData[0]["job_category_id"] = jobCategory.ID

		//Creating Media
		media := models.Media{
			Filename:   "media_filename.ext",
			IdentityID: uuid.MustParse(usersData[0]["id"].(string)),
			URL:        "media_url",
		}
		media.Create(ctx)
		servicesData[0]["work_samples"] = append(servicesData[0]["work_samples"].([]string), media.ID.String())
	})

	It("should create service", func() {
		for i, data := range servicesData {
			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(data)
			req, _ := http.NewRequest("POST", "/services", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", authTokens[0])
			router.ServeHTTP(w, req)
			body := decodeBody(w.Body)
			Expect(w.Code).To(Equal(http.StatusCreated))
			servicesData[i]["id"] = body["id"]
		}
	})

	It("should get all services with pagination", func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/services", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authTokens[0])
		router.ServeHTTP(w, req)
		body := decodeBody(w.Body)
		Expect(w.Code).To(Equal(http.StatusOK))
		Expect(len(body["results"].([]interface{}))).To(Equal(len(servicesData)))
	})

	It("should get service", func() {
		for _, data := range servicesData {
			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(data)
			req, _ := http.NewRequest("GET", fmt.Sprintf("/services/%s", data["id"]), bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", authTokens[0])
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusOK))
			body := decodeBody(w.Body)
			Expect(body["id"]).To(Equal(data["id"]))
		}
	})

	It("should delete service", func() {
		for _, data := range servicesData {
			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(data)
			req, _ := http.NewRequest("DELETE", fmt.Sprintf("/services/%s", data["id"]), bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", authTokens[0])
			router.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusOK))
			body := decodeBody(w.Body)
			bodyExpect(body, gin.H{"message": "success"})
		}
	})
}
