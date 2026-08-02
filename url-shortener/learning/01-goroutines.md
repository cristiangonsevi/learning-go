# Goroutines - El concepto más básico

## ¿Qué es una goroutine?

Es una función que corre **al mismo tiempo** que otras funciones.

## Analogía

```
Sin goroutine (normal):
───────────────────────────────
Hacer café    → 5 min
Hacer tostada → 3 min
TOTAL         → 8 min (uno después de otro)

Con goroutine:
───────────────────────────────
Hacer café    → 5 min ─┐
Hacer tostada → 3 min ─┤ al mismo tiempo
TOTAL         → 5 min (el más lento)
```

## Código

```go
// Función normal (secuencial)
func main() {
    hacerCafe()
    hacerTostada()
    fmt.Println("Listo!")
}

func hacerCafe() {
    fmt.Println("Haciendo café...")
    time.Sleep(5 * time.Second)
    fmt.Println("Café listo")
}

func hacerTostada() {
    fmt.Println("Haciendo tostada...")
    time.Sleep(3 * time.Second)
    fmt.Println("Tostada lista")
}

// Salida:
// Haciendo café...
// (5 segundos)
// Café listo
// Haciendo tostada...
// (3 segundos)
// Tostada lista
// Listo!
// TOTAL: 8 segundos
```

```go
// Con goroutine (paralelo)
func main() {
    go hacerCafe()      // ← Agrega "go" y corre en paralelo
    go hacerTostada()   // ← Agrega "go" y corre en paralelo
    
    time.Sleep(6 * time.Second) // Esperar a que terminen
    fmt.Println("Listo!")
}

// Salida:
// Haciendo café...
// Haciendo tostada...
// (3 segundos)
// Tostada lista
// (2 segundos más)
// Café listo
// Listo!
// TOTAL: 5 segundos
```

## La palabra mágica: `go`

```go
func saludar() {
    fmt.Println("Hola")
}

// Sin goroutine (espera a que termine)
saludar()

// Con goroutine (corre en paralelo, NO espera)
go saludar()
```

## Problema: ¿cómo espero a que terminen?

```go
// MALO - no sabes cuándo terminan
go hacerCafe()
go hacerTostada()
time.Sleep(10 * time.Second)  // ¿Cuánto espero? ¿Y si tarda más?

// BUENO - usar WaitGroup (lo explico en 02-channels.md)
```

## Ejemplo real

```go
// Enviar emails a 100 usuarios (secuencial - lento)
for _, user := range users {
    enviarEmail(user)  // 1 segundo cada uno = 100 segundos
}

// Enviar emails a 100 usuarios (goroutines - rápido)
var wg sync.WaitGroup
for _, user := range users {
    wg.Add(1)
    go func(u User) {
        defer wg.Done()
        enviarEmail(u)  // Todos al mismo tiempo = ~1 segundo
    }(user)
}
wg.Wait()
```

## Resumen

| Concepto | Explicación |
|----------|-------------|
| `go funcion()` | Ejecuta `funcion` en paralelo |
| Goroutine | Función que corre al mismo tiempo que otras |
| ¿Cuándo usar? | Cuando quieres hacer varias cosas a la vez |

## Regla de oro

**Si algo puede hacerse al mismo tiempo, usa `go`.**

Ejemplos:
- ✅ Enviar 100 emails (cada uno independiente)
- ✅ Consultar múltiples APIs (cada una independiente)
- ❌ Hacer café y luego servirlo (depende del primero)
