package metrics

import (
	"context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"homework/internal/models/config"
	"homework/internal/telemetry"
	"net/http"
)

func Listen(ctx context.Context, cfg *config.MetricsConfig, zap *zap.SugaredLogger) {
	telemetry.MustSetup(ctx, cfg.ServiceName)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	zap.Infof("Listening on %s", cfg.Addr)
	if err := http.ListenAndServe(cfg.Addr, mux); err != nil {
		zap.Errorln(err)
	}
}