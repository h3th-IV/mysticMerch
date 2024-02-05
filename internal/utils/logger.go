package utils

import "go.uber.org/zap"

// TODO: change name from ReplaceLogger to Logger after removing internal logger.
var ReplaceLogger, _ = zap.NewDevelopment()
