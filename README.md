# Fechamento automático de leilões

Sistema de leilões em Go. Ao criar um leilão, uma goroutine aguarda o tempo configurado e altera o status dele para fechado no MongoDB.

## Como rodar

É preciso ter Docker e Docker Compose.

```bash
docker compose up --build
```

A API fica em `http://localhost:8080`.

## Variáveis de ambiente

As variáveis ficam em `cmd/auction/.env`. A duração do leilão é definida por:

```env
AUCTION_DURATION=20s
```

O valor usa o formato de `time.ParseDuration`, por exemplo `30s`, `5m` ou `1h`. Se a variável estiver vazia ou inválida, a duração será de 5 minutos.

Depois de alterar o valor, reinicie a aplicação.

## Testando manualmente

Crie um leilão:

```bash
curl -X POST http://localhost:8080/auction \
  -H "Content-Type: application/json" \
  -d '{"product_name":"Notebook","category":"Eletronicos","description":"Notebook usado em bom estado","condition":1}'
```

Liste os leilões para pegar o `id`:

```bash
curl "http://localhost:8080/auction?status=0"
```

Consulte o leilão:

```bash
curl http://localhost:8080/auction/{id}
```

O campo `status` começa em `0` (aberto) e passa para `1` (fechado) quando o tempo de `AUCTION_DURATION` termina.

## Testes

O teste usa o MongoDB do Docker Compose, então suba o banco antes:

```bash
docker compose up -d mongodb
go test -count=1 ./internal/infra/database/auction/...
```

O teste cria um leilão com duração de 2 segundos, confere que ele está aberto, aguarda 3 segundos e confere que ele foi fechado. Ele usa o banco `auctions_test`, que é apagado no final.

Também é possível rodar sem ter Go instalado:

```bash
docker compose run --rm --build --entrypoint "go test -count=1 ./internal/infra/database/auction/..." app
```
