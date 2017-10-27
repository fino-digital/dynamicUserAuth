package restUserAuth

import (
	"net/http"

	"github.com/fino-digital/dynamicUserAuth"
	"github.com/labstack/echo"
)

// FunctionKeyWord is the keyword for the param
const FunctionKeyWord = "function"

// ParamFunction is the param-name for route
const ParamFunction = "/:" + FunctionKeyWord

// StatusNoHost can't find host
const StatusNoHost = 440

// StatusNoFunction can't find function
const StatusNoFunction = 441

// AuthRest is the instance for Auth-REST-impl
type AuthRest struct {
	UserAuth dynamicUserAuth.DynamicUserAuth
}

// Handle handles all functions dynamic by host
func (authRest *AuthRest) Handle(context echo.Context) error {
	function := context.Param(FunctionKeyWord)

	// find correct strategy
	host := context.Request().Host
	strategy, ok := authRest.UserAuth.Stragegies[host]
	if !ok {
		return context.JSON(StatusNoHost, "Can't find host: "+host)
	}

	// find correct function
	strategyFunc, ok := strategy.Functions[function]
	if !ok {
		return context.JSON(StatusNoFunction, "Can't find function: "+function)
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
