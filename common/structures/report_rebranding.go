package structures

type MediaResizeRequest struct {
	Files []MediaResizeFile `json:"files"`
}

type MediaResizeResponse struct {
	Files []MediaResizeFile `json:"files"`
}

type MediaResizeFile struct {
	Id            uint   `json:"id"`
	FileType      string `json:"fileType,omitempty"`
	FileUrl       string `json:"fileUrl"`
	FileReference string `json:"fileReference,omitempty"`
}

type CobrandingRequest struct {
	ReportUrl         string `json:"reportUrl"`
	CobrandedImageUrl string `json:"cobrandedImageUrl"`
	HeaderLength      uint   `json:"headerLength"`
}

type CobrandingResponse struct {
	ReportUrl  string `json:"reportUrl"`
	ReportData string `json:"reportData"`
}
