package utils

import (
	"go.uber.org/zap"
)

// TODO: change name from ReplaceLogger to Logger after removing internal logger.
var ReplaceLogger, _ = zap.NewDevelopment()

// func ReplaceLogger() (*zap.Logger, error) {
// 	z, err := zap.NewDevelopment()
// 	if err != nil {
// 		return nil, fmt.Errorf("error occured setting up zap: %v", err)
// 	}

// 	return z, nil
// }
