package structures

type CommonError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

type ResponseError struct {
	Error string `json:"error"`
}
