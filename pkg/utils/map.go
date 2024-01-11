package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func GetInputAsMap(
	c *gin.Context,
) (map[string]interface{}, error) {
	contentType := c.ContentType()
	if contentType != "application/json" {
		return nil, errors.New("Content-Type must be application/json")
	}

	// Getting the body as a map
	input := make(map[string]interface{})
	err := c.ShouldBindJSON(&input)
	if err != nil {
		return nil, err
	}

	return input, nil
}
