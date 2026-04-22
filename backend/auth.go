package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (svc *serviceContext) authenticate(c *gin.Context) {
	log.Printf("INFO: authenticate request is checking headers")
	log.Printf("Dump all request headers ==================================")
	for name, values := range c.Request.Header {
		for _, value := range values {
			log.Printf("%s=%s\n", name, value)
		}
	}
	log.Printf("END header dump ===========================================")

	computingID := c.GetHeader("remote_user")
	if svc.Dev.user != "" {
		computingID = svc.Dev.user
		log.Printf("INFO: using dev auth user ID: %s", computingID)
	}
	if computingID == "" {
		log.Printf("ERROR: expected auth header not present in request. Not authorized.")
		c.Redirect(http.StatusFound, "/forbidden")
		return
	}
}
