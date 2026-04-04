package internal

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Run() error {
	_ = godotenv.Load()

	host := getEnv("HTTP_HOST", "localhost")
	port := getEnv("HTTP_PORT", "8080")
	addr := host + ":" + port

	metricsTimeout := getEnv("METRICS_TIMEOUT", "3s")
	if _, err := time.ParseDuration(metricsTimeout); err != nil {
		log.Printf("Некорректный METRICS_TIMEOUT='%s', используем 3s", metricsTimeout)
		metricsTimeout = "3s"
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())

	// Prometheus
	reg := prometheus.NewRegistry()
	reg.MustRegister(NewMemStatsCollector())
	r.GET("/metrics", gin.WrapH(promhttp.HandlerFor(reg, promhttp.HandlerOpts{})))

	// API
	r.GET("/api/stats", HandleAPIMetrics)
	r.GET("/gc_percent", HandleGCPercent)
	r.POST("/gc_percent", HandleGCPercent)

	r.GET("/api/config", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"metrics_timeout": metricsTimeout})
	})

	// Статика и главная страница
	r.Static("/static", "./web/static")
	r.GET("/", func(c *gin.Context) {
		c.File("./web/index.html")
	})

	// pprof
	pprofMux := http.NewServeMux()
	pprofMux.Handle("/debug/pprof/", http.DefaultServeMux)
	r.Any("/debug/pprof/*any", gin.WrapH(pprofMux))

	srv := &http.Server{Addr: addr, Handler: r}

	go func() {
		log.Printf("Сервер запущен на http://%s", addr)
		log.Printf("Веб-интерфейс: http://%s/", addr)
		log.Printf("Prometheus: http://%s/metrics", addr)
		log.Printf("GOGC: GET/POST http://%s/gc_percent", addr)
		log.Printf("pprof: http://%s/debug/pprof/", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Останавливаем сервер...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(ctx)
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
