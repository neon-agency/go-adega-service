# Go Adega API

Backend em Go do sistema de adega. Ele expõe as rotas usadas pelo `app-adega` para catálogo, estoque, pedidos, configurações da loja, upload de imagens no GCP, motoboys, rastreamento e relatórios.

## Stack

- Go 1.24
- Echo
- PostgreSQL
- Google Cloud Storage para imagens
- SendGrid para e-mails de motoboys
- Integração com `go-payment-service`

## Como rodar

```bash
cd /home/gabriel/dev/go_adega
cp .env.example .env
go mod download
go run main.go
```

A API sobe por padrão em:

```text
http://localhost:8085
```

Swagger:

```text
http://localhost:8085/swagger
```

Health check:

```text
http://localhost:8085/health
```

## Variáveis de ambiente

Use `.env.example` como base. As principais são:

```env
SERVER_PORT=

POSTGRES_HOST=
POSTGRES_PORT=
POSTGRES_USER=
POSTGRES_PASSWORD=
POSTGRES_DB=
DB_SSL_MODE=
DB_DRIVER=

PAYMENT_SERVICE_URL=
PAYMENT_PROVIDER=

GCS_BUCKET=
GCS_PUBLIC_BASE_URL=
GOOGLE_APPLICATION_CREDENTIALS=

FRONT_APP_URL=

SENDGRID_API_KEY=
SENDGRID_FROM_EMAIL=
SENDGRID_FROM_NAME=Adega Flow
```

Não versionar `.env` nem JSON de service account.

## Migrations

As migrations ficam em:

```text
db/migration
```

O backend não executa migrations automaticamente ao iniciar. Antes de rodar uma versão nova, aplique as migrations no banco configurado no `.env`.

Exemplo usando `psql`:

```bash
set -a
source .env
set +a

PGPASSWORD="$POSTGRES_PASSWORD" psql \
  -h "$POSTGRES_HOST" \
  -p "$POSTGRES_PORT" \
  -U "$POSTGRES_USER" \
  -d "$POSTGRES_DB" \
  -v ON_ERROR_STOP=1 \
  -f db/migration/010_complete_store_settings.up.sql
```

## Principais rotas

- `GET /api/v1/settings/store`: configurações públicas e administrativas da loja
- `PUT /api/v1/settings/store`: salva configurações da loja
- `GET /api/v1/products`: catálogo e produtos do admin
- `POST /api/v1/products`: cria produto
- `PUT /api/v1/products/:id`: edita produto
- `POST /api/v1/uploads/images`: upload de imagem para o bucket GCP
- `POST /api/v1/orders`: cria pedido
- `GET /api/v1/orders`: lista pedidos do admin
- `GET /api/v1/reports`: relatórios
- `GET /api/v1/drivers`: motoboys
- `GET /api/v1/employees`: funcionários

## Testes

```bash
go test ./...
```

## Observações

- Para upload funcionar, configure `GCS_BUCKET`, `GCS_PUBLIC_BASE_URL` e `GOOGLE_APPLICATION_CREDENTIALS`.
- Para pagamentos online, o `go-payment-service` precisa estar rodando e `PAYMENT_SERVICE_URL` deve apontar para ele.
- Para envio de e-mail de motoboy, configure as variáveis do SendGrid.
