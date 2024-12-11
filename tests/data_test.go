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
			"title":                   "sample service",
			"description":             "service desc",
			"payment_currency":        "USD",
			"skills":                  []string{"Skill1"},
			"project_length":          "LESS_THAN_A_DAY",
			"job_category_id":         "",
			"commitment_hours_lower":  "1",
			"commitment_hours_higher": "1",
			"payment_range_lower":     "2",
			"payment_range_higher":    "2",
			"kind":                    "SERVICE",
			"work_samples":            []string{},
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
