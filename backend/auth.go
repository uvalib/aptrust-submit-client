package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type jwtClaims struct {
	ComputeID  string `json:"computeID"`
	CanApprove bool   `json:"canApprove"`
	jwt.RegisteredClaims
}

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
		log.Printf("INFO: using dev user ID: %s", computingID)
	}
	if computingID == "" {
		log.Printf("INFO: expected auth header not present in request. Not authorized.")
		c.Redirect(http.StatusFound, "/forbidden")
		return
	}

	// if not in dev mode check for membership in libra-admins
	canApprove := false
	if svc.Dev.user == "" {
		// Membership format: cn=group_name1;cn=group_name2;...
		membershipStr := c.GetHeader("member")
		if strings.Contains(membershipStr, svc.Group) {
			log.Printf("INFO: user %s is a member of %s and can approve submissions", computingID, svc.Group)
			canApprove = true
		} else {
			log.Printf("INFO: user %s is not a member of %s and cannot approve submissions", computingID, svc.Group)
		}
	} else {
		canApprove = svc.Dev.canApprove
		log.Printf("INFO: dev user approve setting: %t", canApprove)
	}

	signedStr, jwtErr := svc.mintUserJWT(computingID, canApprove)
	if jwtErr != nil {
		log.Printf("ERROR: unable to generate JWT for %s: %s", computingID, jwtErr.Error())
		c.Redirect(http.StatusFound, "/forbidden")
		return
	}

	// Set auth info in a cookie the client can read and pass along in future requests
	c.SetCookie("aptsubmit-client", signedStr, 5, "/", "", false, false)
	c.SetSameSite(http.SameSiteLaxMode)
	c.Redirect(http.StatusFound, "/signedin")
}

func (svc *serviceContext) mintUserJWT(computeID string, canApprove bool) (string, error) {
	expirationTime := time.Now().Add(8 * time.Hour)
	claims := jwtClaims{
		ComputeID:  computeID,
		CanApprove: canApprove,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Issuer:    "aptsubmit-client",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedStr, jwtErr := token.SignedString([]byte(svc.JWTKey))
	if jwtErr != nil {
		return "", jwtErr
	}
	return signedStr, nil
}
