package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/backsoul/gateway/configs"
	"github.com/backsoul/gateway/pkg/types"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var microservices = []types.Microservice{
	{
		Name: "auth",
		Url:  "https://auth.backsoul.xyz/api/v1/",
	},
	{
		Name: "posts",
		Url:  "http://bird/api/v1/",
	},
}

func main() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT"},
		AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.POST(":microservice/:method", func(c *gin.Context) {
		authorization := strings.Split(c.Request.Header["Authorization"][0], " ")[1]
		accessToken := strings.ReplaceAll(authorization, "Bearer ", "")

		token, err := jwt.ParseWithClaims(accessToken, &types.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(configs.Get("JWT_KEY")), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid authorization",
				"data":    err.Error(),
			})
			return
		}
		user := token.Claims.(*types.UserClaims)

		name := c.Param("microservice")
		method := c.Param("method")
		var microservice types.Microservice
		for _, m := range microservices {
			if m.Name == name {
				microservice = m
			}
		}

		bodyAsByteArray, _ := ioutil.ReadAll(c.Request.Body)
		jsonBody := string(bodyAsByteArray)

		payload := map[string]interface{}{"payload": jsonBody, "user": user}
		byts, _ := json.Marshal(payload)
		url := microservice.Url + method
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(byts))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"message": "error",
				"data":    err.Error(),
			})
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		jsonBodyResponse := string(body)
		c.JSON(http.StatusOK, gin.H{
			"result": jsonBodyResponse,
		})
	})
	r.Run()
}
