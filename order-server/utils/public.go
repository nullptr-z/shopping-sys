package utils

import (
	"fmt"
	"math/rand"
	"time"

	"go.uber.org/zap"
)

func GenerateOrderSn(userid int32) string {
	// yyyymmddhhmmss+userid+random(2)
	now := time.Now()
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(89999) + 10000
	Sn := fmt.Sprintf("%d%02d%02d%02d%02d%02d%04d%04d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), userid, randNum)
	zap.S().Info("[GenerateOrderSn] orderSn:", Sn)
	return Sn
}
