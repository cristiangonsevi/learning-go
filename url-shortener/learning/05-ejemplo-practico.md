# Ejemplo Práctico - Todo junto

## El proyecto: Servidor con worker

Vamos a crear un servidor HTTP que:
1. Responde a requests (goroutine del servidor)
2. Tiene un worker que limpia datos viejos (goroutine del worker)
3. Se apaga limpiamente (context + channels)

## Código completo

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "sync"
    "syscall"
    "time"
)

// ============ MODELO ============

type URL struct {
    ID          int
    Code        string
    OriginalURL string
    CreatedAt   time.Time
    ExpiresAt   *time.Time
    IsActive    bool
}

// ============ SIMULACIÓN DE BASE DE DATOS ============

type Database struct {
    urls []URL
    mu   sync.RWMutex
}

func NewDatabase() *Database {
    return &Database{
        urls: []URL{
            {ID: 1, Code: "abc", OriginalURL: "https://google.com", CreatedAt: time.Now(), IsActive: true},
            {ID: 2, Code: "def", OriginalURL: "https://github.com", CreatedAt: time.Now().Add(-24 * time.Hour), IsActive: true},
            {ID: 3, Code: "ghi", OriginalURL: "https://expired.com", CreatedAt: time.Now().Add(-8 * 24 * time.Hour), ExpiresAt: timePtr(time.Now().Add(-1 * 24 * time.Hour)), IsActive: true},
        },
    }
}

func timePtr(t time.Time) *time.Time {
    return &t
}

func (db *Database) GetURL(code string) *URL {
    db.mu.RLock()
    defer db.mu.RUnlock()
    
    for _, url := range db.urls {
        if url.Code == code && url.IsActive {
            return &url
        }
    }
    return nil
}

func (db *Database) DeactivateExpiredURLs() int {
    db.mu.Lock()
    defer db.mu.Unlock()
    
    count := 0
    now := time.Now()
    
    for i := range db.urls {
        if db.urls[i].ExpiresAt != nil && db.urls[i].ExpiresAt.Before(now) && db.urls[i].IsActive {
            db.urls[i].IsActive = false
            count++
        }
    }
    return count
}

func (db *Database) GetStats() (total, active, expired int) {
    db.mu.RLock()
    defer db.mu.RUnlock()
    
    for _, url := range db.urls {
        total++
        if url.IsActive {
            active++
        } else {
            expired++
        }
    }
    return
}

// ============ WORKER ============

type CleanupWorker struct {
    db     *Database
    cancel context.CancelFunc
}

func NewCleanupWorker(db *Database) *CleanupWorker {
    return &CleanupWorker{db: db}
}

func (w *CleanupWorker) Start() {
    ctx, cancel := context.WithCancel(context.Background())
    w.cancel = cancel

    go func() {
        // Ejecutar inmediatamente al arrancar
        w.cleanup()
        
        // Luego cada 30 segundos (para el ejemplo, en producción sería cada hora)
        ticker := time.NewTicker(30 * time.Second)
        defer ticker.Stop()

        for {
            select {
            case <-ticker.C:
                w.cleanup()
            case <-ctx.Done():
                fmt.Println("🛑 Worker apagado")
                return
            }
        }
    }()
    
    fmt.Println("🚀 Worker iniciado")
}

func (w *CleanupWorker) Stop() {
    w.cancel()
}

func (w *CleanupWorker) cleanup() {
    count := w.db.DeactivateExpiredURLs()
    if count > 0 {
        fmt.Printf("🧹 Cleanup: %d URLs desactivadas\n", count)
    } else {
        fmt.Println("✅ Cleanup: no hay URLs expiradas")
    }
}

// ============ HANDLERS ============

type Handler struct {
    db *Database
}

func NewHandler(db *Database) *Handler {
    return &Handler{db: db}
}

func (h *Handler) HandleRedirect(w http.ResponseWriter, r *http.Request) {
    code := r.URL.Path[1:] // Quitar el "/"
    
    url := h.db.GetURL(code)
    if url == nil {
        http.Error(w, "URL no encontrada", http.StatusNotFound)
        return
    }
    
    // Redirigir
    http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}

func (h *Handler) HandleStats(w http.ResponseWriter, r *http.Request) {
    total, active, expired := h.db.GetStats()
    
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, `{"total": %d, "active": %d, "expired": %d}`, total, active, expired)
}

func (h *Handler) HandleHealth(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, `{"status": "ok"}`)
}

// ============ MAIN ============

func main() {
    fmt.Println("🔧 Iniciando servidor...")
    
    // 1. Database
    db := NewDatabase()
    
    // 2. Worker
    cleanupWorker := NewCleanupWorker(db)
    cleanupWorker.Start()
    
    // 3. Handlers
    handler := NewHandler(db)
    
    // 4. Router
    mux := http.NewServeMux()
    mux.HandleFunc("/health", handler.HandleHealth)
    mux.HandleFunc("/stats", handler.HandleStats)
    mux.HandleFunc("/", handler.HandleRedirect)
    
    // 5. Servidor
    server := &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }
    
    // Arrancar servidor en goroutine
    go func() {
        fmt.Println("🌐 Servidor en http://localhost:8080")
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()
    
    // 6. Esperar señal de apagado
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    
    fmt.Println("⏳ Presiona Ctrl+C para apagar")
    <-quit  // Bloquear hasta recibir señal
    
    // 7. Apagar limpiamente
    fmt.Println("\n🛑 Apagando servidor...")
    
    // Apagar worker
    cleanupWorker.Stop()
    
    // Apagar servidor HTTP
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := server.Shutdown(ctx); err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("✅ Servidor apagado correctamente")
}
```

## Cómo probarlo

```bash
# 1. Ejecutar
go run main.go

# 2. En otra terminal, hacer requests
curl http://localhost:8080/abc      # → Redirect a google.com
curl http://localhost:8080/def      # → Redirect a github.com
curl http://localhost:8080/ghi      # → 404 (expirada)
curl http://localhost:8080/stats    # → {"total":3,"active":2,"expired":1}
curl http://localhost:8080/health   # → {"status":"ok"}

# 3. Presionar Ctrl+C para apagar
```

## Qué se ve en consola

```
🔧 Iniciando servidor...
🚀 Worker iniciado
✅ Cleanup: no hay URLs expiradas
🌐 Servidor en http://localhost:8080
⏳ Presiona Ctrl+C para apagar

🧹 Cleanup: 1 URLs desactivadas    ← Cada 30 segundos

^C
🛑 Apagando servidor...
🛑 Worker apagado
✅ Servidor apargado correctamente
```

## Conceptos usados

| Concepto | Dónde se usa |
|----------|--------------|
| Goroutine | Servidor HTTP (`go func()`) |
| Goroutine | Worker (`go func()`) |
| Channel | Señal de apagado (`quit := make(chan os.Signal)`) |
| Channel | Apagar worker (`done := make(chan bool)`) |
| Context | Shutdown del servidor (`context.WithTimeout`) |
| Context | Cancelar worker (`context.WithCancel`) |
| WaitGroup | No lo usamos aquí pero es para esperar múltiples goroutines |
| Mutex | Proteger la base de datos (`sync.RWMutex`) |

## Diagrama

```
main()
  │
  ├──► NewDatabase()
  │
  ├──► worker.Start() ──► goroutine {
  │         │                 ticker: cada 30s
  │         │                 select {
  │         │                   case <-ticker.C: cleanup()
  │         │                   case <-ctx.Done(): return
  │         │                 }
  │         └──────────────► }
  │
  ├──► server.ListenAndServe() ──► goroutine {
  │         │                        Escucha en :8080
  │         └──────────────────────► }
  │
  └──► <-quit (espera Ctrl+C)
          │
          ├── worker.Stop()
          └── server.Shutdown()
```
