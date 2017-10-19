package dynamicUserAuth_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"gitlab.com/fino/finnbroker/core/auth"
	"gitlab.com/fino/finnbroker/finnMiddleware"
)

func TestAuthMiddleware(t *testing.T) {
	// Init TestServer
	targetFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	})
	testServer := httptest.NewServer(targetFunc)
	host := "abc.de"
	testServer.URL = "http://" + host
	defer testServer.Close()

	// build up a testStrategy
	testStrategy := auth.Strategy{
		AuthorizeUser: func(c echo.Context) error {
			return nil
		},
	}

	// new middleware
	authMiddleware := finnMiddleware.NewAuthMiddleware(&auth.FinnAuth{Stragegies: auth.Stragegies{host: &testStrategy}})

	// build request
	router := echo.New()
	request := httptest.NewRequest(echo.GET, testServer.URL, nil)
	rec := httptest.NewRecorder()
	context := router.NewContext(request, rec)

	// TEST
	err := authMiddleware.Handle(echo.WrapHandler(targetFunc))(context)
	if err != nil {
		t.Error(err)
	}
}
