## Файлы для GC-MetricsServer - утилита анализа GC и памяти (runtime, профилирование)  

### 📋 Описание проекта и возможности  

**GC-MetricsServer** - это программа, которая через HTTP-endpoint предоставляет в формате Prometheus метрики памяти  
и сборщика мусора (GC), а также позволяет динамически настраивать процент GC и получать профили через pprof.  

### 📡 Доступные эндпойнты  

  - **/metrics** - метрики в формате Prometheus.
  - **/gc_percent** - получение и изменение текущего значения GOGC.
  - **GET** - возвращает текущий процент.
  - **POST** с параметром percent - устанавливает новый процент.
  - **/debug/pprof/** - стандартные профили pprof (heap, goroutine, CPU и т.д.).  

### 🗂️ Структура проекта  

``` bash
.
├── cmd
│   └── main.go                 # точка входа, запуск сервера
├── internal
│   ├── api.go                  # обработчики /api/stats и /gc_percent
│   ├── collector.go            # Prometheus коллектор и вспомогательные функции
│   └── server.go               # инициализация Gin, graceful shutdown
├── web
│   ├── index.html
│   ├── script.js
│   └── style.css
├── .env
├── compose.yml
├── Dockerfile
├── go.mod
├── go.sum
└── readme.md
```

### 🚀 Быстрый старт  

**Требования:**  
- Docker и Docker Compose  
- Свободный порт: 8080

**Запуск:**  

    docker compose up --build

**После успешного запуска:**  

    http://localhost:8080

Остановка:  

    docker compose down

### ⚙️ Конфигурация  

Адрес сервера задаётся через файл .env в корне проекта.  

    HTTP_HOST=localhost          # хост сервиса
    HTTP_PORT=8080               # порт хоста, на котором работает сервис

### 📦 Метрики  


