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
	"github.com/uvalib/aptrust-submit-bus-definitions/uvaaptsbus"
	"gorm.io/gorm"
)

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

type client struct {
	ID             int64  `json:"id"`
	Identifier     string `json:"identifier"`
	Name           string `json:"name"`
	DefaultStorage string `json:"defaultStorage"`
	ApprovalEmail  string `json:"approvalEmail"`
}

type submissionState struct {
	ID         int64     `json:"id"`
	Submission string    `json:"-"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
}

type submissionFailure struct {
	ID         int64     `json:"id"`
	Submission string    `json:"-"`
	Failure    string    `json:"failure"`
	CreatedAt  time.Time `json:"createdAt"`
}

type file struct {
	ID         int64     `json:"id"`
	Name       string    `json:"fileName"`
	Hash       string    `json:"hash"`
	BagName    string    `json:"bagName" gorm:"primaryKey"`
	FileSize   int64     `json:"fileSize"`
	Submission string    `json:"-" gorm:"primaryKey"`
	CreatedAt  time.Time `json:"createdAt"`
}

type aptFile struct {
	ID         int64      `json:"id"`
	FileName   string     `json:"fileName"`
	Hash       string     `json:"hash"`
	BagName    string     `json:"bagName"`
	FileSize   int64      `json:"fileSize"`
	APTAddedAt *time.Time `json:"aptAddedAt" gorm:"column:apt_added_at"`
	CreatedAt  time.Time  `json:"createdAt"`
}

type submissionConflict struct {
	ID                int64     `json:"id"`
	Submission        string    `json:"-"`
	NewFileID         int64     `json:"-" gorm:"column:new_file"`                          // key into files table
	NewFile           file      `json:"newFile" gorm:"foreignKey:ID;references:NewFileID"` // file data
	Basis             string    `json:"basis"`                                             // source of conflict; aptrust or local
	ConflictingFileID int64     `json:"-" gorm:"column:conflicting_file"`                  // for aptrust basis, key into apt_files; for local key info files table
	LocalConflict     *file     `json:"localConflict,omitempty" gorm:"-"`
	APTConflict       *aptFile  `json:"aptConflict,omitempty" gorm:"-"`
	Ignored           bool      `json:"ignored"`
	CreatedAt         time.Time `json:"createdAt"`
}

type approval struct {
	ID         int64     `json:"id"`
	Submission string    `json:"-"`
	Who        string    `json:"statuwhos"`
	CreatedAt  time.Time `json:"createdAt"`
}

type submission struct {
	ID             int64                 `json:"id"`
	Identifier     string                `json:"identifier"`
	Storage        string                `json:"storage"`
	CollectionName string                `json:"collectionName"`
	Client         string                `json:"-"`
	ClientInfo     client                `json:"client" gorm:"foreignKey:Identifier;references:Client"`
	BagCount       uint                  `json:"bagCount"`
	FileCount      uint                  `json:"fileCount"`
	TotalFileSize  int64                 `json:"totalFileSize"`
	CreatedAt      time.Time             `json:"createdAt"`
	Status         []submissionState     `json:"status" gorm:"foreignKey:Submission;references:Identifier"`
	Failures       []submissionFailure   `json:"failures" gorm:"foreignKey:Submission;references:Identifier"`
	Conflicts      []*submissionConflict `json:"conflicts" gorm:"foreignKey:Submission;references:Identifier"`
	Approval       *approval             `json:"approval,omitempty" gorm:"foreignKey:Submission;references:Identifier"`
}

type bagState struct {
	ID         int64     `json:"id"`
	Submission string    `json:"-" gorm:"primaryKey"`
	BagName    string    `json:"-" gorm:"primaryKey"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
}

type bag struct {
	ID         int64      `json:"id"`
	Submission string     `json:"-"`
	BagName    string     `json:"name"`
	Status     []bagState `json:"status" gorm:"foreignKey:BagName,Submission;references:BagName,Submission"`
	Files      []file     `json:"files" gorm:"foreignKey:BagName,Submission;references:BagName,Submission"`
	CreatedAt  time.Time  `json:"createdAt"`
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
	lateralQ := " JOIN LATERAL (SELECT ss.* FROM submission_states ss WHERE s.identifier = ss.submission  ORDER BY ss.id DESC  LIMIT 1) ss ON true "
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
	sQ += fmt.Sprintf(" where %s order by s.created_at desc offset %d limit %d", strings.Join(conditions, " AND "), startIdx, pageSize)
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

func (svc *serviceContext) getSubmissionDetail(c *gin.Context) {
	submissionID := c.Param("id")
	// submissionID := "sid-d7loa41aq8ss73jdp53g" // TS
	// submissionID := "sid-d7lo3ejsp7mc73ncli60" // approved
	// submissionID := "sid-d7ockedkpgss73ktindg" // confiict
	computeID := getComputeID(c)
	log.Printf("INFO: user %s requests details for submission %s", computeID, submissionID)

	// first load main submission data..
	var sub submission
	if err := svc.DB.Preload("ClientInfo").Preload("Status", func(db *gorm.DB) *gorm.DB {
		return db.Order("submission_states.created_at DESC")
	}).Preload("Approval").Preload("Failures").Preload("Conflicts").Preload("Conflicts.NewFile").
		Where("submissions.identifier=?", submissionID).First(&sub).Error; err != nil {
		log.Printf("ERROR: unable to get submission %s detail: %s", submissionID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// next, get bag / file summary
	if err := svc.DB.Raw("select count(*) bags_count from bags where submission=?", submissionID).Scan(&sub.BagCount).Error; err != nil {
		log.Printf("ERROR: unable to get bag count for submission %s: %s", submissionID, err.Error())
	}
	var fileSummary struct {
		FileCount     uint
		TotalFileSize int64
	}
	if err := svc.DB.Raw("select count(id) file_count, sum(file_size) total_file_size from files where submission=? and name!=? and name!=?",
		submissionID, "aptrust-description.txt", "aptrust-title.txt").
		Scan(&fileSummary).Error; err != nil {
		log.Printf("ERROR: unable to get file summary for submission %s: %s", submissionID, err.Error())
	} else {
		sub.FileCount = fileSummary.FileCount
		sub.TotalFileSize = fileSummary.TotalFileSize
	}

	// if there are conflict, resove the files involved:
	//    basis 'aptrust' uses the apt_files table
	//    basis 'local' is a conflict within the bag uses the files table
	aptIDs := make([]int64, 0)
	localIDs := make([]int64, 0)
	for _, c := range sub.Conflicts {
		if c.Basis == "aptrust" {
			aptIDs = append(aptIDs, c.ConflictingFileID)
		} else {
			localIDs = append(localIDs, c.ConflictingFileID)
		}
	}

	if len(aptIDs) > 0 {
		var aptFiles []aptFile
		if err := svc.DB.Where("id in ?", aptIDs).Find(&aptFiles).Error; err != nil {
			log.Printf("ERROR: unable to get conflicting aptfiles data for submission %s: %s", submissionID, err.Error())
		} else {
			for _, aptF := range aptFiles {
				for _, c := range sub.Conflicts {
					if c.ConflictingFileID == aptF.ID {
						c.APTConflict = &aptF
						break
					}
				}
			}
		}
	}
	if len(localIDs) > 0 {
		var localFiles []file
		if err := svc.DB.Where("id in ?", localIDs).Find(&localFiles).Error; err != nil {
			log.Printf("ERROR: unable to get conflicting local files data for submission %s: %s", submissionID, err.Error())
		} else {
			for _, localF := range localFiles {
				for _, c := range sub.Conflicts {
					if c.ConflictingFileID == localF.ID {
						c.LocalConflict = &localF
						break
					}
				}
			}
		}
	}

	c.JSON(http.StatusOK, sub)
}

func (svc *serviceContext) getSubmissionBags(c *gin.Context) {
	submissionID := c.Param("id")
	computeID := getComputeID(c)
	log.Printf("INFO: user %s requests bag for submission %s", computeID, submissionID)

	var bags []bag
	if err := svc.DB.Debug().Preload("Files", "name != ? and name !=?", "aptrust-description.txt", "aptrust-title.txt").
		Preload("Status", func(db *gorm.DB) *gorm.DB {
			return db.Order("bag_states.created_at DESC")
		}).Where("bags.submission=?", submissionID).
		Find(&bags).Error; err != nil {
		log.Printf("ERROR: unable to get submission %s bags: %s", submissionID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, bags)
}

func (svc *serviceContext) approveSubmission(c *gin.Context) {
	submissionID := c.Param("id")
	computeID := getComputeID(c)
	var req struct {
		Storage string `json:"storage"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("INFO: bad approve request: %s", err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	var sub submission
	if err := svc.DB.Preload("ClientInfo").Where("identifier=?", submissionID).First(&sub).Error; err != nil {
		log.Printf("INFO: unable to load submission %s: %s", submissionID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	extra := struct {
		ComputeID string `json:"computeID"`
		Storage   string `json:"storage"`
	}{
		ComputeID: computeID,
		Storage:   req.Storage,
	}

	log.Printf("INFO: user %s approves submission %s with storage %s", computeID, submissionID, req.Storage)
	if err := svc.publishWorkflowEvent(uvaaptsbus.EventSubmissionApproved, sub.ClientInfo.Identifier, sub.Identifier, extra); err != nil {
		log.Printf("ERROR: unable to publsh approve event for submission %s: %s", submissionID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, "approved")
}

func (svc *serviceContext) cancelSubmission(c *gin.Context) {
	submissionID := c.Param("id")
	computeID := getComputeID(c)

	var sub submission
	if err := svc.DB.Preload("ClientInfo").Where("identifier=?", submissionID).First(&sub).Error; err != nil {
		log.Printf("INFO: unable to load submission %s: %s", submissionID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("INFO: user %s cancels submission %s", computeID, submissionID)
	if err := svc.publishWorkflowEvent(uvaaptsbus.EventSubmissionAbandoned, sub.ClientInfo.Identifier, sub.Identifier, computeID); err != nil {
		log.Printf("ERROR: unable to publsh abandon event for submission %s: %s", submissionID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, "canceled")
}
