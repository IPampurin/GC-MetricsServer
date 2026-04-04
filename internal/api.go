package internal

import (
	"net/http"
	"runtime/debug"
	"strconv"

	"github.com/gin-gonic/gin"
)

// HandleAPIMetrics возвращает JSON с текущими метриками
func HandleAPIMetrics(c *gin.Context) {
	c.JSON(http.StatusOK, GetCurrentMetrics())
}

// HandleGCPercent обрабатывает GET и POST запросы для управления GOGC
func HandleGCPercent(c *gin.Context) {
	switch c.Request.Method {
	case http.MethodGet:
		percent := debug.SetGCPercent(-1)
		debug.SetGCPercent(percent)
		c.JSON(http.StatusOK, gin.H{"gc_percent": percent})

	case http.MethodPost:
		percentStr := c.Query("percent")
		if percentStr == "" {
			c.String(http.StatusBadRequest, "Не указан параметр 'percent'")
			return
		}
		percent, err := strconv.Atoi(percentStr)
		if err != nil {
			c.String(http.StatusBadRequest, "Некорректное значение percent")
			return
		}
		old := debug.SetGCPercent(percent)
		c.JSON(http.StatusOK, gin.H{
			"old_percent": old,
			"new_percent": percent,
		})

	default:
		c.String(http.StatusMethodNotAllowed, "Метод не разрешён")
	}
}
