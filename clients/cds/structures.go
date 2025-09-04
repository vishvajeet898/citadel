package cdsClient

import (
	"github.com/Orange-Health/citadel/common/structures"
)

type VialsApiResponse struct {
	Count           uint                        `json:"count"`
	MasterVialTypes []structures.MasterVialType `json:"master_vial_types"`
}

type LabsApiResponse struct {
	Count uint             `json:"count"`
	Labs  []structures.Lab `json:"result"`
}

type MasterTestsResponse struct {
	Count uint                       `json:"count"`
	Tests []structures.CdsTestMaster `json:"result"`
}

type NrlEnabledMasterTestIdsResponse struct {
	Data []NrlEnabledMasterTestIdsStruct `json:"data"`
}

type NrlEnabledMasterTestIdsStruct struct {
	Id uint `json:"orangeTestId"`
}
