**JRDELPERU Webpage**

# Dependencies

- Golang: Main programming language of the project
- PostgreSQL: Database of the project
- sqlc: Library to handle SQL queries
- Echo: Backend framework
- Templ: HTML templating language for Go
- HTMX 2.x: Frontend library for interactive UI
- Tailwind CSS: CSS framework
- wneessen/go-mail: Library for sending mails with Go 

# Objectives

- Allow admins to dynamically upload text, images, pdfs that will be displayed
  in the view pages.
- Allow admins to do insert, update, delete operations.
- Allow the optional creation of sample data to test the webpage in
  development. (Use `config/sample.go` to get some sample data, we must delete it
  once we have a way to create the sample data directly in the database and
  for images use `assets/static/img/placeholder.webp` as the fallback.
- Implement the functionality of the simple Contact Form using the
  wneessen/go-mail library and Google Captcha. Any needed information like smtp
  server, credentials, sender email, receiver, etc. must be defined using
  environment variables.

# Important information

- The system will be used in Peru. The language must be Spanish.

# Plan

Here is a summary of how I want to build the webpage. I am open to feedback and
improvement.

- Make sure to understand every dependency used in this project. Please, search
  information if needed about these tools and their best practices.
- Explore the project and write a summary in README.md. Most of the UI design
  for the view pages are already defined with Tailwind. Keep the same style for
  other view pages. But, use another style of your choice for the admin pages.
- Documentate any needed environment variable in README.md file.
- Define a consistent data model. We already have defined a basic one, but
  we allow to change or add things because we are in development.
- Build the application progressively, not in one go. For instance, we may
  focus in searching information of the tools in this first chat, ask me if you
  need more information. In the next chat, we can start defining the model,
  build the admin panel, etc.
- Implement a simple user authentication using Go interfaces, so we can easily
  swap the implementation later.
- Follow a consistent style across all go files, sql files, etc.

# Project structure

We use the following basic template for our projects (key files)

```
.
├── Makefile: It orchestrates the server's build (it generates build/server)
├── assets/: It embeds static files into the binary via 'Assets' variable (fs.FS)
│   ├── embed.go
│   ├── embed_dev.go
│   └── static/: It contains the static files
│       ├── css/
│       └── js/
├── cmd/
│   └── server/: The server entrypoint
│       ├── init.go
│       ├── init_dev.go
│       └── main.go
├── db/: It contains migrations and queries used by sqlc to generate repository/
│   ├── migrations/
│   │   ├── 000001_initialize_schema.down.sql
│   │   └── 000001_initialize_schema.up.sql
│   └── query/
├── repository/: Generated files by sqlc
├── service/
├── model/
├── handler/
├── view/
├── sqlc.yml
├── tailwind.css
├── tailwind.config.cjs
...
```

