BEGIN;

-- Libro de Reclamaciones (Código de Protección y Defensa del Consumidor)
CREATE TABLE IF NOT EXISTS complaints (
    complaint_id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    -- Datos del consumidor reclamante
    full_name varchar(255) NOT NULL,
    document_number varchar(50) NOT NULL,
    address varchar(500),
    phone varchar(50),
    email varchar(255),
    -- Identificación del bien contratado
    good_type varchar(20) NOT NULL,           -- 'producto' | 'servicio'
    good_description text NOT NULL,
    -- Datos de la reclamación
    claim_type varchar(20) NOT NULL,          -- 'reclamo' | 'queja'
    detail text NOT NULL,                     -- Detalle de la reclamación
    request text NOT NULL,                    -- Pedido del consumidor
    -- Observaciones de la empresa (uso interno, editable desde el admin)
    company_notes text NOT NULL DEFAULT '',
    -- Fecha de registro declarada por el consumidor
    registered_at date NOT NULL DEFAULT CURRENT_DATE,
    -- Seguimiento interno
    is_resolved boolean NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_complaints_created_at ON complaints(created_at DESC);
CREATE INDEX idx_complaints_is_resolved ON complaints(is_resolved, created_at DESC);

COMMIT;
