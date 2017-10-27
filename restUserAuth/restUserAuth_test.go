package restUserAuth_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"

	"github.com/fino-digital/dynamicUserAuth"
	"github.com/fino-digital/dynamicUserAuth/restUserAuth"
)

func TestRestUserAuthHostCheck(t *testing.T) {
	testData := []struct {
		Host           string
		Route          string
		Body           map[string]interface{}
		ExpectedResult int
	}{
		// Wrong host, correct route
		{
			Host:           "wrongHost",
			Route:          "/ignoreme/NewUser",
			ExpectedResult: restUserAuth.StatusNoHost,
		},
		// Correct host, but wrong route
		{
			Host:           "fino.digital",
			Route:          "/WrongRoute",
			ExpectedResult: http.StatusNotFound,
		},
		// Correct host, but wrong route
		{
			Host:           "fino.digital",
			Route:          "/ignoreme/WrongRoute",
			ExpectedResult: restUserAuth.StatusNoFunction,
		},
		// correct host, correct route, but without body:
		{
			Host:           "fino.digital",
			Route:          "/ignoreme/NewUser",
			ExpectedResult: http.StatusBadRequest,
		},
		// correct host, correct route
		{
			Host:           "fino.digital",
			Route:          "/ignoreme/NewUser",
			ExpectedResult: http.StatusOK,
			Body:           map[string]interface{}{"body": "correct"},
		},
		// correct host, correct route, but wrong body
		// CURRENTLY NOT IMPLEMENTED
		{
			Host:           "fino.digital",
			Route:          "/ignoreme/NewUser",
			Body:           map[string]interface{}{"body": "wrongBody"},
			ExpectedResult: http.StatusOK,
		},
	}

	for testDataIndex, data := range testData {
		strategies := map[string]dynamicUserAuth.Strategy{"fino.digital": dynamicUserAuth.Strategy{
			Functions: map[string]dynamicUserAuth.StrategyFunction{
				"NewUser": dynamicUserAuth.StrategyFunction{
					Resolve: func(c echo.Context, requestMap map[string]interface{}) (interface{}, error) {
						return "hello", nil
					},
				},
			},
		}}

		router := echo.New()
		rest := restUserAuth.AuthRest{UserAuth: dynamicUserAuth.DynamicUserAuth{Stragegies: strategies}}
		router.Any("/ignoreme"+restUserAuth.ParamFunction, rest.Handle)

		// read body if there is something
		var reader io.Reader
		if data.Body != nil {
			bodyByte, _ := json.Marshal(data.Body)
			reader = bytes.NewReader(bodyByte)
		}

		// build router and requet
		req := httptest.NewRequest(echo.POST, data.Route, reader)
		req.Host = data.Host
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Result().StatusCode != data.ExpectedResult {
			t.Errorf("[%d] Actual StatusCode: %d but expected: %d; With body: %s",
				testDataIndex, rec.Result().StatusCode, data.ExpectedResult, rec.Body.String())
		}
	}

}
