# Context - El botón de cancelar

## ¿Qué es un context?

Es un **control remoto** que te permite cancelar operaciones.

## Analogía

```
Sin context:
────────────────────────────────────────
Llamas a una API
La API no responde
Esperas... esperas... esperas... (eternamente)
Tu app se cuelga

Con context:
────────────────────────────────────────
Llamas a una API
Pones alarma de 5 segundos
La API no responde
Alarma suena → cancelas la petición
Tu app sigue funcionando
```

## Los 4 tipos de context

```go
// 1. Vacío (raíz de todo)
ctx := context.Background()

// 2. Con timeout (se cancela después de X tiempo)
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// 3. Con deadline (se cancela en fecha/hora específica)
deadline := time.Now().Add(5 * time.Second)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()

// 4. Cancelable (lo cancelas tú manualmente)
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
```

## Ejemplo 1: Timeout (el más usado)

```go
func llamarAPI(url string) (string, error) {
    // Crear context con 3 segundos de timeout
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()  // ← SIEMPRE defer cancel()

    // Crear request con el context
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return "", err
    }

    // Ejecutar request
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        // ¿Se acabó el tiempo?
        if ctx.Err() == context.DeadlineExceeded {
            return "", fmt.Errorf("la API tardó más de 3 segundos")
        }
        return "", err
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)
    return string(body), nil
}
```

**¿Qué pasa?**
```
1. Creamos context con 3s timeout
2. Hacemos request a la API
3. Si la API responde en < 3s → OK
4. Si la API tarda > 3s → se cancela automáticamente
```

## Ejemplo 2: Cancelar manualmente

```go
func main() {
    ctx, cancel := context.WithCancel(context.Background())

    // Goroutine que trabaja
    go func() {
        for {
            select {
            case <-ctx.Done():  // ← Se ejecuta cuando cancel() es llamado
                fmt.Println("Worker cancelado:", ctx.Err())
                return
            default:
                fmt.Println("Trabajando...")
                time.Sleep(500 * time.Millisecond)
            }
        }
    }()

    // Después de 2 segundos, cancelar
    time.Sleep(2 * time.Second)
    cancel()  // ← Esto envía señal al channel ctx.Done()

    time.Sleep(100 * time.Millisecond) // Esperar a que el worker termine
}
```

**Salida:**
```
Trabajando...
Trabajando...
Trabajando...
Trabajando...
Worker cancelado: context canceled
```

## Ejemplo 3: Query a base de datos

```go
func (r *Repo) GetUser(ctx context.Context, id int) (*User, error) {
    var user User
    query := `SELECT id, name FROM users WHERE id = $1`
    
    // QueryRowContext usa el context para cancelar si es necesario
    err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// En el handler (Gin ya tiene context con timeout)
func (h *Handler) GetUser(c *gin.Context) {
    // c.Request.Context() se cancela si el cliente cierra la conexión
    user, err := h.repo.GetUser(c.Request.Context(), 123)
}
```

## Ejemplo 4: Worker cancelable

```go
type Worker struct {
    cancel context.CancelFunc
}

func NewWorker() *Worker {
    return &Worker{}
}

func (w *Worker) Start() {
    ctx, cancel := context.WithCancel(context.Background())
    w.cancel = cancel

    go func() {
        ticker := time.NewTicker(1 * time.Second)
        defer ticker.Stop()

        for {
            select {
            case <-ticker.C:
                fmt.Println("Ejecutando tarea...")
            case <-ctx.Done():
                fmt.Println("Worker apagado")
                return
            }
        }
    }()
}

func (w *Worker) Stop() {
    w.cancel()  // ← Esto envía señal al channel ctx.Done()
}
```

## Ejemplo 5: Pasar valores (poco recomendado)

```go
// Definir tipo para evitar colisiones
type contextKey string

const userIDKey contextKey = "userID"

// Guardar valor
ctx := context.WithValue(context.Background(), userIDKey, 123)

// Leer valor
userID, ok := ctx.Value(userIDKey).(int)
if !ok {
    fmt.Println("No hay userID")
}
```

**Nota:** Mejor pasar datos como parámetros directos.

## El patrón que necesitas recordar

```go
// 1. Crear context
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()  // ← SIEMPRE defer cancel()

// 2. Pasar context a funciones
result, err := algunaFuncion(ctx)

// 3. Verificar si se canceló
if ctx.Err() != nil {
    fmt.Println("Se canceló:", ctx.Err())
}
```

## Resumen visual

```
context.Background()
       │
       ├─► WithTimeout(5s) ──► Se cancela a los 5s automáticamente
       │         │
       │         └─► ctx.Done() recibe señal
       │
       ├─► WithCancel() ──► Lo cancelas tú con cancel()
       │         │
       │         └─► ctx.Done() recibe señal
       │
       └─► WithValue(key, val) ──► Guarda datos (poco usado)
```

## ¿Cuándo usar cada uno?

| Contexto | Cuándo usarlo |
|----------|---------------|
| `Background()` | Inicio de la app, root del context tree |
| `WithTimeout()` | Llamadas a APIs, queries a DB |
| `WithCancel()` | Workers que necesitas apagar |
| `WithValue()` | Datos del request (userID) - usar poco |

## Regla de oro

**Context = control de vida de operaciones.**

Si algo puede tardar mucho o necesitas cancelarlo, usa context.
