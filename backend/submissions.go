package main

import (
	"encoding/json"
	"fmt"
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

type storageOptions struct {
	ID       int64  `json:"id"`
	Value    string `json:"value"`
	Label    string `json:"label"`
	IsActive bool   `json:"active"`
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

	// First get the count of records matching the query
	lateralQ := " JOIN LATERAL (SELECT ss.* FROM submission_state ss WHERE s.identifier = ss.submission  ORDER BY ss.id DESC  LIMIT 1) ss ON true "
	cntSql := "select count(s.id) as total from submissions s inner join clients c on c.identifier = s.client" + lateralQ

	if includeAuto == false {
		// if approval email is set, this is not an auto-approval item
		conditions = append(conditions, "approval_email <> ''")
	} else {
		// include all submissions regardless of apporval_email
		conditions = append(conditions, "(approval_email = '' OR approval_email <> '')")
	}

	filter := svc.initFilter(filterStr)
	for _, fp := range filter {
		conditions = append(conditions, fmt.Sprintf("%s='%s'", fp.Field, fp.Value))
	}

	if q != "" {
		// it only makes sense to apply query to id and name; others will be done with exact matches
		conditions = append(conditions, "(identifier ~ ? or collection_name ~ ?)")
		countQ = svc.DB.Debug().Raw("select count(id) as total from submissions where ", q, q)
	} else {
		cntSql += " where " + strings.Join(conditions, " AND ")
		countQ = svc.DB.Debug().Raw(cntSql)
	}

	if err := countQ.Scan(&resp.Total).Error; err != nil {
		log.Printf("ERROR: unable to get total submissions: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// now apply the same conditions and get submission data for the page range specified
	sQ := "select s.id, s.identifier, s.storage, collection_name, s.created_at, c.name as client, c.approval_email as approval_email, ss.status as status"
	sQ += " from submissions s inner join clients c on c.identifier = s.client "
	sQ += lateralQ
	sQ += fmt.Sprintf(" where %s order by s.identifier asc offset %d limit %d", strings.Join(conditions, " AND "), startIdx, pageSize)
	qTX := svc.DB.Debug().Raw(sQ)

	if err := qTX.Scan(&resp.Hits).Error; err != nil {
		log.Printf("ERROR: get submissions failed: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

type filterParam struct {
	Field string
	Value string
}

func (svc *serviceContext) initFilter(filterStr string) []filterParam {
	log.Printf("INFO: raw filters [%s]", filterStr)
	out := make([]filterParam, 0)

	if filterStr != "" {
		//  Format: filters=["FIELD_NAME=value","FIELD2=value2"]}'
		var filterRequest []string
		if err := json.Unmarshal([]byte(filterStr), &filterRequest); err != nil {
			log.Printf("ERROR: invalid format for filter %s: %s", filterStr, err.Error())
			return out
		}
		for _, f := range filterRequest {
			bits := strings.Split(f, "=")
			field := strings.TrimSpace(bits[0])
			switch field {
			case "client":
				field = "c.name"
			case "status":
				field = "ss.status"
			}
			out = append(out, filterParam{Field: field, Value: strings.TrimSpace(bits[1])})
		}
	}

	return out
}

func (svc *serviceContext) approveSubmission(c *gin.Context) {
	submissionID := c.Param("id")
	computeID := getComputeID(c)
	log.Printf("INFO: user %s approves submission %s", computeID, submissionID)

	// TODO

	c.String(http.StatusNotImplemented, "not implemented")
}
