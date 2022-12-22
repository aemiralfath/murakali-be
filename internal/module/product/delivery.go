package product

import "github.com/gin-gonic/gin"

type Handlers interface {
	GetCategories(c *gin.Context)
	GetCategoriesByNameLevelOne(c *gin.Context)
	GetCategoriesByNameLevelTwo(c *gin.Context)
	GetCategoriesByNameLevelThree(c *gin.Context)
}