package main

import (
	"github.com/h3th-IV/mysticMerch/internal/api"
	"github.com/h3th-IV/mysticMerch/internal/server"
	"github.com/h3th-IV/mysticMerch/internal/utils"
	"go.uber.org/zap"
)

func main() {
	api.Test()
	utils.ReplaceLogger.Info("Starting Server at", zap.String("port", "8000"))
	server.StartServer()
}
