package s3wrapperClient

type PublicUrlRequestBody struct {
	FileReference string `json:"fileReference"`
}

type OrderFilePublicUrlData struct {
	MediaUrl string `json:"mediaUrl"`
}
