package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/insights/backend/internal/handler"
	"github.com/insights/backend/internal/middleware"
)

func main() {
	h := server.Default(
		server.WithHostPorts("0.0.0.0:8888"),
	)

	h.Use(middleware.CORS())

	h.POST("/api/promql/parse", handler.ParsePromQL)
	h.POST("/api/query/range", handler.QueryRange)

	h.Spin()
}
