package responder

func commonHeaders() map[string]string {
	return map[string]string{
		"Content-Type":                "application/json",
		"Access-Control-Allow-Origin": "*",
	}
}

func Response(code int) interface{} {
	res := &struct {
		StatusCode int `json:"statusCode"`
	}{StatusCode: code}

	return res
}
