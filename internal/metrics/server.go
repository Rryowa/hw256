package metrics

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func Listen(addr string) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	//server := &http.Server{Addr: addr, Handler: mux}

	//go func() {
	//	<-ctx.Done()
	//
	//	log.Warnf("Shutting down server with duration %0.3fs", shutdownDuration.Seconds())
	//	<-time.After(shutdownDuration)
	//
	//	if err := server.Shutdown(context.Background()); err != nil {
	//		log.Errorf("HTTP handler Shutdown: %s", err)
	//	}
	//}()
	//
	//if err := server.ListenAndServe(); err != nil {
	//	log.Errorf("HTTP server ListenAndServe: %s", err)
	//	return err
	//}
	log.Infof("Listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Errorln(err)
	}
	return
}