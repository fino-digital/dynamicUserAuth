package restUserAuth_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
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
		ExpectedResult bool
	}{
		// Wrong host, correct route
		{
			Host:           "wrongHost",
			Route:          "NewUser",
			ExpectedResult: false,
		},
		// Correct host, but wrong host
		{
			Host:           "fino.digital",
			Route:          "WrongRoute",
			ExpectedResult: false,
		},
		// correct host, correct route, but without body:
		{
			Host:           "fino.digital",
			Route:          "NewUser",
			ExpectedResult: false,
		},
		// correct host, correct route, but wrong body
		// CURRENTLY NOT IMPLEMENTED
		{
			Host:           "fino.digital",
			Route:          "NewUser",
			Body:           map[string]interface{}{"body": "wrongBody"},
			ExpectedResult: true,
		},
	}

	for _, data := range testData {
		strategies := map[string]dynamicUserAuth.Strategy{"fino.digital": dynamicUserAuth.Strategy{
			Functions: map[string]dynamicUserAuth.StrategyFunction{
				"NewUser": dynamicUserAuth.StrategyFunction{
					Resolve: func(c echo.Context, requestMap map[string]interface{}) (interface{}, error) {
						return "hello", nil
					},
				},
			},
		}}
		rest := restUserAuth.AuthRest{UserAuth: dynamicUserAuth.DynamicUserAuth{Stragegies: strategies}}

		// read body if there is something
		var reader io.Reader
		if data.Body != nil {
			bodyByte, _ := json.Marshal(data.Body)
			reader = bytes.NewReader(bodyByte)
		}

		// build router and requet
		router := echo.New()
		req := httptest.NewRequest(echo.POST, "/NewUser", reader)
		req.Host = data.Host
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		context := router.NewContext(req, httptest.NewRecorder())

		// call handler
		err := rest.Handle(context)
		if (err != nil) == data.ExpectedResult {
			log.Println(err)
			t.Fail()
		}
	}

}
