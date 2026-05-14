package main

import (
	"flag"
	"log"
)

type devConfig struct {
	user       string
	canApprove bool
	fakeBus    bool
}

type dbConfig struct {
	host string
	port int
	user string
	pass string
	name string
}

type configData struct {
	port            int
	aptrust         string
	jwtKey          string
	busName         string
	eventSourceName string
	database        dbConfig
	group           string
	dev             devConfig
}

func getConfiguration() *configData {
	var config configData
	flag.IntVar(&config.port, "port", 8080, "Port to offer service on")
	flag.StringVar(&config.jwtKey, "jwtkey", "", "JWT signature key")
	flag.StringVar(&config.group, "group", "lib-aptrust-submit-approve-dev", "Grouper group name for submitters")
	flag.StringVar(&config.aptrust, "aptrust", "", "APTrust repo url")

	// DB connection params
	flag.StringVar(&config.database.host, "dbhost", "", "Database host")
	flag.IntVar(&config.database.port, "dbport", 5432, "Database port")
	flag.StringVar(&config.database.name, "dbname", "virgo4", "Database name")
	flag.StringVar(&config.database.user, "dbuser", "v4user", "Database user")
	flag.StringVar(&config.database.pass, "dbpass", "pass", "Database password")

	// event bus
	flag.StringVar(&config.busName, "busname", "", "Event bus name")
	flag.StringVar(&config.eventSourceName, "eventsrc", "", "Event source name")

	// dev mode
	flag.StringVar(&config.dev.user, "devuser", "", "Authorized computing id for dev")
	flag.BoolVar(&config.dev.canApprove, "devapprove", false, "Can dev user approve submissions")
	flag.BoolVar(&config.dev.fakeBus, "devbus", false, "bus dev mode (no events sent out)")

	flag.Parse()

	if config.aptrust == "" {
		log.Fatal("Parameter aptrust is required")
	}
	if config.jwtKey == "" {
		log.Fatal("Parameter jwtkey is required")
	}
	if config.busName == "" {
		log.Fatal("Parameter busname is required")
	}
	if config.eventSourceName == "" {
		log.Fatal("Parameter eventsrc is required")
	}
	if config.database.host == "" {
		log.Fatal("Parameter dbhost is required")
	}
	if config.database.name == "" {
		log.Fatal("Parameter dbname is required")
	}
	if config.database.user == "" {
		log.Fatal("Parameter dbuser is required")
	}
	if config.database.pass == "" {
		log.Fatal("Parameter dbpass is required")
	}

	log.Printf("[CONFIG] port            = [%d]", config.port)
	log.Printf("[CONFIG] aptrust         = [%s]", config.aptrust)
	log.Printf("[CONFIG] group           = [%s]", config.group)
	log.Printf("[CONFIG] dbhost          = [%s]", config.database.host)
	log.Printf("[CONFIG] dbport          = [%d]", config.database.port)
	log.Printf("[CONFIG] dbname          = [%s]", config.database.name)
	log.Printf("[CONFIG] dbuser          = [%s]", config.database.user)
	log.Printf("[CONFIG] eventsrc        = [%s]", config.eventSourceName)
	log.Printf("[CONFIG] busname         = [%s]", config.busName)

	if config.dev.user != "" {
		log.Printf("[CONFIG] devuser         = [%s]", config.dev.user)
		log.Printf("[CONFIG] devapprove      = [%t]", config.dev.canApprove)
	}
	if config.dev.fakeBus {
		log.Printf("[CONFIG] ** dev mode bus - event publishing is disabled **")
	}

	return &config
}
