package views

import (
	"fmt"
	"net/http"
	"net/url"
	"socious/src/apps/auth"
	"socious/src/apps/models"

	"github.com/gin-gonic/gin"
)

func ssoGroup(router *gin.Engine) {
	g := router.Group("sso")

	router.LoadHTMLGlob("src/apps/templates/*")

	g.GET("/login", func(c *gin.Context) {
		redirect_url := c.Query("redirect_url")

		c.HTML(http.StatusOK, "login.html", gin.H{
			"redirect_url": redirect_url,
		})
	})

	g.POST("/login", func(c *gin.Context) {
		redirect_url := c.Query("redirect_url")

		loginForm := new(auth.LoginForm)
		if err := c.ShouldBind(loginForm); err != nil {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"redirect_url": redirect_url,
				"error":        err.Error(),
			})
			return
		}

		u, err := models.GetUserByEmail(loginForm.Email)
		if err != nil {
			fmt.Println("GetUserByEmail", err.Error())
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"redirect_url": redirect_url,
				"error":        "email/password not match",
			})
			return
		}
		if u.Password == nil {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"redirect_url": redirect_url,
				"error":        "email/password not match",
			})
			return
		}
		if err := auth.CheckPasswordHash(loginForm.Password, *u.Password); err != nil {
			fmt.Println(err.Error())
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"redirect_url": redirect_url,
				"error":        "email/password not match",
			})
			return
		}

		tokens, err := auth.GenerateSSOToken(u.ID.String())
		if err != nil {
			fmt.Println(err.Error())
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"redirect_url": redirect_url,
				"error":        err.Error(),
			})
			return
		}
		fmt.Println(tokens)
		//Add redirect_url to query params
		parsedURL, err := url.Parse(redirect_url)
		if err != nil {
			fmt.Println("Error parsing URL:", err)
			return
		}
		queryParams := parsedURL.Query()
		queryParams.Add("access_token", tokens["access_token"])
		parsedURL.RawQuery = queryParams.Encode()
		c.Redirect(http.StatusOK, parsedURL.String())

		return
	})

	g.GET("/register", func(c *gin.Context) {
		redirect_url := c.Query("redirect_url")

		c.HTML(http.StatusOK, "register.html", gin.H{
			"redirect_url": redirect_url,
		})
	})

}
