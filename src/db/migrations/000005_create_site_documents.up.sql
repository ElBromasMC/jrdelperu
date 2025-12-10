BEGIN;

-- Table for site-wide documents (catalogs, brochures, etc.)
CREATE TABLE site_documents (
    document_id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    document_key varchar(50) UNIQUE NOT NULL,
    display_name varchar(255) NOT NULL,
    file_id int REFERENCES static_files ON DELETE SET NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

-- Index for faster lookups by key
CREATE INDEX idx_site_documents_key ON site_documents(document_key);

-- Seed fixed document slots
INSERT INTO site_documents (document_key, display_name) VALUES
('catalogo_vidrios', 'CATÁLOGO DE VIDRIOS - Corporación JR Del Perú S.A.C'),
('catalogo_aluminios', 'CATÁLOGO DE ALUMINIO - Corporación JR Del Perú S.A.C'),
('catalogo_upvc', 'CATÁLOGO DE uPVC - Corporación JR Del Perú S.A.C'),
('afiche_upvc', 'AFICHE DE LOS BENEFICIOS DEL uPVC - Corporación JR Del Perú S.A.C'),
('brochure_empresa', 'BROCHURE DE LA EMPRESA - Corporación JR Del Perú S.A.C');

COMMIT;
