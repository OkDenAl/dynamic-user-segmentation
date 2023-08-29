package tests

type succesResponse struct {
	Message string `json:"message"`
}

type errorResponse struct {
	Data  any    `json:"data"`
	Error string `json:"error"`
}

type dataResponse struct {
	Data  any    `json:"data"`
	Error string `json:"error"`
}

type linkResponse struct {
	Link string `json:"link"`
}
