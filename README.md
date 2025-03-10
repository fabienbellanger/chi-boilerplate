# Chi Boilerplate

A simple boilerplate for [Chi](https://go-chi.io)

[![Go Report Card](https://goreportcard.com/badge/github.com/fabienbellanger/chi-boilerplate)](https://goreportcard.com/report/github.com/fabienbellanger/chi-boilerplate)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=square)](https://pkg.go.dev/github.com/fabienbellanger/chi-boilerplate)

## Sommaire

- [Commands list](#commands-list)
- [Makefile commands](#makefile-commands)
- [Swagger](#swagger)
- [Golang web server in production](#golang-web-server-in-production)
- [Go documentation](#go-documentation)
- [Mesure et performance](#mesure-et-performance)
  - [pprof](#pprof)
  - [trace](#trace)
  - [cover](#cover)
- [TODO](#todo)

## Commands list

| Command                       | Description                 |
|-------------------------------|-----------------------------|
| `<binary> run`                | Start server                |
| `<binary> logs -s`            | Server logs reader          |
| `<binary> logs -d`            | Database (GORM) logs reader |
| `<binary> register`           | Create a new user           |
| `<binary> rabbitmq -i client` | Start RabbitMQ client       |
| `<binary> rabbitmq -i server` | Start RabbitMQ server       |

## Makefile commands

| Makefile command    | Go command                                    | Description                                 |
|---------------------|-----------------------------------------------|---------------------------------------------|
| `make update`       | `go get -u && go mod tidy`                    | Update Go dependencies                      |
| `make serve`        | `go run cmd/main.go`                          | Start the Web server                        |
| `make serve-race`   | `go run --race cmd/main.go`                   | Start the Web server with data races option |
| `make build`        | `go build -o go-url-shortener -v cmd/main.go` | Build application                           |
| `make test`         | `go test -cover ./...`                        | Launch unit tests                           |
| `make test-verbose` | `go test -cover -v ./...`                     | Launch unit tests in verbose mode           |
| `make logs`         | `go run cmd/main.go logs -s`                  | Start server logs reader                    |


## Hot reload

Install [`air`](https://github.com/air-verse/air)

Run:
```bash
air
```
or:
```bash
air | make logs
```

TODO

## Golang web server in production

- [Systemd](https://jonathanmh.com/deploying-go-apps-systemd-10-minutes-without-docker/)
- [ProxyPass](https://evanbyrne.com/blog/go-production-server-ubuntu-nginx)
- [How to Deploy App Using Docker](https://medium.com/@habibridho/docker-as-deployment-tools-5a6de294a5ff)

### Creating a Service for Systemd

```bash
touch /lib/systemd/system/<service name>.service
```

Edit file:

```
[Unit]
Description=<service description>
After=network.target

[Service]
Type=simple
Restart=always
RestartSec=5s
ExecStart=<path to exec with arguments>

[Install]
WantedBy=multi-user.target
```

| Commande                                   | Description        |
|--------------------------------------------|--------------------|
| `systemctl start <service name>.service`   | To launch          |
| `systemctl enable <service name>.service`  | To enable on boot  |
| `systemctl disable <service name>.service` | To disable on boot |
| `systemctl status <service name>.service`  | To show status     |
| `systemctl stop <service name>.service`    | To stop            |

## Database migrations

Install [golang-migrate](https://github.com/golang-migrate/migrate)

### Create a migration
```bash
migrate create -ext sql -dir migrations <migration_name>
```

### Run migrations
```bash
migrate -source file://migrations -database <connection_string> up
migrate -source file://./migrations -database "mysql://root:root@tcp(localhost:3306)/chi_boilerplate" up
```

### Revert migrations
```bash
migrate -source file://migrations -database <connection_string> down
```

## Benchmark

Use [Drill](https://github.com/fcsonline/drill)

```bash
$ drill --benchmark drill.yml --stats --quiet
```

## Go documentation

Install `godoc` (pas dans le répertoire du projet):

```bash
go get -u golang.org/x/tools/...
```

Then run:

```bash
godoc -http=localhost:6060 -play=true -index
```

## Mesure et performance

Go met à disposition de puissants outils pour mesurer les performances des programmes :

- pprof (graph, flamegraph, peek)
- trace
- cover

=> Lien vers une vidéo intéressante [Mesure et optimisation de la performance en Go](https://www.youtube.com/watch?v=jd47gDK-yDc)

Installer `graphviz`

### pprof

Lancer :

```bash
curl http://localhost:3002/debug/pprof/heap?seconds=10 > <fichier à analyser>
curl http://localhost:3002/debug/pprof/heap?seconds=10 -u "username:password" > <fichier à analyser>
```

Puis :

```bash
go tool pprof -http :3012 <fichier à analyser> # Interface web
go tool pprof --nodefraction=0 -http :3012 <fichier à analyser> # Interface web avec tous les noeuds
go tool pprof <fichier à analyser> # Ligne de commande
```

ou :

```bash
go tool pprof -http :3012 -seconds 10 http://localhost:3002/debug/pprof/heap
```

### trace

Lancer :

```bash
go test <package path> -trace=<fichier à analyser>
curl localhost:3002/debug/pprof/trace?seconds=10 > <fichier à analyser>
```

Puis :

```bash
go tool trace <fichier à analyser>
```

### cover

Lancer :

```bash
go test <package path> -covermode=count -coverprofile=./<fichier à analyser>
```

Puis :

```bash
go tool cover -html=<fichier à analyser>
```

## Generate JWT ES384 keys

```bash
mkdir keys

# Private key
openssl ecparam -name secp384r1 -genkey -noout -out keys/private.ec.key

# Public key
openssl ec -in keys/private.ec.key -pubout -out keys/public.ec.pem

# Convert SEC1 private key to PKCS8
openssl pkcs8 -topk8 -nocrypt -in keys/private.ec.key -out keys/private.ec.pem

rm keys/private.ec.key
```

## TODO

- [ ] Rework requests and responses
- [ ] Add user restore route
- [ ] Add CI with Github action
- [ ] Add / Test `Http Rate Limiting Middleware` middleware
- [ ] Add refresh token
- [ ] Add forgotten password
- [ ] Add roles and scopes
- [ ] Add Docker support
  - [ ] Try OpenTelemetry [middleware](https://github.com/gofiber/contrib/tree/main/otelfiber)
  - [ ] Mettre en place la stack Prometheus + Grafana pour la télémétrie
  - [ ] Add Prometheus metrics ([Example](https://github.com/stefanprodan/dockprom))
  - [ ] Create a first user to use API
- [ ] Try test suite [ginkgo](https://github.com/onsi/ginkgo) and [gomega](https://github.com/onsi/gomega)
- [x] Replace SHA2 by argon2/bcrypt for password hashing
- [x] Add `Timeout` middleware
- [x] Add `RealIP` middleware
- [x] Add `CORS` middleware
- [x] Add OpenAPI doc with BasicAuth support
- [x] Add JWT / Auth middleware
- [x] Change / improve logger middleware
- [x] Add unit tests for infrastructure
