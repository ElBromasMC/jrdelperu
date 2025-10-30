# JRDELPERU - Sistema de Gestión de Productos y Servicios

Sistema web para gestión de productos de vidrio, aluminio y uPVC con panel administrativo y formulario de contacto.

## Tecnologías Utilizadas

- **Backend**: Go 1.23+ con Echo v4 framework
- **Base de datos**: PostgreSQL con sqlc para generación de código
- **Templating**: Templ (tipo-seguro, compila a Go)
- **Frontend**: HTMX 2.x + Tailwind CSS v4 + DaisyUI
- **Optimización de Imágenes**: h2non/bimg (libvips) con conversión automática a WebP
- **Email**: wneessen/go-mail
- **Build**: Makefile (orquesta Templ, Tailwind, esbuild, sqlc)

## Dependencias del Sistema

### Requeridas

- **Go 1.23+**
- **PostgreSQL** (versión 12 o superior)
- **Node.js** (para Tailwind CSS y esbuild)
- **libvips 8.12+** - Biblioteca de procesamiento de imágenes
  - Ubuntu/Debian: `sudo apt-get install libvips-dev`
  - macOS: `brew install vips`
  - Alpine Linux: `apk add vips-dev gcc musl-dev`
- **GCC** - Compilador C (requerido para CGO y bimg)
  - Ubuntu/Debian: `sudo apt-get install build-essential`
  - Alpine Linux: `apk add gcc musl-dev`

### Opcionales

- **migrate** - Para aplicar migraciones de base de datos
  - Instalación: https://github.com/golang-migrate/migrate

## Estructura del Proyecto

```
.
├── assets/              # Archivos estáticos embebidos en el binario
│   ├── static/
│   │   ├── css/        # CSS generado por Tailwind
│   │   ├── js/         # JavaScript bundleado por esbuild
│   │   └── img/        # Imágenes estáticas
│   ├── embed.go        # Embeds para producción
│   └── embed_dev.go    # Embeds para desarrollo
├── cmd/server/         # Punto de entrada del servidor
├── config/             # Configuración y datos de ejemplo
├── db/
│   ├── migrations/     # Migraciones SQL (versionadas)
│   └── query/          # Consultas SQL para sqlc
├── repository/         # Código Go generado por sqlc
├── service/            # Lógica de negocio
├── model/              # Modelos de dominio
├── handler/            # Handlers HTTP de Echo
├── view/               # Templates Templ
├── Makefile           # Orquesta el build
├── sqlc.yml           # Configuración de sqlc
└── tailwind.config.cjs # Configuración de Tailwind
```

## Modelo de Datos

### Tablas Principales

#### 1. **admins** - Administradores del sistema
Almacena usuarios con permisos administrativos.

| Campo          | Tipo         | Descripción                    |
|----------------|--------------|--------------------------------|
| admin_id       | int (PK)     | ID autogenerado                |
| username       | varchar(100) | Nombre de usuario único        |
| email          | varchar(255) | Email único                    |
| password_hash  | varchar(255) | Hash bcrypt de la contraseña   |
| is_active      | boolean      | Estado activo/inactivo         |
| created_at     | timestamptz  | Fecha de creación              |
| updated_at     | timestamptz  | Última actualización           |

#### 2. **contact_submissions** - Mensajes del formulario de contacto
Almacena todos los mensajes enviados por usuarios.

| Campo          | Tipo         | Descripción                    |
|----------------|--------------|--------------------------------|
| submission_id  | int (PK)     | ID autogenerado                |
| full_name      | varchar(255) | Nombre completo                |
| email          | varchar(255) | Email de contacto              |
| phone          | varchar(50)  | Teléfono (opcional)            |
| subject        | varchar(500) | Asunto del mensaje             |
| message        | text         | Mensaje completo               |
| is_read        | boolean      | Marca de leído/no leído        |
| created_at     | timestamptz  | Fecha de envío                 |

#### 3. **static_files** - Archivos estáticos (imágenes, PDFs)
Almacena referencias a archivos subidos.

| Campo           | Tipo         | Descripción                           |
|-----------------|--------------|---------------------------------------|
| file_id         | int (PK)     | ID autogenerado                       |
| file_name       | varchar(255) | Nombre único del archivo (hash+ext)   |
| file_type       | varchar(10)  | 'image' o 'pdf'                       |
| mime_type       | varchar(100) | Tipo MIME                             |
| file_size_bytes | bigint       | Tamaño en bytes                       |
| display_name    | varchar(255) | Nombre descriptivo para mostrar al usuario |
| created_at      | timestamptz  | Fecha de subida                       |

#### 4. **categories_tags** - Etiquetas para categorías
Tags para organizar categorías (ej: "Destacado", "Nuevo").

| Campo        | Tipo         | Descripción                    |
|--------------|--------------|--------------------------------|
| tag_id       | int (PK)     | ID autogenerado                |
| tag_name     | varchar(255) | Nombre único del tag           |
| position_num | int          | Orden de visualización         |
| created_at   | timestamptz  | Fecha de creación              |
| updated_at   | timestamptz  | Última actualización           |

#### 5. **categories** - Categorías de productos
Categorías organizadas por tipo de material.

| Campo            | Tipo         | Descripción                        |
|------------------|--------------|------------------------------------|
| category_id      | int (PK)     | ID autogenerado                    |
| material_type    | enum         | 'vidrio', 'aluminio', 'upvc'       |
| slug             | varchar(255) | Slug URL (único por material_type) |
| name             | varchar(255) | Nombre de la categoría             |
| description      | text         | Descripción corta                  |
| long_description | text         | Descripción detallada              |
| image_id         | int (FK)     | Imagen principal (nullable)        |
| tag_id           | int (FK)     | Tag asociado (nullable)            |
| pdf_id           | int (FK)     | PDF técnico (nullable)             |
| created_at       | timestamptz  | Fecha de creación                  |
| updated_at       | timestamptz  | Última actualización               |

#### 6. **category_features** - Características técnicas de categorías
Especificaciones técnicas de cada categoría.

| Campo       | Tipo         | Descripción                    |
|-------------|--------------|--------------------------------|
| feature_id  | int (PK)     | ID autogenerado                |
| category_id | int (FK)     | Categoría asociada             |
| name        | varchar(255) | Nombre de la característica    |
| description | varchar(255) | Descripción/valor              |
| created_at  | timestamptz  | Fecha de creación              |
| updated_at  | timestamptz  | Última actualización           |

#### 7. **items** - Productos individuales
Items/productos dentro de cada categoría.

| Campo            | Tipo         | Descripción                        |
|------------------|--------------|------------------------------------|
| item_id          | int (PK)     | ID autogenerado                    |
| category_id      | int (FK)     | Categoría asociada                 |
| slug             | varchar(255) | Slug URL (único por categoría)     |
| name             | varchar(255) | Nombre del producto                |
| description      | text         | Descripción corta                  |
| long_description | text         | Descripción detallada              |
| image_id         | int (FK)     | Imagen principal (nullable)        |
| created_at       | timestamptz  | Fecha de creación                  |
| updated_at       | timestamptz  | Última actualización               |

#### 8. **item_images** - Galería de imágenes por item
Múltiples imágenes para cada producto.

| Campo        | Tipo     | Descripción                    |
|--------------|----------|--------------------------------|
| item_id      | int (FK) | Item asociado                  |
| image_id     | int (FK) | Imagen en static_files         |
| position_num | int      | Orden de visualización         |

**PK compuesta**: (item_id, image_id)

### Relaciones Principales

- `categories.material_type` → ENUM('vidrio', 'aluminio', 'upvc')
- `categories.image_id` → `static_files.file_id` (ON DELETE SET NULL)
- `categories.tag_id` → `categories_tags.tag_id` (ON DELETE SET NULL)
- `categories.pdf_id` → `static_files.file_id` (ON DELETE SET NULL)
- `category_features.category_id` → `categories.category_id` (ON DELETE CASCADE)
- `items.category_id` → `categories.category_id` (ON DELETE CASCADE)
- `items.image_id` → `static_files.file_id` (ON DELETE SET NULL)
- `item_images.item_id` → `items.item_id` (ON DELETE CASCADE)
- `item_images.image_id` → `static_files.file_id` (ON DELETE CASCADE)

### Índices para Optimización

- `idx_contact_submissions_created_at` - Listado cronológico de mensajes
- `idx_contact_submissions_is_read` - Filtrado de mensajes no leídos
- `idx_admins_username`, `idx_admins_email` - Login de administradores
- `idx_static_files_file_type` - Filtrado por tipo de archivo
- `idx_categories_material_type` - Filtrado por tipo de material
- `idx_categories_slug` - Búsqueda por slug
- `idx_items_slug` - Búsqueda de items por slug
- `idx_item_images_position` - Ordenamiento de galerías

## Panel de Administración

El sistema incluye un panel administrativo completo con las siguientes funcionalidades:

### Gestión de Contenido

- **Categorías**
  - Crear, editar y eliminar categorías de productos
  - Organizar por tipo de material (Vidrio, Aluminio, uPVC)
  - Asignar imágenes principales y PDFs técnicos
  - Agregar etiquetas (Destacado, Nuevo, etc.)
  - Definir características técnicas por categoría

- **Productos (Items)**
  - Crear productos dentro de cada categoría
  - Generación automática de slugs desde el nombre
  - Asignar imagen principal
  - Descripciones cortas y largas
  - Galería de imágenes por producto

- **Etiquetas**
  - Crear y gestionar tags para categorías
  - Ordenar por posición
  - Edición inline con HTMX

### Gestión de Archivos

- **Subida de Imágenes**
  - Formatos soportados: JPG, PNG, WebP, GIF
  - Tamaño máximo: 10 MB
  - **Optimización automática**: Conversión a WebP con calidad 80
  - **Redimensionamiento**: Máximo 2000x2000px manteniendo aspecto
  - Nombres de visualización personalizables

- **Subida de PDFs**
  - Tamaño máximo: 20 MB
  - Nombres de visualización personalizables

- **Gestión**
  - Renombrar archivos con edición inline HTMX
  - Eliminar archivos con confirmación
  - Contador de archivos en tiempo real

### Dashboard

- **Estadísticas en tiempo real**
  - Total de categorías
  - Total de productos
  - Mensajes sin leer
  - Total de archivos (imágenes + PDFs)

### Características Técnicas

- **Interfaz completamente HTMX**: Sin JavaScript personalizado
- **Actualizaciones en tiempo real**: Out of Band (OOB) swaps
- **Validación de formularios**: Client-side y server-side
- **Diseño responsive**: Mobile-first con DaisyUI

## Optimización de Imágenes

El sistema implementa optimización automática de imágenes al momento de la subida:

### Configuración (service/file.go)

```go
const (
    MaxImageWidth    = 2000  // Ancho máximo en píxeles
    MaxImageHeight   = 2000  // Alto máximo en píxeles
    WebPQuality      = 80    // Calidad WebP (1-100)
    StripMetadata    = true  // Eliminar datos EXIF
)
```

### Proceso

1. **Validación**: Verifica tipo MIME y tamaño máximo
2. **Conversión**: Todas las imágenes se convierten a formato WebP
3. **Redimensionamiento**: Reduce dimensiones si exceden límites (mantiene aspecto)
4. **Compresión**: Aplica calidad WebP 80 para balance tamaño/calidad
5. **Limpieza**: Elimina metadatos EXIF para privacidad y menor tamaño

### Beneficios

- **Reducción de tamaño**: 30-80% menor que PNG/JPG
- **Carga más rápida**: Mejor experiencia de usuario
- **SEO**: Tiempos de carga optimizados
- **Ancho de banda**: Menor consumo de recursos

## Configuración de Rutas de Archivos

Todas las rutas de archivos subidos están centralizadas en `config/path.go` para facilitar el mantenimiento:

### Configuración (config/path.go)

```go
const (
    // Rutas base para uploads
    UPLOADS_PATH    = "/uploads"       // URL pública base
    UPLOADS_SAVEDIR = "./uploads"      // Directorio en filesystem

    // Rutas para imágenes
    IMAGES_PATH    = "/uploads/images"    // URL pública de imágenes
    IMAGES_SAVEDIR = "./uploads/images"   // Directorio de imágenes

    // Rutas para PDFs
    PDFS_PATH    = "/uploads/pdfs"      // URL pública de PDFs
    PDFS_SAVEDIR = "./uploads/pdfs"     // Directorio de PDFs
)
```

### Uso en el Código

Estas constantes se utilizan en toda la aplicación:

- **Handlers**: Para servir archivos estáticos (`cmd/server/main.go`)
- **Services**: Para generar URLs de archivos (`service/file.go`)
- **Templates**: Para construir rutas en las vistas (`view/*.templ`)

**IMPORTANTE**: Nunca usar rutas hardcoded. Siempre importar y usar las constantes de `alc/config`.

```go
import (
    "alc/config"
    "path"  // Para URLs (siempre usa /)
)

// Ejemplo: Construir URL de imagen
imageURL := path.Join(config.IMAGES_PATH, fileName)
```

## Variables de Entorno

### Requeridas

- **ENV**: Modo de ejecución
  - `development` - Modo desarrollo con hot reload
  - `production` - Modo producción con assets minificados

- **POSTGRESQL_URL**: URL de conexión a PostgreSQL
  - Formato: `postgres://user:password@host:port/database?sslmode=disable`

- **REL**: Número de release (para cache busting en producción)
  - Ejemplo: `1`, `2`, etc.

- **SESSION_SECRET**: Clave secreta para firmar cookies de sesión
  - Debe ser una cadena aleatoria de al menos 32 caracteres
  - Ejemplo: `mi-clave-super-secreta-de-32-chars-o-mas`
  - IMPORTANTE: Usar una clave diferente en producción

### Para Seeder (Creación de Usuario Admin Inicial)

- **ADMIN_USERNAME**: Nombre de usuario del administrador inicial
  - Ejemplo: `admin`
- **ADMIN_EMAIL**: Email del administrador inicial
  - Ejemplo: `admin@jrdelperu.com`
- **ADMIN_PASSWORD**: Contraseña del administrador inicial
  - Ejemplo: `password123`
  - IMPORTANTE: Usar una contraseña segura en producción

### Para Email (Formulario de Contacto)

- **SMTP_HOST**: Servidor SMTP (ej: `smtp.gmail.com`)
- **SMTP_PORT**: Puerto SMTP (ej: `587` para STARTTLS, `465` para SSL/TLS)
- **SMTP_USERNAME**: Usuario para autenticación SMTP
- **SMTP_PASSWORD**: Contraseña o app password
- **SMTP_FROM_EMAIL**: Email del remitente
- **SMTP_FROM_NAME**: Nombre del remitente
- **SMTP_TO_EMAIL**: Email destinatario para mensajes del formulario

### Para Google reCAPTCHA (opcional)

- **RECAPTCHA_SITE_KEY**: Clave pública de reCAPTCHA
- **RECAPTCHA_SECRET_KEY**: Clave secreta de reCAPTCHA

### Ejemplo de Configuración

```bash
# Aplicación
ENV="development"
REL="1"
SESSION_SECRET="mi-clave-super-secreta-de-32-chars-o-mas"

# Base de datos
POSTGRESQL_URL="postgres://postgres:password@localhost:5432/jrdelperu?sslmode=disable"

# Seeder (para crear usuario admin inicial)
ADMIN_USERNAME="admin"
ADMIN_EMAIL="admin@jrdelperu.com"
ADMIN_PASSWORD="password123"


# Email
SMTP_HOST="smtp.gmail.com"
SMTP_PORT="587"
SMTP_USERNAME="tu-email@gmail.com"
SMTP_PASSWORD="tu-app-password"
SMTP_FROM_EMAIL="noreply@jrdelperu.com"
SMTP_FROM_NAME="JR del Perú"
SMTP_TO_EMAIL="contacto@jrdelperu.com"

# reCAPTCHA (opcional)
RECAPTCHA_SITE_KEY="tu-site-key"
RECAPTCHA_SECRET_KEY="tu-secret-key"
```

## Comandos de Build

**IMPORTANTE**: Los builds requieren `CGO_ENABLED=1` debido a la optimización de imágenes con libvips.

### Desarrollo
```bash
# Build en modo desarrollo (con hot reload)
CGO_ENABLED=1 ENV=development make build/server

# Live reload (watch mode)
CGO_ENABLED=1 ENV=development make live
```

### Producción
```bash
# Build optimizado para producción
CGO_ENABLED=1 ENV=production make build/server
```

### Limpieza
```bash
# Limpiar archivos generados
make clean
```

## Carga de Datos de Prueba

El proyecto incluye un programa seeder que carga datos de ejemplo en la base de datos.

### 1. Crear Usuario Administrador

Primero, crea un usuario administrador para acceder al panel admin:

```bash
# Compilar el seeder de admin
go build -o build/seeder ./cmd/seeder

# Ejecutar (requiere variables de entorno)
POSTGRESQL_URL="postgres://user:pass@localhost:5432/jrdelperu?sslmode=disable" \
ADMIN_USERNAME="admin" \
ADMIN_EMAIL="admin@jrdelperu.com" \
ADMIN_PASSWORD="tu-password-seguro" \
./build/seeder
```

El seeder verificará si el usuario ya existe antes de crearlo.

### 2. Cargar Datos de Ejemplo (Categorías e Items)

Para cargar datos de ejemplo de vidrios, aluminios y uPVC:

```bash
# Compilar el seeder de datos
go build -o build/seed-data ./cmd/seed-data

# Ejecutar (requiere POSTGRESQL_URL)
POSTGRESQL_URL="postgres://user:pass@localhost:5432/jrdelperu?sslmode=disable" \
./build/seed-data
```

Esto cargará:
- **3 Categorías de Vidrios** (Monolítico, Reflectivo, Decorativo)
  - 3 Items para Vidrio Monolítico (2-3mm, 4-5mm, 6-8mm)
  - 3 Características técnicas para Vidrio Monolítico
- **3 Categorías de Aluminios** (Sistema de Fachadas, Carpintería Técnica, Perimetral)
  - 3 Items para Sistema de Fachadas (Estándar, Stick, Frame)
- **3 Categorías de uPVC** (Lumina 60, Lumina 104, Natura 66)
  - 3 Items para Lumina 60 (Proyectante, Batiente, Oscilobatiente)

**Nota:** Los datos de ejemplo se encuentran en `config/sample.go` y serán eliminados una vez que el panel admin esté completamente funcional.

### 3. Acceder al Panel Admin

Una vez cargados los datos:

1. Inicia el servidor: `ENV=development ./build/server`
2. Accede a: `http://localhost:8080/admin/login`
3. Inicia sesión con el usuario admin creado
4. Gestiona categorías, items y características desde el panel

## Migraciones de Base de Datos

Las migraciones están en `db/migrations/` y deben ejecutarse en orden:

1. `000001_initialize_schema` - Schema inicial con todas las tablas
2. `000002_add_display_name_to_files` - Agrega campo display_name a static_files

### Aplicar Migraciones

```bash
# Usando migrate CLI (recomendado)
migrate -path db/migrations -database "$POSTGRESQL_URL" up

# O usando psql directamente (ejecutar en orden)
psql $POSTGRESQL_URL < db/migrations/000001_initialize_schema.up.sql
psql $POSTGRESQL_URL < db/migrations/000002_add_display_name_to_files.up.sql

# Para revertir la última migración
migrate -path db/migrations -database "$POSTGRESQL_URL" down 1
```

**IMPORTANTE**: Si ya tienes la base de datos con el schema inicial (migración 000001), necesitas aplicar la migración 000002 para agregar el campo `display_name` a la tabla `static_files`.

## Seeder - Crear Usuario Administrador

Para crear el primer usuario administrador, usa el seeder:

```bash
# Configurar variables de entorno
export POSTGRESQL_URL="postgres://user:password@localhost:5432/jrdelperu"
export ADMIN_USERNAME="admin"
export ADMIN_EMAIL="admin@jrdelperu.com"
export ADMIN_PASSWORD="tu-password-seguro"

# Compilar y ejecutar el seeder
go build -o build/seeder ./cmd/seeder
./build/seeder
```

El seeder:
- Verifica que el usuario no exista antes de crearlo
- Hashea la contraseña con bcrypt
- Crea un usuario activo listo para usar
- Muestra información del usuario creado

Luego puedes acceder al panel de administración en `/admin/login`

## Generación de Código

### sqlc
```bash
# Generar código Go desde queries SQL
sqlc generate
```

### Templ
```bash
# Generar código Go desde templates .templ
templ generate
```

