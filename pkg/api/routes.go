package api

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/", GetHome)
	router.GET("/static/*path", GetAsset)

	api := router.Group("/api")
	{
		api.Use(dbCheckMiddleware())

		api.GET("/info", GetInfo)
		api.POST("/connect", Connect)
		api.GET("/databases", GetDatabases)
		api.GET("/connection", GetConnectionInfo)
		api.GET("/activity", GetActivity)
		api.GET("/schemas", GetSchemas)
		api.GET("/tables", GetTables)
		api.GET("/tables/:table", GetTable)
		api.GET("/tables/:table/rows", GetTableRows)
		api.GET("/tables/:table/info", GetTableInfo)
		api.GET("/tables/:table/indexes", GetTableIndexes)
		api.GET("/query", RunQuery)
		api.POST("/query", RunQuery)
		api.GET("/explain", ExplainQuery)
		api.POST("/explain", ExplainQuery)
		api.GET("/history", GetHistory)
		api.GET("/bookmarks", GetBookmarks)
	}
}
