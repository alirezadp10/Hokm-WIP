package tests

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
)

func ApiCall(e *echo.Echo, method, url string, routeParams, headers, body map[string]string) (echo.Context, *httptest.ResponseRecorder) {
	marshaledBody, _ := json.Marshal(body)
	req := httptest.NewRequest(method, url, bytes.NewReader(marshaledBody))
	rec := httptest.NewRecorder()
	for headerKey, headerValue := range headers {
		req.Header.Set(headerKey, headerValue)
	}
	c := e.NewContext(req, rec)
	c.SetPath(url)
	for paramKey, paramValue := range routeParams {
		c.SetParamNames(paramKey)
		c.SetParamValues(paramValue)
	}
	return c, rec
}

func convertByteToStringArray(data []byte) []string {
	var result []string
	_ = json.Unmarshal(data, &result)
	return result
}
