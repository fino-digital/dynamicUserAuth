package graphqlUserAuth

import (
	"reflect"

	"github.com/fino-digital/dynamicUserAuth"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/labstack/echo"
)

type AuthSchema struct {
	UserAuth dynamicUserAuth.DynamicUserAuth
}

func (authSchema *AuthSchema) authSchema(c echo.Context) error {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name:   "AuthMutation",
			Fields: authSchema.generateFields(c),
		}),
	})

	handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	}).ServeHTTP(c.Response(), c.Request())
	return err
}

func (authSchema *AuthSchema) generateFields(c echo.Context) *graphql.Fields {
	// generate a type
	var generateType = func(name string, fields map[string]dynamicUserAuth.StrategyField) *graphql.Object {
		typeFields := graphql.Fields{}
		for key, value := range fields {
			typeFields[key] = &graphql.Field{Description: value.Description, Type: ScalarMap[value.Kind()]}
		}
		return graphql.NewObject(graphql.ObjectConfig{Name: name, Fields: fields})
	}

	// generate arguments
	var generateArgs = func(fields map[string]dynamicUserAuth.StrategyField) graphql.FieldConfigArgument {
		argument := graphql.FieldConfigArgument{}
		for key, value := range fields {
			if value.Required {
				argument[key] = &graphql.ArgumentConfig{Type: graphql.NewNonNull(ScalarMap[value.Kind()]), Description: value.Description}
			} else {
				argument[key] = &graphql.ArgumentConfig{Type: ScalarMap[value.Kind()]}
			}
		}
		return argument
	}

	// iterate functions of host
	host := c.Request().Host
	fields := graphql.Fields{}
	for key, value := range authSchema.UserAuth.Stragegies[host].Functions {
		fields[key] = &graphql.Field{
			Type:        generateType("Output", value.Output),
			Args:        generateArgs(value.Input),
			Description: value.Description,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				inputFields := map[string]dynamicUserAuth.StrategyField{}
				for key, value := range value.Input {
					field := p.Args[key].(dynamicUserAuth.StrategyField)
					field.Description = value.Description
					field.Required = value.Required
					inputFields[key] = field
				}
				return value.Resolve(inputFields)
			},
		}
	}
	return &fields
}

var ScalarMap = map[reflect.Kind]*graphql.Scalar{
	reflect.Bool:    graphql.Boolean,
	reflect.Int:     graphql.Int,
	reflect.String:  graphql.String,
	reflect.Float32: graphql.Float,
	reflect.Float64: graphql.Float,
}
