package main

import (
	"context"
	"fmt"

	"github.com/tongium/common-go/pkg/constant"
	"github.com/tongium/common-go/pkg/logutil"
	"go.uber.org/zap"
)

func main() {
	defer logutil.Logger().Sync() // Flush before exit

	err := fmt.Errorf("error for test")

	mainLog := logutil.WithContext(context.TODO())
	mainLog.Info("papaya", zap.Error(err), zap.Int("number", 99))

	appleLog := logutil.WithContext(context.WithValue(context.Background(), constant.RequestIDContextKey, "apple"))
	appleLog.Debug("apple")

	orangeLog := logutil.WithContext(context.WithValue(context.Background(), constant.RequestIDContextKey, "orange"))
	orangeLog.Error("orange")

	appleLog.Info("test again")
	mainLog.Warn("test", zap.Error(err), zap.String("word", "amen"))
}
