package main

import (
	"context"
	"fmt"

	"github.com/tongium/common-go/pkg/constant"
	"github.com/tongium/common-go/pkg/logutil"
	"go.uber.org/zap"
)

func main() {
	logutil.New("info", "console")
	defer logutil.Logger().Sync()

	err := fmt.Errorf("error for test")

	mainLog := logutil.WithContext(context.TODO())
	mainLog.Info("papaya", zap.Error(err), zap.Int("number", 99))

	appleLog := logutil.WithContext(context.WithValue(context.Background(), constant.RequestIDContextKey, "apple"))
	appleLog.Info("apple")

	orangeLog := logutil.WithContext(context.WithValue(context.Background(), constant.RequestIDContextKey, "orange"))
	orangeLog.Info("orange")

	appleLog.Info("test again")
	mainLog.Info("test", zap.Error(err), zap.String("word", "amen"))
}
