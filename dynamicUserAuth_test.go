package dynamicUserAuth_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fino-digital/dynamicUserAuth"
	"github.com/labstack/echo"
)

func TestAuthMiddleware(t *testing.T) {
	// Init TestServer
	targetFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	})
	testServer := httptest.NewServer(targetFunc)
	host := "fino.digital"
	testServer.URL = "http://" + host
	defer testServer.Close()

	// build up a testStrategy
	testStrategy := dynamicUserAuth.Strategy{
		AuthorizeUser: func(c echo.Context) error {
			return nil
		},
	}

	// new middleware
	authMiddleware := dynamicUserAuth.NewAuthMiddleware(&dynamicUserAuth.DynamicUserAuth{Stragegies: dynamicUserAuth.Stragegies{host: testStrategy}})

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

func TestAllowedAddrSet(t *testing.T) {
	// Init TestServer
	targetFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	})
	testServer := httptest.NewServer(targetFunc)
	host := "fino.digital"
	testServer.URL = "http://" + host
	defer testServer.Close()

	// build up a testStrategy
	testStrategy := dynamicUserAuth.Strategy{
		AuthorizeUser: func(c echo.Context) error {
			return errors.New("it shouldn't pass this function")
		},
		AllowedAddrSet: map[string]struct{}{"192.0.2.1:1234": {}},
	}

	// new middleware
	authMiddleware := dynamicUserAuth.NewAuthMiddleware(&dynamicUserAuth.DynamicUserAuth{Stragegies: dynamicUserAuth.Stragegies{host: testStrategy}})

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
