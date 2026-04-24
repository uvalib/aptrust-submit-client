package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type searchHit struct {
	ID             int64     `json:"id"`
	Identifier     string    `json:"identifier"`
	Storage        string    `json:"storage"`
	CollectionName string    `json:"collectionName"`
	Client         string    `json:"client"`
	ApprivalEmail  string    `json:"approvalEmail"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"createdAt"`
}

type searchResponse struct {
	Total int64       `json:"total"`
	Hits  []searchHit `json:"hits"`
}

func (svc *serviceContext) getSubmissions(c *gin.Context) {
	q := strings.TrimSpace(c.Query("q"))
	filterStr := c.Query("filters")
	computeID := getComputeID(c)
	startIdx, _ := strconv.ParseInt(c.Query("start"), 10, 64)
	pageSize, _ := strconv.ParseInt(c.Query("limit"), 10, 64)
	if pageSize == 0 {
		pageSize = 30
	}
	log.Printf("INFO: user %s requests submissions starting from  %d limit %d with query [%s] and filters [%s]",
		computeID, startIdx, pageSize, q, filterStr)

	resp := searchResponse{}
	if err := svc.DB.Raw("select count(id) as total from submissions").Scan(&resp.Total).Error; err != nil {
		log.Printf("ERROR: unable to get total submissions: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	statusSubQ := " (select ss.status from submission_state ss where ss.submission = s.identifier order by ss.created_at desc limit 1) as status "
	sQ := "select s.id, s.identifier, s.storage, collection_name, s.created_at, c.name as client, c.approval_email as approval_email, "
	sQ += statusSubQ
	sQ += "from submissions s inner join clients c on c.identifier = s.client "
	sQ += fmt.Sprintf("order by s.identifier asc offset %d limit %d", startIdx, pageSize)

	if err := svc.DB.Debug().Raw(sQ).Scan(&resp.Hits).Error; err != nil {
		log.Printf("ERROR: get submissions failed: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (svc *serviceContext) approveSubmission(c *gin.Context) {
	submissionID := c.Param("id")
	computeID := getComputeID(c)
	log.Printf("INFO: user %s approves submission %s", computeID, submissionID)

	// TODO

	c.String(http.StatusNotImplemented, "not implemented")
}
