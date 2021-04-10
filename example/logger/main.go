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

	err := fmt.Errorf("just test")

	log := logutil.WithContext(context.TODO())
	log.Info("test", zap.Error(err), zap.Int("data", 1))

	log2 := logutil.WithContext(context.WithValue(context.Background(), constant.RequestIDContextKey, "2"))
	log2.Info("test")

	log3 := logutil.WithContext(context.WithValue(context.Background(), constant.RequestIDContextKey, "3"))
	log3.Info("test")

	log2.Info("test again")
	log.Info("test", zap.Error(err), zap.Int("data", 2))
}
