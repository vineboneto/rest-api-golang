# CRUD em golang com uma arquitetura minimalista e pratica

- [Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

## Executar em desenvolvimento

```bash
cp .env.example .env
# Dependencias
docker compose up -d
# Migrations
go run cmd/migrate/main.go migrate -up

# Instalar Swag Cmd
go get -u github.com/swaggo/swag/cmd/swag
# Gerar Documentação
swag init --parseDependency -g ./cmd/server/main.go -o docs

# Documentação
curl http://localhost:8080/api/v1/docs/index.html
```
