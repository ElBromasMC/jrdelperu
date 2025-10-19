
BEGIN;
    
CREATE TYPE material_type AS ENUM ('vidrio', 'aluminio', 'upvc');

CREATE TABLE IF NOT EXISTS static_files (
    file_id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    file_name varchar(255) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS categories_tags (
    tag_id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    tag_name varchar(255) UNIQUE NOT NULL,
    position_num int NOT NULL
);

CREATE TABLE IF NOT EXISTS categories (
    category_id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    material_type material_type NOT NULL,
    slug varchar(255) NOT NULL,
    name varchar(255) NOT NULL,
    description text NOT NULL,
    long_description text NOT NULL,
    image_id int REFERENCES static_files ON DELETE SET NULL,
    tag_id int REFERENCES categories_tags ON DELETE SET NULL,
    UNIQUE (material_type, slug)
);

CREATE TABLE IF NOT EXISTS category_features (
    feature_id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    category_id int NOT NULL REFERENCES categories ON DELETE CASCADE,
    name varchar(255) NOT NULL,
    description varchar(255) NOT NULL,
    UNIQUE (category_id, name)
);

CREATE TABLE IF NOT EXISTS items (
    item_id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    category_id int NOT NULL REFERENCES categories ON DELETE CASCADE,
    slug varchar(255) NOT NULL,
    name varchar(255) NOT NULL,
    description text NOT NULL,
    long_description text NOT NULL,
    image_id int REFERENCES static_files ON DELETE SET NULL,
    UNIQUE (category_id, slug)
);

CREATE TABLE IF NOT EXISTS item_images (
    item_id int NOT NULL REFERENCES items ON DELETE CASCADE,
    image_id int NOT NULL REFERENCES static_files ON DELETE CASCADE,
    position_num int NOT NULL,
    PRIMARY KEY (item_id, image_id)
);

COMMIT;

