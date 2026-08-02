# Workers - Tareas que corren en background

## ¿Qué es un worker?

Es una goroutine que hace una tarea **repetitiva** en background.

## Analogía

```
Worker = Empleado que hace lo mismo cada día

┌─────────────────────────────────────────────────┐
│                CLEANUP WORKER                    │
├─────────────────────────────────────────────────┤
│                                                  │
│  8:00 AM → Limpia URLs expiradas                 │
│  8:01 AM → Duerme hasta mañana                   │
│  ...                                             │
│  8:00 AM → Limpia URLs expiradas                 │
│  8:01 AM → Duerme hasta mañana                   │
│  ...                                             │
│                                                  │
└─────────────────────────────────────────────────┘
```

## Ejemplo 1: Worker simple

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // Iniciar worker
    go cleanupWorker()
    
    // El main sigue haciendo otras cosas
    fmt.Println("Servidor arrancado")
    time.Sleep(10 * time.Second) // Simular que el servidor corre
}

func cleanupWorker() {
    for {
        fmt.Println("Ejecutando cleanup...")
        // Aquí va tu lógica de limpieza
        time.Sleep(5 * time.Second) // Esperar 5 segundos
    }
}
```

**Salida:**
```
Servidor arrancado
Ejecutando cleanup...
Ejecutando cleanup...
Ejecutando cleanup...
```

## Ejemplo 2: Worker con ticker (como cronjob)

```go
func cleanupWorker() {
    // Ticker = setInterval en JavaScript
    ticker := time.NewTicker(24 * time.Hour) // Cada 24 horas
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:  // Cuando el ticker dispara
            fmt.Println("Ejecutando cleanup...")
            // Aquí va tu lógica
        }
    }
}
```

## Ejemplo 3: Worker con shutdown (controlado)

```go
type CleanupWorker struct {
    done chan bool  // Channel para apagarlo
}

func NewCleanupWorker() *CleanupWorker {
    return &CleanupWorker{
        done: make(chan bool),
    }
}

func (w *CleanupWorker) Start() {
    go func() {
        ticker := time.NewTicker(24 * time.Hour)
        defer ticker.Stop()

        for {
            select {
            case <-ticker.C:
                w.cleanup()
            case <-w.done:  // Cuando llaman a Stop()
                fmt.Println("Worker apagado")
                return
            }
        }
    }()
}

func (w *CleanupWorker) Stop() {
    w.done <- true  // Enviar señal de apagado
}

func (w *CleanupWorker) cleanup() {
    fmt.Println("Limpiando URLs expiradas...")
    // Tu lógica aquí
}
```

**Uso:**
```go
func main() {
    worker := NewCleanupWorker()
    worker.Start()
    
    // El servidor corre...
    time.Sleep(10 * time.Second)
    
    // Apagar limpiamente
    worker.Stop()
}
```

## Ejemplo 4: Worker con context (recomendado)

```go
type CleanupWorker struct {
    cancel context.CancelFunc
}

func NewCleanupWorker() *CleanupWorker {
    return &CleanupWorker{}
}

func (w *CleanupWorker) Start() {
    ctx, cancel := context.WithCancel(context.Background())
    w.cancel = cancel

    go func() {
        ticker := time.NewTicker(24 * time.Hour)
        defer ticker.Stop()

        for {
            select {
            case <-ticker.C:
                w.cleanup(ctx)
            case <-ctx.Done():
                fmt.Println("Worker apagado")
                return
            }
        }
    }()
}

func (w *CleanupWorker) Stop() {
    w.cancel()
}

func (w *CleanupWorker) cleanup(ctx context.Context) {
    // Pasar context a las queries de DB
    // Si el worker se apaga durante cleanup, las queries se cancelan
    fmt.Println("Limpiando URLs expiradas...")
}
```

## Ejemplo 5: Worker que hace cleanup a medianoche

```go
func (w *CleanupWorker) Start() {
    ctx, cancel := context.WithCancel(context.Background())
    w.cancel = cancel

    go func() {
        for {
            // Calcular cuánto falta para medianoche
            now := time.Now()
            next := time.Date(
                now.Year(), now.Month(), now.Day()+1,
                0, 0, 0, 0,  // Medianoche
                now.Location(),
            )
            duration := time.Until(next)
            
            fmt.Printf("Próximo cleanup en: %v\n", duration)
            
            // Esperar hasta medianoche
            select {
            case <-time.After(duration):
                w.cleanup(ctx)
            case <-ctx.Done():
                fmt.Println("Worker apagado")
                return
            }
        }
    }()
}
```

## Ejemplo 6: Múltiples workers

```go
func main() {
    // Worker de cleanup
    cleanupWorker := NewCleanupWorker()
    cleanupWorker.Start()
    
    // Worker de estadísticas
    statsWorker := NewStatsWorker()
    statsWorker.Start()
    
    // Worker de notificaciones
    notifWorker := NewNotifWorker()
    notifWorker.Start()
    
    // Esperar señal de apagado
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    // Apagar todos limpiamente
    cleanupWorker.Stop()
    statsWorker.Stop()
    notifWorker.Stop()
}
```

## Ejemplo 7: En main.go con servidor

```go
func main() {
    // 1. Config
    cfg := config.Load()
    
    // 2. Database
    db := setupDB(cfg)
    defer db.Close()
    
    // 3. Workers
    cleanupWorker := worker.NewCleanupWorker(db)
    cleanupWorker.Start()
    defer cleanupWorker.Stop()
    
    // 4. Servidor HTTP
    router := setupRouter(cfg, db)
    server := &http.Server{
        Addr:    ":" + cfg.Port,
        Handler: router,
    }
    
    // Arrancar servidor en goroutine
    go func() {
        if err := server.ListenAndServe(); err != nil {
            log.Fatal(err)
        }
    }()
    
    fmt.Println("Servidor en puerto", cfg.Port)
    
    // 5. Esperar señal de apagado
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    fmt.Println("Apagando...")
    cleanupWorker.Stop()
    server.Shutdown(context.Background())
}
```

## Resumen visual

```
main.go
    │
    ├──► Servidor HTTP (goroutine)
    │         │
    │         └── Escucha requests
    │
    ├──► Worker 1 (goroutine)
    │         │
    │         └── Duerme → ejecuta → repite
    │
    ├──► Worker 2 (goroutine)
    │         │
    │         └── Duerme → ejecuta → repite
    │
    └──► main() espera señal de apagado
```

## ¿Cuándo usar workers?

| Tarea | ¿Worker? |
|-------|----------|
| Limpiar URLs expiradas cada día | ✅ Sí |
| Enviar emails pendientes cada hora | ✅ Sí |
| Actualizar estadísticas cada 5 min | ✅ Sí |
| Responder a un request HTTP | ❌ No (usa handler) |
| Procesar una imagen subida | ❌ No (usa goroutine puntual) |

## Regla de oro

**Worker = tarea repetitiva en background.**

Si algo debe hacerse periódicamente sin que nadie lo pida, es un worker.
