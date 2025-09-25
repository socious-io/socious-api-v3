package tests_test

import (
	"socious/src/apps/models"

	"github.com/gin-gonic/gin"
)

var (
	authTokens        = []string{}
	authRefreshTokens = []string{}

	usersData = []*models.User{
		{
			Username:         "test",
			Email:            "test@test.com",
			IdentityVerified: true,
			Events:           []string{},
			Tags:             []string{},
		},
		{
			Username:         "test2",
			Email:            "test2@test.com",
			IdentityVerified: true,
			Events:           []string{},
			Tags:             []string{},
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
			"payment_type":      "CRYPTO",
			"commitment_period": "MONTHLY",
		},
	}
)
