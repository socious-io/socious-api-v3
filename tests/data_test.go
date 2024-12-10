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
		{
			"first_name": "TestName2",
			"last_name":  "TestLastName2",
			"username":   "test2",
			"email":      "test2@test.com",
			"password":   "test123456",
		},
	}

	jobCategoryData = []gin.H{
		{
			"name":                "OTHER",
			"hourly_wage_dollars": 20.1,
		},
	}

	servicesData = []gin.H{
		{
			"title":               "sample service",
			"description":         "",
			"payment_currency":    "",
			"skills":              []string{"Skill1"},
			"job_category_id":     "",
			"service_total_hours": 1,
			"service_price":       2,
			"service_length":      "LESS_THAN_A_DAY",
			"work_samples":        []string{},
		},
	}

	contractsData = []gin.H{
		{
			"title":             "sample contract",
			"description":       "test desc",
			"total_amount":      2000,
			"currency":          "USD",
			"type":              "PAID",
			"commitment_period": "MONTHLY",
		},
	}
)
