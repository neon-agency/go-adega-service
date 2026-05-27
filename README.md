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

Em Cloud Run, a aplicação usa automaticamente a variável `PORT` injetada pela plataforma, normalmente `8080`. Não configure `SERVER_PORT` no Cloud Run a menos que você também altere a porta do container no serviço.

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
SERVER_PORT=8085

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

O backend executa as migrations automaticamente ao iniciar, usando `file://db/migration`. Por isso, ao rodar via Docker, a imagem precisa incluir a pasta `db/migration`.

Se precisar aplicar uma migration manualmente para debug, use `psql`:

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

## Docker

Build da imagem:

```bash
docker build -t go-adega-service .
```

Execução local usando variáveis do `.env`:

```bash
docker run --rm --env-file .env -p 8085:8085 go-adega-service
```

Simulando Cloud Run localmente:

```bash
docker run --rm --env-file .env -e PORT=8080 -p 8080:8080 go-adega-service
```

Se usar upload para GCP localmente, monte o JSON da service account e aponte `GOOGLE_APPLICATION_CREDENTIALS` para o caminho dentro do container:

```bash
docker run --rm --env-file .env \
  -e GOOGLE_APPLICATION_CREDENTIALS=/secrets/gcp-service-account.json \
  -v "$PWD/gcp-service-account.json:/secrets/gcp-service-account.json:ro" \
  -p 8085:8085 \
  go-adega-service
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

- Para upload funcionar localmente, configure `GCS_BUCKET`, `GCS_PUBLIC_BASE_URL` e `GOOGLE_APPLICATION_CREDENTIALS`. No GCP, prefira a service account do próprio serviço.
- Para pagamentos online, o `go-payment-service` precisa estar rodando e `PAYMENT_SERVICE_URL` deve apontar para ele.
- Para envio de e-mail de motoboy, configure as variáveis do SendGrid.
