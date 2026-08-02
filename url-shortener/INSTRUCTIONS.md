# Plan Maestro - Acortador de URLs

Arquitectura, flujos y decisiones técnicas. Sin código, solo diseño.

---

## 🏗️ Arquitectura General

```
┌─────────────────────────────────────────────────────────────┐
│                      Cliente (Browser/API)                  │
└─────────────────────┬───────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────┐
│                      Load Balancer (opcional)               │
│                      (Nginx/HAProxy)                        │
└─────────────────────┬───────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────┐
│                   API Gateway (Gin)                        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐   │
│  │ Middleware  │→ │  Router     │→ │  Rate Limiter   │   │
│  │ (Logger,    │  │  (Endpoints)│  │  (Token Bucket) │   │
│  │  Recovery,  │  │             │  │                 │   │
│  │  CORS)      │  │             │  │                 │   │
│  └─────────────┘  └─────────────┘  └─────────────────┘   │
└─────────────────────┬─────────────────────────────────────┘
                      │
        ┌─────────────┼─────────────┐
        ▼             ▼             ▼
┌──────────────┐ ┌──────────┐ ┌─────────────┐
│  PostgreSQL  │ │   Redis  │ │  Background │
│  (Datos)     │ │  (Caché) │ │  Worker     │
└──────────────┘ └──────────┘ └─────────────┘
```

---

## 📁 Estructura del Proyecto

```
url-shortener/
├── cmd/
│   └── server/
│       └── main.go              # Punto de entrada
├── internal/
│   ├── config/
│   │   └── config.go            # Carga de variables de entorno
│   ├── handler/
│   │   ├── auth.go              # Handlers de autenticación
│   │   ├── url.go               # Handlers de URLs
│   │   └── health.go            # Health check
│   ├── middleware/
│   │   ├── auth.go              # JWT middleware
│   │   ├── cors.go              # CORS middleware
│   │   ├── logger.go            # Request logging
│   │   ├── ratelimit.go         # Rate limiting
│   │   └── recovery.go          # Panic recovery
│   ├── model/
│   │   ├── user.go              # Modelo de usuario
│   │   ├── url.go               # Modelo de URL
│   │   └── visit.go             # Modelo de visita
│   ├── repository/
│   │   ├── user_repo.go         # Interfaz + implementación
│   │   ├── url_repo.go          # Interfaz + implementación
│   │   └── visit_repo.go        # Interfaz + implementación
│   ├── service/
│   │   ├── auth_service.go      # Lógica de autenticación
│   │   ├── url_service.go       # Lógica de URLs
│   │   └── stats_service.go     # Lógica de estadísticas
│   ├── cache/
│   │   └── redis.go             # Cliente Redis
│   ├── worker/
│   │   └── cleanup.go           # Worker de limpieza
│   └── router/
│       └── router.go            # Definición de rutas
├── migrations/
│   ├── 001_create_users.up.sql
│   ├── 001_create_users.down.sql
│   ├── 002_create_urls.up.sql
│   ├── 002_create_urls.down.sql
│   ├── 003_create_visits.up.sql
│   └── 003_create_visits.down.sql
├── docker-compose.yml
├── Dockerfile
├── .env.example
├── go.mod
└── go.sum
```

---

## 📦 Dependencias (Librerías a usar)

| Propósito | Librería | Comando de instalación |
|-----------|----------|------------------------|
| HTTP Framework | `github.com/gin-gonic/gin` | `go get github.com/gin-gonic/gin` |
| SQL Extensions | `github.com/jmoiron/sqlx` | `go get github.com/jmoiron/sqlx` |
| PostgreSQL Driver | `github.com/lib/pq` | `go get github.com/lib/pq` |
| Redis Client | `github.com/redis/go-redis/v9` | `go get github.com/redis/go-redis/v9` |
| JWT | `github.com/golang-jwt/jwt/v5` | `go get github.com/golang-jwt/jwt/v5` |
| Config (env) | `github.com/joho/godotenv` | `go get github.com/joho/godotenv` |
| Migrations | `github.com/golang-migrate/migrate/v4` | `go get github.com/golang-migrate/migrate/v4` |
| UUID | `github.com/google/uuid` | `go get github.com/google/uuid` |
| Validación | `github.com/go-playground/validator/v10` | `go get github.com/go-playground/validator/v10` |
| Logging | `log/slog` (stdlib) | No requiere instalación |
| Password Hash | `golang.org/x/crypto` | `go get golang.org/x/crypto` |

---

## 🔧 Uso de sqlx (Patrones de Repositorio)

### Conexión a la base de datos

```go
// En main.go o config
db, err := sqlx.Connect("postgres", "host=localhost port=5432 user=urlshortener dbname=urlshortener password=secret sslmode=disable")
if err != nil {
    log.Fatal(err)
}
defer db.Close()
```

### Patrón de Repository con sqlx

```go
// Definir interfaz
type URLRepository interface {
    Create(ctx context.Context, url *model.URL) error
    GetByCode(ctx context.Context, code string) (*model.URL, error)
    GetByUserID(ctx context.Context, userID int64, page, limit int) ([]model.URL, error)
    Update(ctx context.Context, url *model.URL) error
    Delete(ctx context.Context, code string) error
}

// Implementación
type urlRepository struct {
    db *sqlx.DB
}

func NewURLRepository(db *sqlx.DB) URLRepository {
    return &urlRepository{db: db}
}
```

### Ejemplos de queries con sqlx

```go
// INSERT
func (r *urlRepository) Create(ctx context.Context, url *model.URL) error {
    query := `INSERT INTO urls (code, original_url, user_id, is_public, expires_at)
              VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
    return r.db.QueryRowContext(ctx, query,
        url.Code, url.OriginalURL, url.UserID, url.IsPublic, url.ExpiresAt,
    ).Scan(&url.ID, &url.CreatedAt)
}

// SELECT single
func (r *urlRepository) GetByCode(ctx context.Context, code string) (*model.URL, error) {
    var url model.URL
    query := `SELECT id, code, original_url, user_id, is_public, expires_at, created_at, updated_at
              FROM urls WHERE code = $1`
    err := r.db.GetContext(ctx, &url, query, code)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil // No encontrado
        }
        return nil, err
    }
    return &url, nil
}

// SELECT many
func (r *urlRepository) GetByUserID(ctx context.Context, userID int64, page, limit int) ([]model.URL, error) {
    var urls []model.URL
    offset := (page - 1) * limit
    query := `SELECT id, code, original_url, user_id, is_public, expires_at, created_at, updated_at
              FROM urls WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
    err := r.db.SelectContext(ctx, &urls, query, userID, limit, offset)
    return urls, err
}

// UPDATE
func (r *urlRepository) Update(ctx context.Context, url *model.URL) error {
    query := `UPDATE urls SET original_url = $1, is_public = $2, expires_at = $3, updated_at = NOW()
              WHERE id = $4`
    _, err := r.db.ExecContext(ctx, query, url.OriginalURL, url.IsPublic, url.ExpiresAt, url.ID)
    return err
}

// DELETE
func (r *urlRepository) Delete(ctx context.Context, code string) error {
    query := `DELETE FROM urls WHERE code = $1`
    _, err := r.db.ExecContext(ctx, query, code)
    return err
}
```

---

## 📊 Modelo de Datos (Esquema UML)

### Tablas principales

```
┌─────────────────┐          ┌─────────────────────────┐
│     users        │          │         urls             │
├─────────────────┤          ├─────────────────────────┤
│ id: BIGINT PK   │◄─────────│ id: BIGINT PK           │
│ email: TEXT UNIQ│          │ code: VARCHAR(10) UNIQ  │
│ password: TEXT  │          │ original_url: TEXT      │
│ created_at: TS  │          │ user_id: BIGINT FK      │
│ updated_at: TS  │          │ is_public: BOOLEAN      │
└─────────────────┘          │ expires_at: TIMESTAMP?  │
                             │ created_at: TIMESTAMP   │
                             │ updated_at: TIMESTAMP   │
                             └─────────────────────────┘
                                      │
                                      │ 1
                                      │
                                      │ *
┌─────────────────┐          ┌─────────────────────────┐
│  url_stats      │          │    visit_logs           │
├─────────────────┤          ├─────────────────────────┤
│ url_id: BIGINT  │◄─────────│ id: BIGINT PK           │
│ visits: INT     │          │ url_id: BIGINT FK       │
│ last_visit: TS  │          │ visitor_ip: TEXT        │
└─────────────────┘          │ user_agent: TEXT        │
                             │ referer: TEXT           │
                             │ country: VARCHAR(50)    │
                             │ visited_at: TIMESTAMP   │
                             └─────────────────────────┘
```

### Relaciones

1. **users → urls**: Un usuario puede tener muchas URLs (1:N). Si user_id es NULL → URL anónima
2. **urls → url_stats**: Cada URL tiene una estadística (1:1)
3. **urls → visit_logs**: Una URL puede tener muchos visit logs (1:N)

### SQL de Migraciones

```sql
-- 001_create_users.up.sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 002_create_urls.up.sql
CREATE TABLE urls (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(10) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    is_public BOOLEAN DEFAULT true,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_urls_code ON urls(code);
CREATE INDEX idx_urls_user_id ON urls(user_id);
CREATE INDEX idx_urls_expires_at ON urls(expires_at);

-- 003_create_visits.up.sql
CREATE TABLE url_stats (
    url_id BIGINT PRIMARY KEY REFERENCES urls(id) ON DELETE CASCADE,
    visits INT DEFAULT 0,
    last_visit TIMESTAMP
);

CREATE TABLE visit_logs (
    id BIGSERIAL PRIMARY KEY,
    url_id BIGINT REFERENCES urls(id) ON DELETE CASCADE,
    visitor_ip VARCHAR(45),
    user_agent TEXT,
    referer TEXT,
    country VARCHAR(50),
    visited_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_visit_logs_url_id ON visit_logs(url_id);
CREATE INDEX idx_visit_logs_visited_at ON visit_logs(visited_at);
```

### SQL de Rollback (Down Migrations)

```sql
-- 003_create_visits.down.sql
DROP TABLE IF EXISTS visit_logs;
DROP TABLE IF EXISTS url_stats;

-- 002_create_urls.down.sql
DROP TABLE IF EXISTS urls;

-- 001_create_users.down.sql
DROP TABLE IF EXISTS users;
```

### Cómo ejecutar migraciones

```bash
# Instalar CLI de migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Ejecutar migraciones
migrate -path migrations -database "postgres://urlshortener:secure_password@localhost:5432/urlshortener?sslmode=disable" up

# Rollback última migración
migrate -path migrations -database "postgres://urlshortener:secure_password@localhost:5432/urlshortener?sslmode=disable" down 1

# Ver versión actual
migrate -path migrations -database "postgres://urlshortener:secure_password@localhost:5432/urlshortener?sslmode=disable" version

# Forzar a una versión (si algo falla)
migrate -path migrations -database "postgres://urlshortener:secure_password@localhost:5432/urlshortener?sslmode=disable" force 3
```

---

## 🔌 API Endpoints (Referencia Completa)

| Método | Ruta | Auth | Descripción | Request Body | Response |
|--------|------|------|-------------|--------------|----------|
| `POST` | `/api/auth/register` | No | Registrar usuario | `{email, password}` | `{token, user}` |
| `POST` | `/api/auth/login` | No | Login | `{email, password}` | `{token, user}` |
| `POST` | `/api/urls` | Opcional | Crear URL corta | `{url, expires_in?, is_public?}` | `{code, short_url, expires_at}` |
| `GET` | `/api/urls` | Sí | Listar mis URLs | - | `{urls: [...]}` |
| `GET` | `/api/urls/:code` | No | Obtener info de URL | - | `{code, original_url, visits, ...}` |
| `PUT` | `/api/urls/:code` | Sí | Editar URL | `{original_url?, is_public?, expires_in?}` | `{url: {...}}` |
| `DELETE` | `/api/urls/:code` | Sí | Eliminar URL | - | `{message: "deleted"}` |
| `GET` | `/api/urls/:code/stats` | Opcional | Estadísticas | - | `{visits, last_visit, countries, ...}` |
| `GET` | `/:code` | No | Redirigir | - | `302 Redirect` |
| `GET` | `/health` | No | Health check | - | `{status: "ok"}` |

---

## 🔄 Flujo de Redirección (GET /:code)

```
┌──────────────────────────────────────────────────────────────┐
│                    FLUJO DE REDIRECCIÓN                      │
└──────────────────────────────────────────────────────────────┘

1. Usuario → GET /abc123

2. ¿Está en caché Redis?
   ├─ Sí → Obtener URL de Redis → Ir al paso 5
   └─ No → Ir al paso 3

3. Buscar en PostgreSQL por código
   ├─ ¿Existe? → Ir al paso 4
   └─ No existe → Error 404

4. Guardar en Redis (TTL: 1 hora)
   ✅ Cachear la URL

5. ¿La URL ha expirado?
   ├─ Sí → Error 410 Gone
   └─ No → Ir al paso 6

6. Incrementar contador de visitas (async)
   📝 Registrar en tabla visit_logs
   📝 Incrementar url_stats

7. Redirigir (302 Found)
   🔄 Location: original_url
```

### Detalles importantes del flujo

- **Redis TTL**: 1 hora para URLs populares. Si una URL tiene muchos hits, se mantiene en caché
- **Escritura asíncrona**: Las visitas se registran en background para no bloquear la redirección
- **Códigos HTTP**:
  - 302 Found: Redirección temporal (puedes cambiar a 301 si quieres permanente)
  - 404 Not Found: URL no existe
  - 410 Gone: URL expirada

---

## 📝 Flujo de Creación de URL (POST /api/urls)

```
┌──────────────────────────────────────────────────────────────┐
│                   FLUJO DE CREACIÓN DE URL                   │
└──────────────────────────────────────────────────────────────┘

1. Cliente → POST /api/urls
   Body: { "url": "https://...", "expires_in": 7, "is_public": true }

2. Validar URL (sintaxis)
   ├─ ¿URL válida? → Ir al paso 3
   └─ No válida → Error 400

3. Validar URL (existencia real)
   ✅ Hacer HEAD/GET a la URL para verificar que existe
   ├─ Si responde 200/301/302 → Ir al paso 4
   └─ Si no responde → Error 422 (No se puede acceder)

4. ¿Usuario está autenticado?
   ├─ Sí → Usuario registrado
   │   ├─ expires_in: NULL → URL permanente
   │   └─ expires_in: N → Calcula fecha de expiración
   │
   └─ No → Usuario anónimo
       └─ expires_at: NOW() + 7 días

5. Generar código único (6-8 caracteres)
   🔑 Algoritmo: Base62 (0-9, A-Z, a-z)
   - Intentar 10 veces si hay colisión
   - Si todas fallan, generar con timestamp

6. Guardar en PostgreSQL
   📝 Insertar en tabla urls
   📝 Crear registro en url_stats (visits: 0)

7. Opcional: Guardar en Redis
   ✅ Si es URL pública, cachear inmediatamente

8. Devolver respuesta
   📤 { "code": "abc123", "short_url": "...", "expires_at": "..." }
```

### Algoritmo de generación de código

```
INPUT:  ID del usuario (o nil)
OUTPUT: Código único de 6-8 caracteres

1. Tomar timestamp actual (nanosegundos)
2. Mezclar con un número aleatorio
3. Codificar en Base62 (alfanumérico)
4. Si el código ya existe → repetir con diferente random
5. Si después de 10 intentos → usar UUID corto
```

---

## 📋 Flujo de Listar Mis URLs (GET /api/urls)

```
┌──────────────────────────────────────────────────────────────┐
│                   FLUJO DE LISTAR URLs                       │
└──────────────────────────────────────────────────────────────┘

1. Usuario autenticado → GET /api/urls
   Query params opcionales: ?page=1&limit=20&sort=created_at&order=desc

2. Verificar JWT token
   ├─ Token inválido → Error 401
   └─ Token válido → Ir al paso 3

3. Buscar URLs del usuario en PostgreSQL
   SELECT * FROM urls
   WHERE user_id = ?
   ORDER BY created_at DESC
   LIMIT ? OFFSET ?

4. Contar total de URLs del usuario
   SELECT COUNT(*) FROM urls WHERE user_id = ?

5. Devolver respuesta paginada
   📤 {
        "urls": [...],
        "total": 45,
        "page": 1,
        "limit": 20,
        "has_next": true
      }
```

---

## ✏️ Flujo de Editar URL (PUT /api/urls/:code)

```
┌──────────────────────────────────────────────────────────────┐
│                    FLUJO DE EDITAR URL                       │
└──────────────────────────────────────────────────────────────┘

1. Usuario autenticado → PUT /api/urls/:code
   Body: { "original_url": "https://nueva-url.com", "is_public": false }

2. Verificar JWT token
   ├─ Token inválido → Error 401
   └─ Token válido → Ir al paso 3

3. Buscar URL por código en PostgreSQL
   ├─ No existe → Error 404
   └─ Existe → Ir al paso 4

4. Verificar propiedad
   ├─ ¿user_id coincide con el token? → Ir al paso 5
   └─ No es el dueño → Error 403

5. Validar nueva URL (si se proporciona)
   ├─ ¿URL válida? → Ir al paso 6
   └─ No válida → Error 400

6. Actualizar campos en PostgreSQL
   - original_url (si se proporcionó)
   - is_public (si se proporcionó)
   - expires_at (si se proporcionó expires_in)
   - updated_at = NOW()

7. Invalidar caché en Redis
   🗑️ Eliminar key url:{code}

8. Devolver respuesta
   📤 { "url": { "code": "...", "original_url": "...", ... } }
```

---

## 🗑️ Flujo de Eliminar URL (DELETE /api/urls/:code)

```
┌──────────────────────────────────────────────────────────────┐
│                   FLUJO DE ELIMINAR URL                      │
└──────────────────────────────────────────────────────────────┘

1. Usuario autenticado → DELETE /api/urls/:code

2. Verificar JWT token
   ├─ Token inválido → Error 401
   └─ Token válido → Ir al paso 3

3. Buscar URL por código en PostgreSQL
   ├─ No existe → Error 404
   └─ Existe → Ir al paso 4

4. Verificar propiedad
   ├─ ¿user_id coincide con el token? → Ir al paso 5
   └─ No es el dueño → Error 403

5. Eliminar de PostgreSQL
   🗑️ DELETE FROM urls WHERE code = ? (CASCADE elimina stats y visits)

6. Eliminar de Redis
   🗑️ Eliminar url:{code} y url:{code}:stats

7. Devolver respuesta
   📤 { "message": "URL deleted successfully" }
```

---

## 🔐 Flujo de Autenticación

```
┌──────────────────────────────────────────────────────────────┐
│                    REGISTRO DE USUARIO                      │
└──────────────────────────────────────────────────────────────┘

1. POST /api/auth/register
   Body: { "email": "...", "password": "..." }

2. Validar email (formato y unicidad)
   ├─ ¿Ya existe? → Error 409
   └─ No existe → Ir al paso 3

3. Hashear password (bcrypt)
   🔐 Cost: 10 (balance entre seguridad y rendimiento)

4. Guardar en PostgreSQL

5. Generar JWT (24h de expiración)
   📝 Payload: { "user_id": id, "email": email }

6. Devolver: { "token": "...", "user": { ... } }

┌──────────────────────────────────────────────────────────────┐
│                     LOGIN DE USUARIO                        │
└──────────────────────────────────────────────────────────────┘

1. POST /api/auth/login
   Body: { "email": "...", "password": "..." }

2. Buscar usuario por email
   ├─ ¿Existe? → Ir al paso 3
   └─ No existe → Error 401

3. Verificar password (bcrypt compare)
   ├─ ¿Coincide? → Ir al paso 4
   └─ No coincide → Error 401

4. Generar JWT (24h expiración)

5. Devolver: { "token": "...", "user": { ... } }
```

### Estructura del JWT

```
Header:
{
  "alg": "HS256",
  "typ": "JWT"
}

Payload:
{
  "user_id": 123,
  "email": "user@example.com",
  "exp": 1234567890,
  "iat": 1234567890
}
```

---

## 🧹 Flujo de Limpieza de URLs Expiradas

```
┌──────────────────────────────────────────────────────────────┐
│               CLEANUP WORKER (Background)                    │
└──────────────────────────────────────────────────────────────┘

1. Worker se ejecuta cada hora

2. Query: SELECT id, code, expires_at
   FROM urls
   WHERE expires_at IS NOT NULL
   AND expires_at < NOW()

3. Por cada URL expirada:
   a. Eliminar de PostgreSQL
   b. Eliminar de Redis (si existe)
   c. Registrar en log de auditoría

4. ¿Muchas URLs expiradas? → Ejecutar en batches (1000 por lote)

5. Dormir 1 hora → Repetir
```

---

## 🗂️ Esquema de Redis (Caché)

### Estructura de datos en Redis

```
┌──────────────────────────────────────────────────────────────┐
│                     REDIS KEY STRUCTURE                     │
└──────────────────────────────────────────────────────────────┘

1. Key: url:{code}
   Value: JSON de la URL
   TTL: 3600 segundos (1 hora)
   Ejemplo: url:abc123 → {"id":1,"original":"...","expires_at":"..."}

2. Key: url:{code}:stats
   Value: Estadísticas en caché
   TTL: 300 segundos (5 minutos)
   Ejemplo: url:abc123:stats → {"visits":1234,"last_visit":"..."}

3. Key: rate:{ip}
   Value: Contador de requests
   TTL: 60 segundos
   Ejemplo: rate:192.168.1.1 → 45 (para rate limiting)
```

### Estrategia de caché

| Acción | Qué hacer |
|--------|-----------|
| **GET /:code** | 1. Buscar en Redis<br>2. Si no está, buscar en PostgreSQL<br>3. Guardar en Redis (TTL 1h) |
| **POST /api/urls** | 1. Guardar en PostgreSQL<br>2. Si es pública, guardar en Redis |
| **PUT /api/urls/:code** | 1. Actualizar en PostgreSQL<br>2. Invalidar caché (eliminar de Redis) |
| **DELETE /api/urls/:code** | 1. Eliminar de PostgreSQL<br>2. Eliminar de Redis |
| **Visita a URL** | 1. Incrementar en PostgreSQL (async)<br>2. Invalidar stats en Redis |

---

## 📊 Flujo de Estadísticas

```
┌──────────────────────────────────────────────────────────────┐
│                    ESTADÍSTICAS DE URL                       │
└──────────────────────────────────────────────────────────────┘

1. Usuario (autenticado) → GET /api/urls/:code/stats

2. Verificar permisos:
   a. ¿La URL es pública? → Ver estadísticas
   b. ¿Es el dueño? → Ver estadísticas detalladas
   c. ¿Anónimo y privada? → Error 403

3. Buscar en Redis (stats)
   ├─ Sí → Devolver datos cacheados
   └─ No → Ir al paso 4

4. Consultar en PostgreSQL:
   SELECT * FROM url_stats WHERE url_id = ?
   SELECT * FROM visit_logs WHERE url_id = ? ORDER BY visited_at DESC LIMIT 100

5. Calcular métricas:
   - Total de visitas
   - Última visita
   - Visitas por día (última semana)
   - Países (top 5)
   - Referers (top 5)

6. Guardar en Redis (TTL: 5 min)

7. Devolver JSON con estadísticas
```

### Ejemplo de respuesta de estadísticas

```json
{
  "code": "abc123",
  "original_url": "https://example.com/very-long-url",
  "total_visits": 1234,
  "last_visit": "2024-01-15T10:30:00Z",
  "created_at": "2024-01-01T00:00:00Z",
  "expires_at": null,
  "visits_by_day": [
    {"date": "2024-01-15", "count": 45},
    {"date": "2024-01-14", "count": 38}
  ],
  "top_countries": [
    {"country": "Mexico", "count": 500},
    {"country": "USA", "count": 300}
  ],
  "top_referrers": [
    {"referer": "twitter.com", "count": 200},
    {"referer": "direct", "count": 150}
  ]
}
```

---

## ⚙️ Rate Limiting

### Estrategia: Token Bucket por IP

```
┌──────────────────────────────────────────────────────────────┐
│                    RATE LIMITING FLOW                        │
└──────────────────────────────────────────────────────────────┘

1. Request llega → Extraer IP del cliente

2. Buscar en Redis: rate:{ip}
   ├─ No existe → Crear con TTL 60s, contador = 1
   └─ Existe → Ir al paso 3

3. ¿Contador < límite?
   ├─ Sí → Incrementar contador, permitir request
   └─ No → Error 429 Too Many Requests

4. Headers de respuesta:
   X-RateLimit-Limit: 100
   X-RateLimit-Remaining: 55
   X-RateLimit-Reset: 1234567890
```

### Límites por tipo

| Tipo de usuario | Requests por minuto |
|-----------------|---------------------|
| Anónimo | 30 |
| Registrado | 100 |
| Admin | Sin límite |

---

## 🎨 Diagrama de Secuencia: Creación de URL (Usuario Registrado)

```
Cliente          API (Gin)         PostgreSQL        Redis
   │                │                   │               │
   │ POST /api/urls │                   │               │
   │───────────────►│                   │               │
   │                │                   │               │
   │                │ 1. Validar URL    │               │
   │                │ (sintaxis)        │               │
   │                │                   │               │
   │                │ 2. Verificar      │               │
   │                │ autenticación     │               │
   │                │                   │               │
   │                │ 3. Calcular exp.  │               │
   │                │ (NULL = perm)     │               │
   │                │                   │               │
   │                │ 4. Generar código │               │
   │                │ (Base62 + random) │               │
   │                │                   │               │
   │                │ 5. INSERT url     │               │
   │                │──────────────────►│               │
   │                │                   │               │
   │                │ 6. INSERT stats   │               │
   │                │──────────────────►│               │
   │                │                   │               │
   │                │ 7. Cachear (si   │               │
   │                │    es pública)    │               │
   │                │──────────────────────────────────►│
   │                │                   │               │
   │ 8. Response    │                   │               │
   │◄───────────────│                   │               │
```

---

## 🎨 Diagrama de Secuencia: Redirección

```
Cliente          API (Gin)         Redis        PostgreSQL     Worker
   │                │                 │               │           │
   │ GET /abc123    │                 │               │           │
   │───────────────►│                 │               │           │
   │                │                 │               │           │
   │                │ 1. Buscar en   │               │           │
   │                │ Redis          │               │           │
   │                │────────────────►│               │           │
   │                │                 │               │           │
   │                │ 2. ¿Cache miss? │               │           │
   │                │◄────────────────│               │           │
   │                │                 │               │           │
   │                │ 3. Buscar en    │               │           │
   │                │ PostgreSQL      │               │           │
   │                │────────────────────────────────►│           │
   │                │                 │               │           │
   │                │ 4. Guardar en   │               │           │
   │                │ Redis (TTL 1h)  │               │           │
   │                │────────────────►│               │           │
   │                │                 │               │           │
   │                │ 5. ¿URL expiró? │               │           │
   │                │ → 410 Gone      │               │           │
   │                │                 │               │           │
   │                │ 6. Registrar    │               │           │
   │                │ visita (async)  │               │           │
   │                │────────────────────────────────────────────────►│
   │                │                 │               │           │
   │ 7. 302 Found   │                 │               │           │
   │◄───────────────│                 │               │           │
   │                │                 │               │           │
   │                │                 │               │  8. Log   │
   │                │                 │               │  visita   │
   │                │                 │               │◄──────────│
```

---

## ❌ Manejo de Errores

### Formato estándar de error

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "URL not found",
    "details": "No URL exists with code 'abc123'"
  }
}
```

### Códigos de error HTTP

| HTTP Code | Error Code | Cuándo usar |
|-----------|------------|-------------|
| 400 | `BAD_REQUEST` | Body inválido, parámetros faltantes |
| 400 | `INVALID_URL` | URL con formato inválido |
| 401 | `UNAUTHORIZED` | Token faltante o inválido |
| 403 | `FORBIDDEN` | Sin permisos para acceder al recurso |
| 404 | `NOT_FOUND` | URL o recurso no encontrado |
| 409 | `CONFLICT` | Email ya registrado |
| 410 | `GONE` | URL expirada |
| 422 | `UNPROCESSABLE` | URL no accesible (no responde) |
| 429 | `RATE_LIMITED` | Demasiadas requests |
| 500 | `INTERNAL_ERROR` | Error del servidor |

---

## 🚀 Estrategia de Despliegue

### Arquitectura de producción

```
┌─────────────────────────────────────────────────────────────┐
│                   DOCKER COMPOSE                            │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────────┐   ┌──────────────┐   ┌──────────────┐  │
│  │   API (Go)   │   │ PostgreSQL   │   │    Redis     │  │
│  │   Port:8080  │   │   Port:5432  │   │   Port:6379  │  │
│  └──────────────┘   └──────────────┘   └──────────────┘  │
│         │                                                   │
│         ▼                                                   │
│  ┌──────────────┐                                          │
│  │    Nginx     │ (Proxy inverso + SSL)                    │
│  │   Port:443   │                                          │
│  └──────────────┘                                          │
└─────────────────────────────────────────────────────────────┘
```

### Variables de entorno

```
# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=urlshortener
DB_PASSWORD=secure_password
DB_NAME=urlshortener

# Redis
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT
JWT_SECRET=super_secret_key_32_bytes_min

# Server
PORT=8080
CORS_ALLOWED_ORIGINS=https://example.com

# API
BASE_URL=https://tu-dominio.com
DEFAULT_EXPIRATION_DAYS=7
MAX_EXPIRATION_DAYS=30

# Rate Limiting
RATE_LIMIT_ANONYMOUS=30
RATE_LIMIT_AUTHENTICATED=100
```

### Ejemplo de docker-compose.yml

```yaml
version: '3.8'

services:
  api:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      - postgres
      - redis
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=urlshortener
      - DB_PASSWORD=secure_password
      - DB_NAME=urlshortener
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - JWT_SECRET=${JWT_SECRET}
      - PORT=8080

  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: urlshortener
      POSTGRES_PASSWORD: secure_password
      POSTGRES_DB: urlshortener
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

volumes:
  postgres_data:
  redis_data:
```

### Ejemplo de Dockerfile

```dockerfile
# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/server/main.go

# Run stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations
EXPOSE 8080
CMD ["./main"]
```

### Comandos útiles de Docker

```bash
# Levantar servicios
docker-compose up -d

# Ver logs
docker-compose logs -f api

# Detener servicios
docker-compose down

# Reconstruir después de cambios
docker-compose up -d --build

# Ejecutar migraciones dentro del contenedor
docker-compose exec api ./main migrate up
```

---

## 📋 Checklist de funcionalidades

### MVP (Mínimo Producto Viable) - Semana 1-2

- [ ] Servidor HTTP con Gin
- [ ] Endpoint POST /api/urls (crear URL)
- [ ] Endpoint GET /:code (redirección)
- [ ] Conexión a PostgreSQL
- [ ] Generación de códigos (Base62)
- [ ] Expiración de 7 días para anónimos
- [ ] Logs básicos

### Funcionalidades core - Semana 3

- [ ] Autenticación JWT (registro/login)
- [ ] Endpoint GET /api/urls (listar mis URLs)
- [ ] Endpoint PUT /api/urls/:code (editar)
- [ ] Endpoint DELETE /api/urls/:code (eliminar)
- [ ] Expiración configurable para usuarios
- [ ] URLs públicas/privadas

### Mejoras de rendimiento - Semana 4

- [ ] Redis para caché de URLs
- [ ] Rate limiting por IP
- [ ] Visitas registradas async (goroutine)
- [ ] Cleanup worker (URLs expiradas)
- [ ] Endpoint GET /api/urls/:code/stats

### Extras (si hay tiempo)

- [ ] Tests unitarios e integración
- [ ] Swagger/OpenAPI documentation
- [ ] Docker + docker-compose
- [ ] CI/CD (GitHub Actions)
- [ ] Frontend simple (HTML + CSS + JS)
- [ ] QR code para cada URL

---

## 📊 Diagrama de Casos de Uso

```
┌─────────────────────────────────────────────────────────────┐
│                  CASOS DE USO                               │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─────────────────┐        ┌──────────────────────┐      │
│  │   Usuario       │        │   Usuario            │      │
│  │   Anónimo       │        │   Registrado         │      │
│  └────────┬────────┘        └──────────┬───────────┘      │
│           │                             │                   │
│           ▼                             ▼                   │
│  ┌─────────────────┐        ┌──────────────────────┐      │
│  │ • Crear URL     │        │ • Crear URL          │      │
│  │   (expira 7d)   │        │   (permanente/conf)  │      │
│  │ • Redirigir     │        │ • Redirigir          │      │
│  │ • Ver estadíst. │        │ • Ver estadísticas   │      │
│  │   (si pública)  │        │ • Listar mis URLs    │      │
│  └─────────────────┘        │ • Editar URLs        │      │
│                              │ • Eliminar URLs      │      │
│                              └──────────────────────┘      │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## 🎯 Resumen ejecutivo para reclutadores

**¿Qué hace este proyecto?**
Un acortador de URLs con autenticación, estadísticas y expiración configurable.

**¿Qué tecnologías usa?**
Go (Gin), PostgreSQL, Redis, Docker, JWT.

**¿Qué problemas resuelve?**
- URLs largas y difíciles de compartir
- URLs temporales para casos de uso específicos
- Seguimiento de audiencia con estadísticas detalladas
- Control de acceso (público/privado)

**¿Qué patrones de diseño usa?**
- Repository pattern (abstracción de DB)
- Dependency injection
- Middleware pattern (autenticación, logging, rate limiting)
- Worker pattern (cleanup de URLs expiradas)
- Caching pattern (Redis)

**¿Qué habilidades demuestra?**
- Backend con Go (Gin)
- Bases de datos SQL (PostgreSQL)
- Caché (Redis)
- Concurrencia (goroutines, channels)
- Autenticación (JWT)
- Arquitectura limpia y escalable
- Testing y documentación

---

## 🚀 Orden de Implementación

1. **Configurar el proyecto** - Crear estructura de carpetas y go.mod
2. **Docker Compose** - PostgreSQL + Redis + API
3. **Migraciones** - Crear tablas en PostgreSQL
4. **Config** - Cargar variables de entorno
5. **Modelos** - Definir structs de Go
6. **Repository** - Capa de acceso a datos
7. **Service** - Lógica de negocio
8. **Handlers** - Endpoints HTTP
9. **Router** - Definir rutas y middleware
10. **Middleware** - Auth, CORS, Logger, Rate Limit
11. **Cache** - Integrar Redis
12. **Worker** - Cleanup de URLs expiradas
13. **Tests** - Unitarios e integración
14. **Documentación** - Swagger/OpenAPI
15. **Deploy** - Dockerfile + CI/CD
