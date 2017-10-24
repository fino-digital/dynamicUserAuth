package dynamicUserAuth

import (
	"errors"
	"reflect"

	"github.com/labstack/echo"
)

// Stragegies is the map of pointer for strategies.
// key: host, value: strategy
type Stragegies map[string]*Strategy

// DynamicUserAuth holds all stragegies for different products.
// Expand this for new products.
type DynamicUserAuth struct {
	// Stragegies holds host to strategy
	Stragegies Stragegies
}

// StrategyField describes a field for input or output of a strategie
type StrategyField struct {
	reflect.Type
	Description string
	Required    bool
}

// StrategyFunction can be for example "newUser"
type StrategyFunction struct {
	Description string
	Input       map[string]StrategyField
	Output      map[string]StrategyField
	Resolve     func(map[string]StrategyField) (map[string]StrategyField, error)
}

// Strategy represent a strategy for one product.
// Implement a new strategy for a new product
type Strategy struct {
	Functions     map[string]StrategyFunction
	AuthorizeUser echo.HandlerFunc
}

// AuthMiddleware is the middleare for all auth-stuff.
type AuthMiddleware struct {
	dynamicUserAuth *DynamicUserAuth
}

// NewAuthMiddleware creates a new authMiddleware.
// this function is here to force to get all requirements
func NewAuthMiddleware(dynamicUserAuth *DynamicUserAuth) *AuthMiddleware {
	return &AuthMiddleware{dynamicUserAuth: dynamicUserAuth}
}

// Handle handles the auth-process.
// Use this for all save-endpoints.
func (authMiddleware *AuthMiddleware) Handle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		host := context.Request().Host
		// Check first if strategy for this host exist.
		// If-else-construct is confused (`return next(context)` should be at the end).
		// - If you find a better way, plz go for it!
		if strategy, ok := authMiddleware.dynamicUserAuth.Stragegies[host]; ok {
			if err := strategy.AuthorizeUser(context); err != nil {
				return err
			}
			return next(context)
		}
		return errors.New("can't find strategy")
	}
}
