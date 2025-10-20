# JRDELPERU Web page

## Prerequisites

- Podman and Podman Compose

## .env file example

> [!IMPORTANT]
> The database is not created automatically and must be created within database
> container.
> `createdb -U postgres jrdelperu`
> The scheme is applied using
> `migrate -database ${POSTGRESQL_URL} -path db/migrations up`

```shell
# Env for the application
POSTGRESQL_URL="postgres://postgres:LlaveSecreta01@db:5432/jrdelperu?sslmode=disable"
REL="1"
SESSION_SECRET="your-secret-key-at-least-32-characters-long"

# Env for the database
POSTGRES_PASSWORD="LlaveSecreta01"

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

# Seeder
ADMIN_USERNAME="admin"
ADMIN_EMAIL="admin@jrdelperu.com"
ADMIN_PASSWORD="123456"
```

## Live reload (development)

```shell
bin/live-dev
```

## Run (production)

```shell
bin/up-prod
```

