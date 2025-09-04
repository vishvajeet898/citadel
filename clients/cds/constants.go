package cdsClient

import (
	"github.com/Orange-Health/citadel/conf"
)

var (
	Config              = conf.GetConfig()
	CdsApiKey    string = Config.GetString("cds.api_key")
	CdsBaseUrl   string = Config.GetString("cds.base_url")
	CdsUserEmail string = Config.GetString("cds.user_email")
)

const (
	PANEL_DETAILS            = "PANEL_DETAILS"
	INVESTIGATION_DETAILS    = "INVESTIGATION_DETAILS"
	DEPARTMENT_MAPPING       = "DEPARTMENT_MAPPING"
	CONSTITUENTS_DATA        = "CONSTITUENTS_DATA"
	DEPENDENT_INVESTIGATIONS = "DEPENDENT_INVESTIGATIONS"
	VIALS                    = "VIALS"
	LAB_MASTERS              = "LAB_MASTERS"
	COLLECTION_SEQUENCE      = "COLLECTION_SEQUENCE"
	MASTER_TESTS             = "MASTER_TESTS"
	TEST_DEDUPLICATION       = "TEST_DEDUPLICATION"
	NRL_ENABLED_TESTS        = "NRL_ENABLED_TESTS"
)

var URL_MAP = map[string]string{
	TEST_DEDUPLICATION:       "/api/v1/test-deduplication",
	PANEL_DETAILS:            "/api/v1/investigation-master/investigation-panel-details",
	INVESTIGATION_DETAILS:    "/api/v1/investigation-master/investigation-details",
	DEPARTMENT_MAPPING:       "/api/v1/test-masters/department-mapping",
	CONSTITUENTS_DATA:        "/api/v1/test-masters/fetch",
	DEPENDENT_INVESTIGATIONS: "/api/v1/investigation-master/dependent-investigations",
	VIALS:                    "/api/v1/vials",
	LAB_MASTERS:              "/api/v1/lab-masters",
	COLLECTION_SEQUENCE:      "/api/v1/sequencer/fetch-collection-order",
	MASTER_TESTS:             "/api/v1/test-masters",
	NRL_ENABLED_TESTS:        "/api/v1/test-city-meta/nrl-enabled",
}
