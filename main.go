package main

import (
	"fmt"
	"go/types"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/backsoul/gateway/configs"
	"github.com/backsoul/gateway/pkg/types"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type Microservice struct {
	Url  string
	Name string
}

var microservices = []Microservice{
	{
		Name: "auth",
		Url:  "https://auth.backsoul.xyz",
	},
	{
		Name: "posts",
		Url:  "https://posts.backsoul.xyz",
	},
}

func main() {
	r := gin.Default()
	r.GET(":microservice", func(c *gin.Context) {
		accessToken := strings.ReplaceAll(string(c.Get("Authorization")), "Bearer ", "")
		token, err := jwt.ParseWithClaims(accessToken, &types.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(configs.Get("JWT_KEY")), nil
		})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"message": "Errot parse jwt",
				"data":    err.Error(),
			})
		}
		claims := token.Claims.(*types)

		// get microservice
		name := c.Param("microservice")
		var microservice Microservice
		for _, m := range microservices {
			if m.Name == name {
				microservice = m
			}
		}

		resp, err := http.Get(microservice.Url)
		if err != nil {
			fmt.Printf("error making http request: %s\n", err)
			os.Exit(1)
		}

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("error reading")
				fmt.Printf(err.Error())
			}
			bodyString := string(bodyBytes)
			c.JSON(http.StatusOK, gin.H{
				"message": bodyString,
			})
		} else {

		}
	})
	r.Run()
}
