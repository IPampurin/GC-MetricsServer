package main

import (
	"log"

	"github.com/IPampurin/GC-MetricsServer/internal"
)

func main() {

	if err := internal.Run(); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
