package restUserAuth

import (
	"net/http"
	"strings"

	"github.com/fino-digital/dynamicUserAuth"
	"github.com/labstack/echo"
)

// AuthRest is the instance for Auth-REST-impl
type AuthRest struct {
	UserAuth    dynamicUserAuth.DynamicUserAuth
	IgnoreRoute string
}

// Handle handles all functions dynamic by host
func (authRest *AuthRest) Handle(context echo.Context) error {
	// find correct strategy
	host := context.Request().Host
	strategy, ok := authRest.UserAuth.Stragegies[host]
	if !ok {
		return context.JSON(http.StatusMethodNotAllowed, "Can't find host: "+host)
	}

	// find correct function
	path := strings.Replace(context.Request().URL.String(), authRest.IgnoreRoute, "", 1)
	path = strings.Replace(path, "/", "", 1)
	strategyFunc, ok := strategy.Functions[path]
	if !ok {
		return context.JSON(http.StatusMethodNotAllowed, "Can't find route: "+path)
	}

	// get body
	requestMap := new(map[string]interface{})
	if err := context.Bind(requestMap); err != nil {
		return context.JSON(http.StatusMethodNotAllowed, err)
	}

	// call resolve of function
	returnInterf, err := strategyFunc.Resolve(context, (*requestMap))
	if err != nil {
		return context.JSON(http.StatusMethodNotAllowed, err)
	}

	return context.JSON(http.StatusOK, returnInterf)
}
