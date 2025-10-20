
BEGIN;

-- Tipo ENUM para materiales
CREATE TYPE material_type AS ENUM ('vidrio', 'aluminio', 'upvc');

-- Tabla de archivos estáticos (imágenes y PDFs)
CREATE TABLE IF NOT EXISTS static_files (
    file_id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    file_name varchar(255) UNIQUE NOT NULL,
    file_type varchar(10),
    mime_type varchar(100),
    file_size_bytes bigint,
    created_at timestamptz NOT NULL DEFAULT NOW()
);

-- Tabla de administradores para autenticación
CREATE TABLE IF NOT EXISTS admins (
    admin_id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    username varchar(100) UNIQUE NOT NULL,
    email varchar(255) UNIQUE NOT NULL,
    password_hash varchar(255) NOT NULL,
    is_active boolean DEFAULT true,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

-- Tabla para almacenar mensajes del formulario de contacto
CREATE TABLE IF NOT EXISTS contact_submissions (
    submission_id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    full_name varchar(255) NOT NULL,
    email varchar(255) NOT NULL,
    phone varchar(50),
    subject varchar(500) NOT NULL,
    message text NOT NULL,
    is_read boolean DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT NOW()
);

-- Tabla de etiquetas para categorías
CREATE TABLE IF NOT EXISTS categories_tags (
    tag_id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    tag_name varchar(255) UNIQUE NOT NULL,
    position_num int NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

-- Tabla de categorías de productos
CREATE TABLE IF NOT EXISTS categories (
    category_id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    material_type material_type NOT NULL,
    slug varchar(255) NOT NULL,
    name varchar(255) NOT NULL,
    description text NOT NULL,
    long_description text NOT NULL,
    image_id int REFERENCES static_files ON DELETE SET NULL,
    tag_id int REFERENCES categories_tags ON DELETE SET NULL,
    pdf_id int REFERENCES static_files ON DELETE SET NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    UNIQUE (material_type, slug)
);

-- Tabla de características técnicas de categorías
CREATE TABLE IF NOT EXISTS category_features (
    feature_id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    category_id int NOT NULL REFERENCES categories ON DELETE CASCADE,
    name varchar(255) NOT NULL,
    description varchar(255) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    UNIQUE (category_id, name)
);

-- Tabla de productos individuales (items)
CREATE TABLE IF NOT EXISTS items (
    item_id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    category_id int NOT NULL REFERENCES categories ON DELETE CASCADE,
    slug varchar(255) NOT NULL,
    name varchar(255) NOT NULL,
    description text NOT NULL,
    long_description text NOT NULL,
    image_id int REFERENCES static_files ON DELETE SET NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    UNIQUE (category_id, slug)
);

-- Tabla de galería de imágenes por item
CREATE TABLE IF NOT EXISTS item_images (
    item_id int NOT NULL REFERENCES items ON DELETE CASCADE,
    image_id int NOT NULL REFERENCES static_files ON DELETE CASCADE,
    position_num int NOT NULL,
    PRIMARY KEY (item_id, image_id)
);

-- Índices para optimización de consultas
CREATE INDEX idx_static_files_file_type ON static_files(file_type);
CREATE INDEX idx_admins_username ON admins(username);
CREATE INDEX idx_admins_email ON admins(email);
CREATE INDEX idx_contact_submissions_created_at ON contact_submissions(created_at DESC);
CREATE INDEX idx_contact_submissions_is_read ON contact_submissions(is_read);
CREATE INDEX idx_categories_tags_position ON categories_tags(position_num);
CREATE INDEX idx_categories_material_type ON categories(material_type);
CREATE INDEX idx_categories_slug ON categories(material_type, slug);
CREATE INDEX idx_items_slug ON items(category_id, slug);
CREATE INDEX idx_item_images_position ON item_images(item_id, position_num);

COMMIT;

