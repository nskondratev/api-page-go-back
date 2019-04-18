package handler

type responseEnvelope struct {
	Data interface{} `json:"data"`
}

type paginationResponseEnvelope struct {
	Data  interface{} `json:"data"`
	Total int         `json:"total"`
}

type errorResponseEnvelope struct {
	Error string `json:"error"`
}
