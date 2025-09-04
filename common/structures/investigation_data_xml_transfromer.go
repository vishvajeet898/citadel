package structures

import "encoding/xml"

type AttuneInvestigationResults struct {
	XMLName              xml.Name                     `xml:"InvestigationResults"`
	InvestigationDetails []AttuneInvestigationDetails `xml:"InvestigationDetails"`
}

type AttuneInvestigationDetails struct {
	XMLName            xml.Name       `xml:"InvestigationDetails"`
	InvestigationName  string         `xml:"InvestigationName"`
	InvestigationID    string         `xml:"InvestigationID"`
	ReportStatus       string         `xml:"ReportStatus"`
	ClinicalHistory    string         `xml:"ClinicalHistory"`
	Gross              string         `xml:"Gross"`
	SampleName         string         `xml:"SampleName"`
	InCubPeriod        string         `xml:"InCubPeriod"`
	CultureStainType   string         `xml:"CultureStainType"`
	Colonycount        string         `xml:"Colonycount"`
	CultureReport      string         `xml:"CultureReport"`
	ResistanceDetected string         `xml:"ResistanceDetected"`
	StainDetails       []StainDetails `xml:"StainDetails"`
	OrganismIsolated   string         `xml:"OrganismIsolated"`
	OrganDetails       []OrganDetails `xml:"OrganDetails"`
}

type StainDetails struct {
	XMLName xml.Name `xml:"StainDetails"`
	Stain   []Stain  `xml:"Stain"`
}
type Stain struct {
	XMLName xml.Name `xml:"Stain"`
	Type    string   `xml:"Type,attr"`
	Result  string   `xml:"Result,attr"`
}

type OrganDetails struct {
	XMLName xml.Name `xml:"OrganDetails"`
	Organ   []Organ  `xml:"Organ"`
}

type Organ struct {
	XMLName     xml.Name `xml:"Organ"`
	Name        string   `xml:"Name,attr"`
	ColonyCount string   `xml:"ColonyCount,attr"`
	DrugName    string   `xml:"DrugName,attr"`
	NameSeq     string   `xml:"NameSeq,attr"`
	Family      string   `xml:"Family,attr"`
	FamilySeq   string   `xml:"FamilySeq,attr"`
	Sensitivity string   `xml:"Sensitivity,attr"`
	Level       string   `xml:"Level,attr"`
	Zone        string   `xml:"Zone,attr"`
}

type TransformedInvestigationResults struct {
	ReportStatus       string                   `json:"report_status"`
	SampleType         string                   `json:"sample_type"`
	IncubationPeriod   string                   `json:"incubation_period"`
	CultureStainResult string                   `json:"culture_stain_result"`
	CultureReport      string                   `json:"culture_report"`
	ColonyCount        string                   `json:"colony_count"`
	ResistanceDetected string                   `json:"resistance_detected"`
	StainResultDetails []StainResultDetail      `json:"stain_result_details"`
	GrowthData         []MicrorganismGrowthData `json:"growth_data"`
	ClinicalHistory    string                   `json:"clinical_history"`
	Gross              string                   `json:"gross"`
}

type StainResultDetail struct {
	Type   string `json:"type"`
	Result string `json:"result"`
}

type Sensitivity struct {
	Antimicrobial string `json:"antimicrobial"`
	Sensitivity   string `json:"sensitivity"`
	Level         string `json:"level"`
	MicValue      int    `json:"mic_value"`
	DisplayOrder  int    `json:"display_order"`
}

type MicrorganismGrowthData struct {
	Microorganism      string        `json:"microorganism"`
	ColonyCount        string        `json:"colony_count"`
	Family             string        `json:"family"`
	FamilyDisplayOrder *int          `json:"family_display_order"`
	Sensitivity        []Sensitivity `json:"sensitivity"`
}
