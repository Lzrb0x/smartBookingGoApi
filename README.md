# Smart Booking Go API

MVP inicial para agendamentos de barbearia, usando Go + Gin + sqlx.

## Como rodar

```sh
cp .env.example .env # ajuste conforme necessário
docker-compose up -d
go run ./...
```

## Migrations

As migrations SQL estão em `migrations/`. Utilize sua ferramenta preferida (por exemplo, `golang-migrate`) apontando para o Postgres do `docker-compose`.
