package main

import (
	"fmt"
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
	DB      *gorm.DB
	Events  eventContext
	Dev     devConfig
}

func initializeService(version string, cfg *configData) *serviceContext {
	ctx := serviceContext{
		Version: version,
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
		Version string `json:"version"`
	}{
		Version: ver,
	}
	c.JSON(http.StatusOK, resp)
}
