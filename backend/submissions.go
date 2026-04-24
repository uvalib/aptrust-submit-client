package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type client struct {
	ID             int64  `json:"id"`
	Identifier     string `json:"identifier"`
	Name           string `json:"name"`
	DefaultStorage string `json:"defaultStorage"`
	ApprovalEmail  string `json:"approvalEmail"`
}

type searchHit struct {
	ID             int64     `json:"id"`
	Identifier     string    `json:"identifier"`
	Storage        string    `json:"storage"`
	CollectionName string    `json:"collectionName"`
	Client         string    `json:"client"`
	ApprovalEmail  string    `json:"approvalEmail"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"createdAt"`
}

type searchResponse struct {
	Total int64       `json:"total"`
	Hits  []searchHit `json:"hits"`
}

func (svc *serviceContext) getSubmissions(c *gin.Context) {
	computeID := getComputeID(c)
	q := strings.TrimSpace(c.Query("q"))
	filterStr := c.Query("filters")
	includeAuto, _ := strconv.ParseBool(c.Query("includeauto"))
	startIdx, _ := strconv.ParseInt(c.Query("start"), 10, 64)
	pageSize, _ := strconv.ParseInt(c.Query("limit"), 10, 64)
	if pageSize == 0 {
		pageSize = 30
	}
	log.Printf("INFO: user %s requests submissions starting from  %d limit %d with query [%s] and filters [%s]",
		computeID, startIdx, pageSize, q, filterStr)

	resp := searchResponse{}
	var countQ *gorm.DB
	conditions := make([]string, 0)

	// First get the count based on search q and include auto accepted
	cntSql := "select count(s.id) as total from submissions s inner join clients c on c.identifier = s.client"

	if includeAuto == false {
		// if approval email is set, this is not an auto-approval item
		conditions = append(conditions, "approval_email <> ''")
	} else {
		// include all submissions regardless of apporval_email
		conditions = append(conditions, "(approval_email = '' OR approval_email <> '')")
	}

	if q != "" {
		// it only makes sense to apply query to id and name; others will be done with exact matches
		conditions = append(conditions, "(identifier ~ ? or collection_name ~ ?)")
		countQ = svc.DB.Raw("select count(id) as total from submissions where ", q, q)
	} else {
		cntSql += " where " + strings.Join(conditions, " AND ")
		countQ = svc.DB.Raw(cntSql)
	}

	if err := countQ.Scan(&resp.Total).Error; err != nil {
		log.Printf("ERROR: unable to get total submissions: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// now apply the same conditions and get submission data for the page range specified
	statusSubQ := " (select ss.status from submission_state ss where ss.submission = s.identifier order by ss.created_at desc limit 1) as status "
	sQ := "select s.id, s.identifier, s.storage, collection_name, s.created_at, c.name as client, c.approval_email as approval_email, "
	sQ += statusSubQ
	sQ += "from submissions s inner join clients c on c.identifier = s.client "
	sQ += " where " + strings.Join(conditions, " AND ")
	qTX := svc.DB.Debug().Raw(sQ)
	qTX.Order("s.identifier asc")
	qTX.Offset(int(startIdx))
	qTX.Limit(int(pageSize))

	if err := qTX.Scan(&resp.Hits).Error; err != nil {
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
