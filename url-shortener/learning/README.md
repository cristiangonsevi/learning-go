# Guía de Aprendizaje Go - Concurrencia

## Orden de lectura

1. **[01-goroutines.md](./01-goroutines.md)** - Funciones que corren al mismo tiempo
2. **[02-channels.md](./02-channels.md)** - Cómo se comunican las goroutines
3. **[03-context.md](./03-context.md)** - El botón de cancelar
4. **[04-workers.md](./04-workers.md)** - Tareas repetitivas en background
5. **[05-ejemplo-practico.md](./05-ejemplo-practico.md)** - Todo junto en un proyecto

## Resumen rápido

```
goroutine  → "haz esto al mismo tiempo"
channel    → "habla con otra goroutine"
context    → "cancela cuando sea necesario"
worker     → "repite esta tarea cada X tiempo"
```

## Analogía del restaurante

```
main()           = Gerente (coordina todo)
goroutine        = Empleado (hace tareas en paralelo)
channel          = Timbre de pedido (comunicación entre empleados)
context          = Temporizador (cancela si algo tarda mucho)
worker           = Empleado de limpieza (limpia cada cierto tiempo)
```

## Ejemplo en tu proyecto URL Shortener

```go
main.go
    │
    ├──► Servidor HTTP (goroutine)
    │         │
    │         └── Responde a POST /api/urls
    │             Responde a GET /:code
    │
    ├──► Cleanup Worker (goroutine)
    │         │
    │         └── Cada 24h: desactiva URLs expiradas
    │
    └──► main() espera Ctrl+C
              │
              └── Apaga worker y servidor limpiamente
```
