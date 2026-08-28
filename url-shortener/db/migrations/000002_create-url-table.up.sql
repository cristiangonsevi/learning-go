CREATE TABLE urls (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    url TEXT NOT NULL,  -- VARCHAR(255) es muy corto para URLs largas
    short_url VARCHAR(100) NOT NULL UNIQUE,  -- Debe ser UNIQUE
    visit_count BIGINT DEFAULT 0 NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    last_visited_at TIMESTAMPTZ,
    user_id UUID,  -- Si tienes autenticación
    is_active BOOLEAN DEFAULT TRUE NOT NULL,
    expires_at TIMESTAMPTZ,  -- Para URLs con expiración
    title TEXT,  -- Título de la página (opcional)
    description TEXT  -- Descripción (opcional)
);

-- Índices para mejorar rendimiento
CREATE INDEX idx_urls_short_url ON urls(short_url);
CREATE INDEX idx_urls_user_id ON urls(user_id);
CREATE INDEX idx_urls_created_at ON urls(created_at);
CREATE INDEX idx_urls_is_active ON urls(is_active);

-- Trigger para actualizar updated_at automáticamente
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_urls_updated_at 
    BEFORE UPDATE ON urls 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();