package mutation

import (
	"project_layout/internal/modules/auth/graphql/output"
	userRepo "project_layout/internal/modules/auth/repository"

	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
)

func SignUpMutationType(
	types map[string]*graphql.Object,
	userRepo userRepo.UserRepositoryInterface,
	logger *logrus.Logger,
) *graphql.Field {
	return &graphql.Field{
		Type: types["auth"],
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			err := userRepo.CreateUser()
			if err != nil {
				return nil, err
			}
			return output.AuthOutput{
				AccessToken: "token",
			}, nil
		},
	}
}
