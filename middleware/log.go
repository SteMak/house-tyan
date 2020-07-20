package middleware

import (
	"time"

	"github.com/SteMak/house-tyan/app"
	"github.com/sirupsen/logrus"
)

func Log(logger *logrus.Logger) app.HandlerFunc {
	return func(ctx *app.Context) {
		start := time.Now()
		ctx.Next()
		logger.WithFields(
			logrus.Fields{
				"duration": time.Since(start),
			},
		).Trace()
	}
}
