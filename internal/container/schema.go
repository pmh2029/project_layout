package container

import (
	"project_layout/internal/modules/auth/graphql/mutation"
	"project_layout/internal/modules/auth/graphql/output"
	"project_layout/internal/modules/auth/graphql/query"

	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func NewGraphQLSchema(
	repositories *RepositoryContainer,
	db *gorm.DB,
	logger *logrus.Logger,
) (graphql.Schema, error) {
	outputTypes := make(map[string]*graphql.Object)
	for _, graphqlType := range []*graphql.Object{
		output.AuthOutputType(outputTypes, logger),
	} {
		outputTypes[graphqlType.Name()] = graphqlType
	}

	return graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"all": query.GetAllUsersQueryType(
					outputTypes,
					logger,
				),
			},
		}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name: "Mutation",
			Fields: graphql.Fields{
				"signup": mutation.SignUpMutationType(
					outputTypes,
					logger,
				),
			},
		}),
	})
}
