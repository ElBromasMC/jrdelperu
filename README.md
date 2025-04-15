# go-webserver-template

Jrdelperu website

## Development environment

### Prerequisites

* Podman
* Podman compose

### .env file example

```shell
REL="1"
```

### Live reload

```shell
bin/live-dev
```

## Production environment

### Prerequisites

* Docker

### Docker compose .env file example

```shell
WEBSERVER_HOSTNAME=domain.tld
REL=1
```

### Run

```shell
bin/up-prod
```

