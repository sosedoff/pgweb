package api

import (
	"github.com/gin-gonic/gin"

	"github.com/sosedoff/pgweb/pkg/command"
	"github.com/sosedoff/pgweb/pkg/metrics"
)

func SetupMiddlewares(group *gin.RouterGroup) {
	if command.Opts.Cors {
		group.Use(corsMiddleware())
	}

	group.Use(dbCheckMiddleware())
}

func SetupRoutes(router *gin.Engine) {
	root := router.Group(command.Opts.Prefix)

	root.GET("/", gin.WrapH(GetHome(command.Opts.Prefix)))
	root.GET("/static/*path", gin.WrapH(GetAssets(command.Opts.Prefix)))
	root.GET("/connect/:resource", ConnectWithBackend)

	api := root.Group("/api")
	SetupMiddlewares(api)

	if command.Opts.Sessions {
		api.GET("/sessions", GetSessions)
	}

	api.GET("/info", GetInfo)
	api.POST("/connect", Connect)
	api.POST("/disconnect", Disconnect)
	api.POST("/switchdb", SwitchDb)
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
	api.GET("/tables_stats", GetTablesStats)
	api.GET("/functions/:id", GetFunction)
	api.GET("/query", RunQuery)
	api.POST("/query", RunQuery)
	api.GET("/explain", ExplainQuery)
	api.POST("/explain", ExplainQuery)
	api.GET("/analyze", AnalyzeQuery)
	api.POST("/analyze", AnalyzeQuery)
	api.GET("/history", GetHistory)
	api.GET("/bookmarks", GetBookmarks)
	api.GET("/export", DataExport)
	api.GET("/local_queries", requireLocalQueries(), GetLocalQueries)
	api.GET("/local_queries/:id", requireLocalQueries(), RunLocalQuery)
	api.POST("/local_queries/:id", requireLocalQueries(), RunLocalQuery)
}

func SetupMetrics(engine *gin.Engine) {
	if command.Opts.MetricsEnabled && command.Opts.MetricsAddr == "" {
		engine.GET("/metrics", gin.WrapH(metrics.Handler()))
	}
}
