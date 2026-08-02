# Channels - Cómo se comunican las goroutines

## ¿Qué es un channel?

Es un **tubo** por donde pasas datos entre goroutines.

## Analogía

```
Goroutine A                    Goroutine B
    │                              │
    │   "Hola, aquí tienes datos"  │
    │ ─────── channel ──────────► │
    │                              │
```

Es como un tubo de banco:
- Tú metes dinero por un lado (goroutine A)
- El cajero recibe por el otro lado (goroutine B)

## Código básico

```go
// Crear un channel
ch := make(chan string)

// Enviar datos (meter en el tubo)
ch <- "hola"  // ← Esto BLOQUEA hasta que alguien reciba

// Recibir datos (sacar del tubo)
mensaje := <-ch  // ← Esto BLOQUEA hasta que alguien envíe
```

## Ejemplo 1: Dos goroutines hablando

```go
func main() {
    ch := make(chan string)  // Crear channel

    go func() {
        time.Sleep(2 * time.Second)
        ch <- "Ya terminé!"  // Enviar mensaje
    }()

    fmt.Println("Esperando...")
    mensaje := <-ch  // Recibir mensaje (bloquea hasta que llegue)
    fmt.Println(mensaje)
}

// Salida:
// Esperando...
// (2 segundos)
// Ya terminé!
```

## Ejemplo 2: WaitGroup con channels

```go
func main() {
    done := make(chan bool)  // Channel para avisar que terminó

    go func() {
        fmt.Println("Haciendo café...")
        time.Sleep(3 * time.Second)
        done <- true  // Avisar que terminé
    }()

    <-done  // Esperar a que avise
    fmt.Println("Café listo!")
}
```

## Ejemplo 3: Múltiples mensajes

```go
func main() {
    ch := make(chan string)

    // Productor (envía 3 mensajes)
    go func() {
        ch <- "mensaje 1"
        ch <- "mensaje 2"
        ch <- "mensaje 3"
        close(ch)  // Cerrar channel cuando termine
    }()

    // Consumidor (recibe hasta que se cierre)
    for msg := range ch {
        fmt.Println(msg)
    }
}

// Salida:
// mensaje 1
// mensaje 2
// mensaje 3
```

## Ejemplo 4: Channel con buffer

```go
// Sin buffer (bloquea hasta que alguien reciba)
ch := make(chan string)
ch <- "hola"  // ← BLOQUEA aquí si nadie recibe

// Con buffer (puedes enviar sin que reciban inmediatamente)
ch := make(chan string, 3)  // Guarda hasta 3 mensajes
ch <- "mensaje 1"  // No bloquea
ch <- "mensaje 2"  // No bloquea
ch <- "mensaje 3"  // No bloquea
ch <- "mensaje 4"  // ← BLOQUEA porque el buffer está lleno
```

## Ejemplo 5: select (escuchar múltiples channels)

```go
func main() {
    ch1 := make(chan string)
    ch2 := make(chan string)

    go func() {
        time.Sleep(2 * time.Second)
        ch1 <- "Canal 1"
    }()

    go func() {
        time.Sleep(1 * time.Second)
        ch2 <- "Canal 2"
    }()

    // select espera a que CUALQUIER channel reciba datos
    select {
    case msg := <-ch1:
        fmt.Println("Recibí de ch1:", msg)
    case msg := <-ch2:
        fmt.Println("Recibí de ch2:", msg)
    }
}

// Salida (ch2 llega primero):
// Recibí de ch2: Canal 2
```

## Ejemplo 6: select con timeout

```go
func main() {
    ch := make(chan string)

    go func() {
        time.Sleep(5 * time.Second)
        ch <- "resultado"
    }()

    select {
    case msg := <-ch:
        fmt.Println("Recibí:", msg)
    case <-time.After(3 * time.Second):  // Timeout de 3 segundos
        fmt.Println("Se acabó el tiempo!")
    }
}

// Salida:
// Se acabó el tiempo!
```

## Resumen visual

```
make(chan T)        → Crear channel de tipo T
make(chan T, 5)     → Crear channel con buffer de 5

ch <- valor         → Enviar valor al channel
valor := <-ch       → Recibir valor del channel
close(ch)           → Cerrar channel (no más envíos)

for range ch       → Recibir hasta que se cierre

select {
case v := <-ch1:    → Escuchar ch1
case v := <-ch2:    → Escuchar ch2
case <-time.After:  → Timeout
}
```

## ¿Cuándo usar qué?

| Situación | Solución |
|-----------|----------|
| Goroutine necesita avisar que terminó | `done <- true` |
| Goroutine necesita enviar datos | `ch <- datos` |
| Escuchar múltiples sources | `select` |
| Timeout | `select` con `time.After` |
| Buffer de mensajes | `make(chan T, n)` |

## Regla de oro

**Channel = comunicación entre goroutines.**

Si dos goroutines necesitan hablar, usan un channel.
