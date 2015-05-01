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

		api.POST("/connect", GetConnect)
		api.GET("/databases", GetGetDatabases)
		api.GET("/connection", GetConnectionInfo)
		api.GET("/activity", GetActivity)
		api.GET("/schemas", GetGetSchemas)
		api.GET("/tables", GetGetTables)
		api.GET("/tables/:table", GetGetTable)
		api.GET("/tables/:table/rows", GetGetTableRows)
		api.GET("/tables/:table/info", GetGetTableInfo)
		api.GET("/tables/:table/indexes", GetTableIndexes)
		api.GET("/query", GetRunQuery)
		api.POST("/query", GetRunQuery)
		api.GET("/explain", GetExplainQuery)
		api.POST("/explain", GetExplainQuery)
		api.GET("/history", GetHistory)
		api.GET("/bookmarks", GetBookmarks)
	}
}
