package restUserAuth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/fino-digital/dynamicUserAuth"
	"github.com/labstack/echo"
)

// AuthRest is the instance for Auth-REST-impl
type AuthRest struct {
	UserAuth dynamicUserAuth.DynamicUserAuth
}

// Handle handles all functions dynamic by host
func (authRest *AuthRest) Handle(context echo.Context) error {
	// find correct strategy
	host := context.Request().Host
	strategy, ok := authRest.UserAuth.Stragegies[host]
	if !ok {
		return errors.New("Can't find host: " + host)
	}

	// find correct function
	path := strings.Replace(context.Request().URL.String(), "/", "", 1)
	strategyFunc, ok := strategy.Functions[path]
	if !ok {
		return errors.New("Can't find route: " + path)
	}

	// get body
	requestMap := new(map[string]interface{})
	if err := context.Bind(requestMap); err != nil {
		return err
	}

	// call resolve of function
	returnInterf, err := strategyFunc.Resolve(context, (*requestMap))
	if err != nil {
		return err
	}

	return context.JSON(http.StatusOK, returnInterf)
}
