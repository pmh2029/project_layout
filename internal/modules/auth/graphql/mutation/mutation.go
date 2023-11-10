package mutation

import (
	"project_layout/internal/modules/auth/graphql/output"

	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
)

func SignUpMutationType(
	types map[string]*graphql.Object,
	logger *logrus.Logger,
) *graphql.Field {
	return &graphql.Field{
		Type: types["auth"],
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return output.AuthOutput{
				AccessToken: "token",
			}, nil
		},
	}
}
