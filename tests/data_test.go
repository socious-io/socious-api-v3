package tests_test

import "github.com/gin-gonic/gin"

var (
	intKey            = ""
	authTokens        = []string{}
	authRefreshTokens = []string{}

	usersData = []gin.H{
		{
			"first_name": "TestName",
			"last_name":  "TestLastName",
			"username":   "test",
			"email":      "test@test.com",
			"password":   "test123456",
		},
	}

	servicesData = []gin.H{
		{
			"title":               "sample service",
			"description":         "",
			"payment_currency":    "",
			"skills":              []string{"Skill1"},
			"job_category_id":     "282bd9ef-73cf-4c4c-bcf0-09615930d408",
			"service_total_hours": 1,
			"service_price":       2,
			"service_length":      "LESS_THAN_A_DAY",
			"work_samples":        []string{"0001ef01-4d4a-4665-b73a-5f558ffaa2a0"},
		},
	}
)
