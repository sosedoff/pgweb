package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sosedoff/pgweb/pkg/command"
)

func SetupMiddlewares(group *gin.RouterGroup) {
	if command.Opts.Debug {
		group.Use(requestInspectMiddleware())
	}

	group.Use(dbCheckMiddleware())
}

func SetupRoutes(router *gin.Engine) {
	router.GET("/", GetHome)
	router.GET("/static/*path", GetAsset)

	api := router.Group("/api")
	{
		SetupMiddlewares(api)

		if command.Opts.Sessions {
			api.GET("/sessions", GetSessions)
		}

		api.GET("/info", GetInfo)
		api.POST("/connect", Connect)
		api.GET("/databases", GetDatabases)
		api.GET("/connection", GetConnectionInfo)
		api.GET("/activity", GetActivity)
		api.GET("/schemas", GetSchemas)
		api.GET("/objects", GetObjects)
		api.GET("/tables/:table", GetTable)
		api.GET("/tables/:table/rows", GetTableRows)
		api.GET("/tables/:table/info", GetTableInfo)
		api.GET("/tables/:table/indexes", GetTableIndexes)
		api.GET("/tables/:table/constraints", GetTableConstraints)
		api.GET("/query", RunQuery)
		api.POST("/query", RunQuery)
		api.GET("/explain", ExplainQuery)
		api.POST("/explain", ExplainQuery)
		api.GET("/history", GetHistory)
		api.GET("/bookmarks", GetBookmarks)
	}
}
