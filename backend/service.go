package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/uvalib/aptrust-submit-bus-definitions/uvaaptsbus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type eventContext struct {
	DevMode     bool
	BusName     string
	EventSource string
	Bus         uvaaptsbus.UvaBus
}

type serviceContext struct {
	Version string
	JWTKey  string
	Group   string
	DB      *gorm.DB
	Events  eventContext
	Dev     devConfig
}

func initializeService(version string, cfg *configData) *serviceContext {
	ctx := serviceContext{
		Version: version,
		JWTKey:  cfg.jwtKey,
		Group:   cfg.group,
		Dev:     cfg.dev,
	}

	log.Printf("INFO: connecting to database...")
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d",
		cfg.database.user, cfg.database.pass, cfg.database.name, cfg.database.host, cfg.database.port)
	gdb, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	ctx.DB = gdb
	log.Printf("INFO: db connect success")

	ctx.Events.DevMode = cfg.dev.fakeBus
	ctx.Events.BusName = cfg.busName
	ctx.Events.EventSource = cfg.eventSourceName
	if cfg.dev.fakeBus == false {
		log.Printf("INFO: configure event bus [%s] with source [%s]", cfg.busName, cfg.eventSourceName)
		busCfg := uvaaptsbus.UvaBusConfig{BusName: cfg.busName, Source: cfg.eventSourceName, Log: nil}
		bus, err := uvaaptsbus.NewUvaBus(busCfg)
		if err != nil {
			log.Fatalf("unable to init event bus with name %s and source %s: %s", cfg.busName, cfg.eventSourceName, err.Error())
		}
		ctx.Events.Bus = bus
	}

	return &ctx
}

func (svc *serviceContext) getVersion(c *gin.Context) {
	vMap := svc.lookupVersion()
	c.JSON(http.StatusOK, vMap)
}

func (svc *serviceContext) lookupVersion() map[string]string {
	build := "unknown"
	// working directory is the bin directory, and build tag is in the root
	files, _ := filepath.Glob("../buildtag.*")
	if len(files) == 1 {
		build = strings.Replace(files[0], "../buildtag.", "", 1)
	}

	vMap := make(map[string]string)
	vMap["version"] = svc.Version
	vMap["build"] = build
	return vMap
}

func (svc *serviceContext) healthCheck(c *gin.Context) {
	type hcResp struct {
		Healthy bool   `json:"healthy"`
		Message string `json:"message,omitempty"`
	}
	hcMap := make(map[string]hcResp)
	hcMap["aptsubmit-client"] = hcResp{Healthy: true}

	c.JSON(http.StatusOK, hcMap)
}

func (svc *serviceContext) getConfig(c *gin.Context) {
	verInfo := svc.lookupVersion()
	ver := fmt.Sprintf("v%s-%s", verInfo["version"], verInfo["build"])
	resp := struct {
		Version            string           `json:"version"`
		SubmissionStatuses []string         `json:"submissionStatuses"`
		Clients            []client         `json:"clients"`
		StorageOptions     []storageOptions `json:"storageOptions"`
	}{
		Version:            ver,
		SubmissionStatuses: []string{"abandoned", "building", "complete", "error", "incomplete", "pending-approval", "pending-ingest", "registered", "submitting", "validating"},
	}

	var clients []client
	if err := svc.DB.Find(&clients).Error; err != nil {
		log.Printf("ERROR: unable to load clients info: %s", err.Error())
	} else {
		resp.Clients = clients
	}
	var storage []storageOptions
	if err := svc.DB.Where("is_active=?", true).Find(&storage).Error; err != nil {
		log.Printf("ERROR: unable to load storage options: %s", err.Error())
	} else {
		resp.StorageOptions = storage
	}
	c.JSON(http.StatusOK, resp)
}

func (svc *serviceContext) logClientError(c *gin.Context) {
	agent := c.Request.UserAgent()
	body, _ := io.ReadAll(c.Request.Body)
	build := "unknown"
	files, _ := filepath.Glob("../buildtag.*")
	if len(files) == 1 {
		build = strings.Replace(files[0], "../buildtag.", "", 1)
	}
	if strings.Contains(agent, "Googlebot") || strings.Contains(agent, "Applebot") {
		log.Printf("CLIENT (bot) WARN (Version %s.%s): %s", svc.Version, build, body)
	} else {
		log.Printf("CLIENT ERROR (Version %s.%s): %s", svc.Version, build, body)
	}
}

func (svc *serviceContext) publishWorkflowEvent(eventName string, clientId string, submissionId string, extra any) error {
	extraBytes, err := json.Marshal(extra)
	if err != nil {
		return err
	}
	pl := uvaaptsbus.UvaWorkflowEvent{SubmissionId: submissionId, Extra: fmt.Sprintf("%s", extraBytes)}
	detail, err := pl.Serialize()
	if err != nil {
		return err
	}

	event := uvaaptsbus.UvaBusEvent{
		EventName: eventName,
		ClientId:  clientId,
		Detail:    detail,
	}

	if svc.Dev.fakeBus {
		log.Printf("INFO: using fake dev bus to publish %+v", event)
	} else {
		log.Printf("INFO: publishing %+v", event)
		if err := svc.Events.Bus.PublishEvent(&event); err != nil {
			return err
		}
	}
	return nil
}
