package internal

import (
	"runtime"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Структура для JSON-ответа /api/stats
type MetricsJSON struct {
	AllocBytes      string  `json:"alloc_bytes"`       // используемая память (человеко‑читаемый вид)
	TotalAllocBytes string  `json:"total_alloc_bytes"` // всего выделено за всё время
	NumGC           uint32  `json:"num_gc"`            // количество сборок GC
	LastGCTime      string  `json:"last_gc_time"`      // время последнего GC в формате "2006-01-02 15:04:05.000"
	GCCPUFraction   float64 `json:"gc_cpu_fraction"`   // доля CPU на GC
	GCPercents      int     `json:"gc_percent"`        // текущее значение GOGC
}

// Кастомный коллектор для Prometheus
type MemStatsCollector struct {
	allocBytes      *prometheus.Desc
	totalAllocBytes *prometheus.Desc
	numGC           *prometheus.Desc
	lastGCTime      *prometheus.Desc
	gcCPUFraction   *prometheus.Desc
	gcPercent       *prometheus.Desc
}

// Новый коллектор
func NewMemStatsCollector() *MemStatsCollector {
	return &MemStatsCollector{
		allocBytes: prometheus.NewDesc(
			"go_memstats_alloc_bytes",
			"Количество байт, выделенных и используемых в данный момент (HeapAlloc)",
			nil, nil,
		),
		totalAllocBytes: prometheus.NewDesc(
			"go_memstats_total_alloc_bytes",
			"Общее количество выделенных байт за всё время (TotalAlloc)",
			nil, nil,
		),
		numGC: prometheus.NewDesc(
			"go_gc_num_gc",
			"Количество завершённых циклов GC",
			nil, nil,
		),
		lastGCTime: prometheus.NewDesc(
			"go_gc_last_gc_time_seconds",
			"Время последнего GC в секундах с эпохи Unix",
			nil, nil,
		),
		gcCPUFraction: prometheus.NewDesc(
			"go_gc_cpu_fraction",
			"Доля CPU, потраченная на GC с момента запуска программы",
			nil, nil,
		),
		gcPercent: prometheus.NewDesc(
			"go_gc_percent",
			"Текущее значение GOGC (процент)",
			nil, nil,
		),
	}
}

// Describe реализует интерфейс Collector
func (c *MemStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.allocBytes
	ch <- c.totalAllocBytes
	ch <- c.numGC
	ch <- c.lastGCTime
	ch <- c.gcCPUFraction
	ch <- c.gcPercent
}

// Collect собирает метрики из runtime
func (c *MemStatsCollector) Collect(ch chan<- prometheus.Metric) {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	// Получаем текущий GC процент без изменения
	gcPercent := debug.SetGCPercent(-1)
	debug.SetGCPercent(gcPercent)

	ch <- prometheus.MustNewConstMetric(c.allocBytes, prometheus.GaugeValue, float64(ms.Alloc))
	ch <- prometheus.MustNewConstMetric(c.totalAllocBytes, prometheus.CounterValue, float64(ms.TotalAlloc))
	ch <- prometheus.MustNewConstMetric(c.numGC, prometheus.CounterValue, float64(ms.NumGC))
	ch <- prometheus.MustNewConstMetric(c.lastGCTime, prometheus.GaugeValue, float64(ms.LastGC)/1e9)
	ch <- prometheus.MustNewConstMetric(c.gcCPUFraction, prometheus.GaugeValue, ms.GCCPUFraction)
	ch <- prometheus.MustNewConstMetric(c.gcPercent, prometheus.GaugeValue, float64(gcPercent))
}

// formatBytes переводит байты в KB/MB/GB с двумя знаками
func formatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return strconv.FormatUint(b, 10) + " B"
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return strconv.FormatFloat(float64(b)/float64(div), 'f', 2, 64) + " " + []string{"KB", "MB", "GB", "TB"}[exp] + "B"
}

// GetCurrentMetrics возвращает актуальные метрики в удобном для JSON виде
func GetCurrentMetrics() MetricsJSON {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	gcPercent := debug.SetGCPercent(-1)
	debug.SetGCPercent(gcPercent)

	lastGCTimeStr := "ещё не было"
	if ms.LastGC != 0 {
		t := time.Unix(0, int64(ms.LastGC))
		lastGCTimeStr = t.Format("2006-01-02 15:04:05.000")
	}

	return MetricsJSON{
		AllocBytes:      formatBytes(ms.Alloc),
		TotalAllocBytes: formatBytes(ms.TotalAlloc),
		NumGC:           ms.NumGC,
		LastGCTime:      lastGCTimeStr,
		GCCPUFraction:   ms.GCCPUFraction,
		GCPercents:      gcPercent,
	}
}
