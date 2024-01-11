package handler

import (
	"fmt"
	"net/http"
	"project_layout/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

func (h *AuthHandler) SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		input, err := utils.GetInputAsMap(c)
		if err != nil {
			h.logger.Errorf("get input as map error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		validationResult, err := h.validator.Validate("auth_signup.json", input)
		if err != nil {
			h.logger.Errorf("get validator error: %v", err)

			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if validationResult != nil {
			h.logger.Errorf("get validator result error: %v", validationResult.Errors())

			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("%v", validationResult.Errors()),
			})
			return
		}

		result := graphql.Do(graphql.Params{
			Schema:  h.graphql,
			Context: c,
			RequestString: `
				mutation {
					signup {
						access_token
					}
				}
			`,
		})

		if result.HasErrors() {
			h.logger.Infof("error: %v", result.Errors[0])
		}

		c.JSON(http.StatusOK, gin.H{
			"data": result.Data.(map[string]interface{})["signup"],
		})
	}
}
