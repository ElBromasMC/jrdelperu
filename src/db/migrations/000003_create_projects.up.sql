
BEGIN;

-- Tabla de proyectos completados
CREATE TABLE IF NOT EXISTS projects (
    project_id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    slug varchar(255) UNIQUE NOT NULL,
    description varchar(255) NOT NULL,  -- Nombre del proyecto (ej: "Hotel Real Palace")
    location varchar(255) NOT NULL,     -- Ubicación (ej: "Churín - Oyón")
    period varchar(50) NOT NULL,        -- Periodo (ej: "2019")
    area_m2 numeric(10,2),              -- Metros cuadrados (ej: 214.00)
    service text NOT NULL,              -- Descripción del servicio
    display_order int NOT NULL DEFAULT 0,
    is_visible boolean DEFAULT true,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

-- Tabla de imágenes asociadas a proyectos
CREATE TABLE IF NOT EXISTS project_images (
    project_id int NOT NULL REFERENCES projects ON DELETE CASCADE,
    image_id int NOT NULL REFERENCES static_files ON DELETE CASCADE,
    display_order int NOT NULL DEFAULT 0,
    is_featured boolean DEFAULT false,
    PRIMARY KEY (project_id, image_id)
);

-- Índices para optimización de consultas
CREATE INDEX idx_projects_slug ON projects(slug);
CREATE INDEX idx_projects_display_order ON projects(display_order);
CREATE INDEX idx_projects_is_visible ON projects(is_visible);
CREATE INDEX idx_project_images_display_order ON project_images(project_id, display_order);
CREATE INDEX idx_project_images_featured ON project_images(project_id, is_featured);

COMMIT;
