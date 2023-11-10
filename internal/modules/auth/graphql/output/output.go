package output

import (
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
)

type AuthOutput struct {
	AccessToken string
}

func AuthOutputType(
	types map[string]*graphql.Object,
	logger *logrus.Logger,
) *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "auth",
			Fields: graphql.FieldsThunk(func() graphql.Fields {
				return graphql.Fields{
					"access_token": &graphql.Field{
						Type: graphql.String,
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							return p.Source.(AuthOutput).AccessToken, nil
						},
					},
				}
			}),
		},
	)
}
