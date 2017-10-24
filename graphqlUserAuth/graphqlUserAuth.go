package graphqlUserAuth

import (
	"log"
	"reflect"

	"github.com/fino-digital/dynamicUserAuth"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/labstack/echo"
)

type AuthSchema struct {
	UserAuth dynamicUserAuth.DynamicUserAuth
}

func (authSchema *AuthSchema) AuthSchema(c echo.Context) error {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name:   "AuthMutation",
			Fields: generateFields(authSchema.UserAuth.Stragegies, c),
		}),
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name:   "Query",
			Fields: graphql.Fields{"Test": &graphql.Field{Type: graphql.String}},
		}),
	})
	if err != nil {
		log.Println(err)
	}

	handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	}).ServeHTTP(c.Response(), c.Request())
	return err
}

func generateFields(strategies dynamicUserAuth.Stragegies, c echo.Context) graphql.Fields {
	// generate a type
	var generateType = func(name string, fields map[string]dynamicUserAuth.StrategyField) *graphql.Object {
		typeFields := graphql.Fields{}
		for key, value := range fields {
			typeFields[key] = &graphql.Field{Description: value.Description, Type: ScalarMap[value.Kind()], Name: name}
		}
		return graphql.NewObject(graphql.ObjectConfig{Name: name, Fields: typeFields})
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
	for key, value := range strategies[host].Functions {
		fields[key] = &graphql.Field{
			Type:        generateType("Output_"+key, value.Output),
			Args:        generateArgs(value.Input),
			Description: value.Description,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return value.Resolve(c, p.Args)
			},
		}
	}
	return fields
}

var ScalarMap = map[reflect.Kind]*graphql.Scalar{
	reflect.Bool:    graphql.Boolean,
	reflect.Int:     graphql.Int,
	reflect.Uint:    graphql.Int,
	reflect.String:  graphql.String,
	reflect.Float32: graphql.Float,
	reflect.Float64: graphql.Float,
}
