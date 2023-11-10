package mount

import (
	"project_layout/internal/pkg/container"
	"project_layout/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// MountAll mount all dependencies
func MountAll(
	gin *gin.Engine,
	db *gorm.DB,
	logger *logrus.Logger,
) error {
	// initialize validation
	varlidator, err := validator.NewJsonSchemaValidator()
	if err != nil {
		logger.Infof("Failed to create a JSON Schema validator: %v", err.Error())
		return err
	}

	// initialize repositories container
	repositories := container.NewRepositoryContainer(db, logger)

	// initialize graphql schema
	graphqlSchema, err := container.NewGraphQLSchema(&repositories, db, logger)
	if err != nil {
		logger.Infof("Failed to create a GraphQL schema: %v", err.Error())
		return err
	}

	// initialize handler container
	handlerContainer := container.NewHandlerContainer(varlidator, graphqlSchema, logger)

	// set up hanlder container
	container.SetupHandler(gin, &handlerContainer)
	return nil
}
