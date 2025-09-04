package structures

// CommonAPIResponse is the common API response structure.
// swagger:response commonAPIResponse
type CommonAPIResponse struct {
	// The message of the response.
	// example: "Fetched Successfully"
	Message string `json:"message,omitempty"`
	// The data of the response.
	Data interface{} `json:"data,omitempty"`
	// The error of the response.
	// example: "Error while fetching"
	Error string `json:"error,omitempty"`
	// Validation Errors
	// example: {"field": "error message"}
	ValidationErrors map[string]interface{} `json:"validation_errors,omitempty"`
}
