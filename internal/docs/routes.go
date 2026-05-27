package docs

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	e.GET("/openapi.json", openapi)
	e.GET("/swagger", swagger)
}

func swagger(c echo.Context) error {
	html := `<!doctype html>
<html lang="pt-BR">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Adega API - Swagger</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
      window.ui = SwaggerUIBundle({
        url: '/openapi.json',
        dom_id: '#swagger-ui',
        deepLinking: true,
        persistAuthorization: true
      });
    </script>
  </body>
</html>`
	return c.HTML(http.StatusOK, html)
}

func openapi(c echo.Context) error {
	return c.JSON(http.StatusOK, spec())
}

func spec() map[string]any {
	return map[string]any{
		"openapi": "3.0.3",
		"info": map[string]any{
			"title":       "Adega API",
			"version":     "1.0.0",
			"description": "API do sistema de adega. As descrições das rotas indicam onde cada endpoint é usado no front app-adega.",
		},
		"servers": []map[string]string{
			{"url": "http://localhost:8085", "description": "Backend local"},
		},
		"tags": []map[string]string{
			{"name": "Health", "description": "Saúde da API"},
			{"name": "Auth", "description": "Autenticação do painel administrativo"},
			{"name": "Products", "description": "Catálogo e estoque de produtos"},
			{"name": "Orders", "description": "Pedidos do cliente e kanban administrativo"},
			{"name": "Settings", "description": "Configurações públicas da loja"},
			{"name": "Uploads", "description": "Upload de imagens para bucket GCP"},
			{"name": "Reports", "description": "Relatórios e gráficos administrativos"},
			{"name": "People", "description": "Motoboys e funcionários"},
			{"name": "Tracking", "description": "Entregas e rastreamento"},
			{"name": "Metrics", "description": "Indicadores operacionais"},
		},
		"paths":      paths(),
		"components": components(),
	}
}

func paths() map[string]any {
	return map[string]any{
		"/health": map[string]any{
			"get": op("Health", "Health check", "Verifica se a API está no ar. Uso no app-adega: não é chamado por tela; útil para debug, deploy e monitoramento.", nil, "HealthResponse", false),
		},
		"/api/v1/auth/register": map[string]any{
			"post": op("Auth", "Cadastrar usuário administrativo", "Cria a conta administrativa inicial da loja. Uso no app-adega: tela LoginPage em modo 'Criar conta'.", "RegisterRequest", "AuthSession", false),
		},
		"/api/v1/auth/login": map[string]any{
			"post": op("Auth", "Login administrativo", "Autentica o dono/administrador. Uso no app-adega: LoginPage ao clicar em Entrar.", "LoginRequest", "AuthSession", false),
		},
		"/api/v1/auth/me": map[string]any{
			"get": op("Auth", "Sessão atual", "Retorna o usuário logado a partir do token. Uso no app-adega: AuthContext/authController.restoreSession ao restaurar sessão.", nil, "User", true),
		},
		"/api/v1/auth/logout": map[string]any{
			"post": op("Auth", "Logout", "Finaliza a sessão no cliente. Uso no app-adega: botão sair no AdminLayout.", nil, nil, true),
		},
		"/api/v1/settings/store": map[string]any{
			"get": op("Settings", "Buscar configurações da loja", "Retorna nome, contato, endereço, taxa, regras comerciais e métodos de pagamento. Uso no app-adega: productsRepository.getCatalog monta o catálogo e CheckoutScreen; SettingsPage carrega a aba Pagamentos.", nil, "StoreSettings", false),
			"put": op("Settings", "Atualizar configurações da loja", "Atualiza dados públicos, endereço, horários, mídia, regras de entrega e métodos de pagamento online/na entrega. Uso no app-adega: SettingsPage aba Pagamentos salva o que aparece no checkout.", "UpdateStoreSettingsRequest", "StoreSettings", true),
		},
		"/api/v1/products": map[string]any{
			"get": withParams(op("Products", "Listar produtos", "Lista produtos ativos para cliente ou todos para admin quando admin=true. Uso no app-adega: MenuScreen/ProductScreen via catálogo, ProductsPage e InventoryPage.", nil, "ProductList", false), []map[string]any{
				query("category", "string", "Filtra por categoria."),
				query("admin", "boolean", "Quando true, inclui produtos inativos para o painel."),
				query("mode", "string", "Parâmetro enviado pelo front para modo de negócio; mantido para compatibilidade."),
			}),
			"post": op("Products", "Criar produto", "Cadastra um produto. Uso no app-adega: ProductsPage, drawer de novo produto.", "CreateProductRequest", "Product", true),
		},
		"/api/v1/products/{id}": map[string]any{
			"get":    withParams(op("Products", "Detalhar produto", "Busca um produto por ID. Uso no app-adega: não é chamado hoje; ProductScreen usa o produto já carregado no catálogo.", nil, "Product", false), []map[string]any{pathParam("id", "string", "ID do produto.")}),
			"put":    withParams(op("Products", "Atualizar produto", "Atualiza dados, preço, custo, estoque mínimo, imagem e disponibilidade. Uso no app-adega: ProductsPage, drawer de edição.", "UpdateProductRequest", "Product", true), []map[string]any{pathParam("id", "string", "ID do produto.")}),
			"delete": withParams(op("Products", "Excluir produto", "Remove um produto. Uso no app-adega: ProductsPage, ação de excluir no drawer.", nil, nil, true), []map[string]any{pathParam("id", "string", "ID do produto.")}),
		},
		"/api/v1/products/{id}/availability": map[string]any{
			"patch": withParams(op("Products", "Ativar ou desativar produto", "Altera se o produto aparece para venda. Uso no app-adega: ProductsPage, Toggle de disponibilidade.", "AvailabilityRequest", nil, true), []map[string]any{pathParam("id", "string", "ID do produto.")}),
		},
		"/api/v1/products/{id}/stock-movements": map[string]any{
			"post": withParams(op("Products", "Movimentar estoque", "Registra entrada, venda, ajuste ou perda. Uso no app-adega: InventoryPage ao ajustar estoque.", "StockMovementRequest", nil, true), []map[string]any{pathParam("id", "string", "ID do produto.")}),
		},
		"/api/v1/uploads/images": map[string]any{
			"post": map[string]any{
				"tags":        []string{"Uploads"},
				"summary":     "Enviar imagem",
				"description": "Envia imagem para o bucket GCP configurado e retorna a URL pública. Uso no app-adega: ProductsPage, área de imagem do drawer de produto.",
				"requestBody": map[string]any{
					"required": true,
					"content": map[string]any{
						"multipart/form-data": map[string]any{
							"schema": map[string]any{
								"type": "object",
								"properties": map[string]any{
									"file":   map[string]any{"type": "string", "format": "binary"},
									"folder": map[string]any{"type": "string", "example": "products"},
								},
								"required": []string{"file"},
							},
						},
					},
				},
				"responses": responses("UploadResponse"),
				"security":  bearer(),
			},
		},
		"/api/v1/orders": map[string]any{
			"get": withParams(op("Orders", "Listar pedidos", "Lista pedidos com cliente, endereço e itens. Uso no app-adega: OrdersKanbanPage, colunas do kanban administrativo.", nil, "OrderList", true), []map[string]any{
				query("status", "string", "Filtra por status do backend, como paid, separating ou delivered."),
			}),
			"post": op("Orders", "Criar pedido", "Cria pedido do cliente, reserva estoque, valida payment_mode/payment_method contra SettingsPage e só aciona o serviço de pagamento para pagamento online. Uso no app-adega: MobilePage/CheckoutScreen ao confirmar pedido.", "CreateOrderRequest", "Order", false),
		},
		"/api/v1/orders/{id}": map[string]any{
			"get":    withParams(op("Orders", "Detalhar pedido", "Busca pedido completo por ID. Uso no app-adega: não é chamado pela tela atual; kanban usa a listagem.", nil, "Order", true), []map[string]any{pathParam("id", "string", "ID do pedido.")}),
			"delete": withParams(op("Orders", "Cancelar pedido", "Marca pedido como canceled. Uso no app-adega: OrdersKanbanPage, ação Cancelar no card.", nil, nil, true), []map[string]any{pathParam("id", "string", "ID do pedido.")}),
		},
		"/api/v1/orders/{id}/status": map[string]any{
			"patch": withParams(op("Orders", "Atualizar status do pedido", "Avança ou move pedido entre status. Uso no app-adega: OrdersKanbanPage, botão Avançar e drag-and-drop entre colunas.", "UpdateStatusRequest", "Order", true), []map[string]any{pathParam("id", "string", "ID do pedido.")}),
		},
		"/api/v1/reports": map[string]any{
			"get": withParams(op("Reports", "Relatório de performance", "Retorna séries, pagamentos, top produtos e categorias. Uso no app-adega: ReportsPage, KPIs e gráficos.", nil, "ReportData", true), []map[string]any{
				query("mode", "string", "Modo enviado pelo front: adega ou comida."),
				query("period", "string", "Período: today, 7d, 30d ou 90d."),
			}),
		},
		"/api/v1/metrics/overview": map[string]any{
			"get": op("Metrics", "Resumo operacional", "Retorna entradas, saídas, lucro bruto, lucro líquido e pedidos. Uso no app-adega: não é chamado pelo front atual; rota mantida para painéis/resumos operacionais.", nil, "Overview", true),
		},
		"/api/v1/drivers": map[string]any{
			"get":  op("People", "Listar motoboys", "Lista motoboys cadastrados. Uso no app-adega: não é chamado pelo front atual; rota preparada para tela de equipe/motoboys.", nil, "PersonList", true),
			"post": op("People", "Cadastrar motoboy", "Cadastra motoboy, senha provisória e envio por SendGrid quando configurado. Uso no app-adega: não é chamado pelo front atual; rota preparada para gestão de motoboys.", "CreatePersonRequest", "Person", true),
		},
		"/api/v1/drivers/login": map[string]any{
			"post": op("People", "Login do motoboy", "Autentica motoboy por e-mail e senha. Uso no app-adega: não é chamado pelo front atual; rota prevista para app/tela do entregador.", "DriverLoginRequest", "DriverLoginResponse", false),
		},
		"/api/v1/drivers/{id}": map[string]any{
			"put": withParams(op("People", "Atualizar motoboy", "Atualiza cadastro e limite de entregas do motoboy. Uso no app-adega: não é chamado pelo front atual.", "UpdatePersonRequest", "Person", true), []map[string]any{pathParam("id", "string", "ID do motoboy.")}),
		},
		"/api/v1/drivers/{id}/password": map[string]any{
			"patch": withParams(op("People", "Trocar senha do motoboy", "Troca a senha provisória no primeiro acesso. Uso no app-adega: não é chamado pelo front atual; rota prevista para tela do entregador.", "DriverChangePasswordRequest", nil, false), []map[string]any{pathParam("id", "string", "ID do motoboy.")}),
		},
		"/api/v1/employees": map[string]any{
			"get":  op("People", "Listar funcionários", "Lista funcionários administrativos/operacionais. Uso no app-adega: não é chamado pelo front atual; SettingsPage ainda usa dados locais na aba Equipe.", nil, "PersonList", true),
			"post": op("People", "Cadastrar funcionário", "Cria funcionário. Uso no app-adega: não é chamado pelo front atual; rota preparada para aba Equipe.", "CreatePersonRequest", "Person", true),
		},
		"/api/v1/employees/{id}": map[string]any{
			"put": withParams(op("People", "Atualizar funcionário", "Atualiza funcionário. Uso no app-adega: não é chamado pelo front atual; rota preparada para aba Equipe.", "UpdatePersonRequest", "Person", true), []map[string]any{pathParam("id", "string", "ID do funcionário.")}),
		},
		"/api/v1/orders/{orderId}/delivery": map[string]any{
			"post": withParams(op("Tracking", "Criar entrega do pedido", "Cria código de rastreio e opcionalmente atribui motoboy. Uso no app-adega: não é chamado pelo front atual; rota prevista para fluxo de entrega do admin.", "CreateDeliveryRequest", "TrackingInfo", true), []map[string]any{pathParam("orderId", "string", "ID do pedido.")}),
		},
		"/api/v1/deliveries/available": map[string]any{
			"get": op("Tracking", "Entregas disponíveis", "Lista entregas que o motoboy pode pegar. Uso no app-adega: não é chamado pelo front atual; rota prevista para tela do entregador.", nil, "TrackingList", false),
		},
		"/api/v1/drivers/{driverId}/deliveries": map[string]any{
			"get": withParams(op("Tracking", "Entregas do motoboy", "Lista entregas atribuídas a um motoboy. Uso no app-adega: não é chamado pelo front atual; rota prevista para tela do entregador.", nil, "TrackingList", false), []map[string]any{pathParam("driverId", "string", "ID do motoboy.")}),
		},
		"/api/v1/drivers/{driverId}/deliveries/{code}/claim": map[string]any{
			"post": withParams(op("Tracking", "Motoboy pegar entrega", "Atribui entrega ao motoboy validando limite e localização no estabelecimento. Uso no app-adega: não é chamado pelo front atual; rota prevista para tela do entregador.", "ClaimDeliveryRequest", nil, false), []map[string]any{pathParam("driverId", "string", "ID do motoboy."), pathParam("code", "string", "Código de rastreio.")}),
		},
		"/api/v1/drivers/{driverId}/deliveries/{code}/location": map[string]any{
			"patch": withParams(op("Tracking", "Enviar localização do motoboy", "Atualiza posição GPS da entrega validando motoboy. Uso no app-adega: não é chamado pelo front atual; rota prevista para tela do entregador.", "LocationUpdateRequest", nil, false), []map[string]any{pathParam("driverId", "string", "ID do motoboy."), pathParam("code", "string", "Código de rastreio.")}),
		},
		"/api/v1/tracking/{code}": map[string]any{
			"get": withParams(op("Tracking", "Rastrear pedido", "Retorna status, cliente, endereço, motoboy e última localização. Uso no app-adega: não é chamado pelo front atual; rota prevista para página pública de rastreio.", nil, "TrackingInfo", false), []map[string]any{pathParam("code", "string", "Código de rastreio.")}),
		},
		"/api/v1/tracking/{code}/location": map[string]any{
			"patch": withParams(op("Tracking", "Atualizar localização por código", "Atualiza status/localização usando o código público. Uso no app-adega: não é chamado pelo front atual; preferir rota validada por motoboy.", "LocationUpdateRequest", nil, false), []map[string]any{pathParam("code", "string", "Código de rastreio.")}),
		},
	}
}

func op(tag, summary, description string, requestSchema any, responseSchema any, secured bool) map[string]any {
	out := map[string]any{
		"tags":        []string{tag},
		"summary":     summary,
		"description": description,
		"responses":   responses(responseSchema),
	}
	if requestSchema != nil {
		out["requestBody"] = jsonBody(requestSchema)
	}
	if secured {
		out["security"] = bearer()
	}
	return out
}

func withParams(operation map[string]any, params []map[string]any) map[string]any {
	operation["parameters"] = params
	return operation
}

func jsonBody(schema any) map[string]any {
	return map[string]any{
		"required": true,
		"content": map[string]any{
			"application/json": map[string]any{
				"schema": schemaRef(schema),
			},
		},
	}
}

func responses(schema any) map[string]any {
	if schema == nil {
		return map[string]any{"204": map[string]string{"description": "Sem conteúdo"}}
	}
	return map[string]any{
		"200": map[string]any{
			"description": "Sucesso",
			"content": map[string]any{
				"application/json": map[string]any{"schema": schemaRef(schema)},
			},
		},
		"400": errorResponse("Requisição inválida"),
		"401": errorResponse("Não autorizado"),
	}
}

func errorResponse(description string) map[string]any {
	return map[string]any{
		"description": description,
		"content": map[string]any{
			"application/json": map[string]any{"schema": ref("ErrorResponse")},
		},
	}
}

func schemaRef(schema any) map[string]any {
	if name, ok := schema.(string); ok {
		return ref(name)
	}
	return schema.(map[string]any)
}

func ref(name string) map[string]any {
	return map[string]any{"$ref": "#/components/schemas/" + name}
}

func bearer() []map[string][]string {
	return []map[string][]string{{"bearerAuth": {}}}
}

func pathParam(name, typ, description string) map[string]any {
	return map[string]any{"name": name, "in": "path", "required": true, "description": description, "schema": map[string]string{"type": typ}}
}

func query(name, typ, description string) map[string]any {
	return map[string]any{"name": name, "in": "query", "required": false, "description": description, "schema": map[string]string{"type": typ}}
}

func components() map[string]any {
	return map[string]any{
		"securitySchemes": map[string]any{
			"bearerAuth": map[string]string{"type": "http", "scheme": "bearer"},
		},
		"schemas": schemas(),
	}
}

func schemas() map[string]any {
	return map[string]any{
		"ErrorResponse":               obj(map[string]any{"error": str("Mensagem de erro"), "message": str("Mensagem alternativa")}, nil),
		"HealthResponse":              obj(map[string]any{"status": str("ok"), "service": str("go_adega"), "timestamp": str("2026-05-24T10:00:00-03:00")}, nil),
		"LoginRequest":                obj(map[string]any{"email": str("admin@adega.com"), "password": str("admin123")}, []string{"email", "password"}),
		"RegisterRequest":             obj(map[string]any{"name": str("Gabriel Melo"), "email": str("gabriel@adega.com"), "phone": str("(11) 99999-9999"), "storeName": str("Adega Flow"), "password": str("senha123")}, []string{"name", "email", "phone", "storeName", "password"}),
		"AuthSession":                 obj(map[string]any{"token": str("admin:uuid"), "user": ref("User")}, []string{"token", "user"}),
		"User":                        obj(map[string]any{"id": str("uuid"), "name": str("Administrador"), "email": str("admin@adega.com"), "role": str("owner"), "storeName": str("Adega Flow"), "initials": str("AD"), "color": str("#6E1F2C")}, []string{"id", "name", "email", "role"}),
		"StoreSettings":               obj(map[string]any{"id": str("uuid"), "store_name": str("Adega Flow"), "phone": str("5511931153811"), "whatsapp": str("5511931153811"), "logo_url": str("https://storage.googleapis.com/bucket/logo.png"), "banner_url": str("https://storage.googleapis.com/bucket/banner.png"), "address_street": str("Rua da Adega"), "address_number": str("100"), "address_neighborhood": str("Centro"), "address_city": str("Sao Paulo"), "address_state": str("SP"), "address_zip_code": str("01000-000"), "delivery_fee_cents": integer(790), "free_delivery_from_cents": integer(12000), "min_order_cents": integer(0), "latitude": number(-23.55), "longitude": number(-46.63), "driver_pickup_radius_meters": integer(150), "opening_time": str("09:00:00"), "closing_time": str("22:00:00"), "is_open": boolSchema(true), "accept_online_pix": boolSchema(true), "accept_online_card": boolSchema(true), "accept_delivery_pix": boolSchema(false), "accept_delivery_card": boolSchema(true), "accept_delivery_cash": boolSchema(true)}, nil),
		"UpdateStoreSettingsRequest":  obj(map[string]any{"store_name": str("Adega Flow"), "phone": str("5511931153811"), "whatsapp": str("5511931153811"), "logo_url": str("https://storage.googleapis.com/bucket/logo.png"), "banner_url": str("https://storage.googleapis.com/bucket/banner.png"), "address_street": str("Rua da Adega"), "address_number": str("100"), "address_complement": str(""), "address_neighborhood": str("Centro"), "address_city": str("Sao Paulo"), "address_state": str("SP"), "address_zip_code": str("01000-000"), "delivery_fee_cents": integer(790), "free_delivery_from_cents": integer(12000), "min_order_cents": integer(0), "latitude": number(-23.55), "longitude": number(-46.63), "driver_pickup_radius_meters": integer(150), "opening_time": str("09:00"), "closing_time": str("22:00"), "is_open": boolSchema(true), "accept_online_pix": boolSchema(true), "accept_online_card": boolSchema(true), "accept_delivery_pix": boolSchema(false), "accept_delivery_card": boolSchema(true), "accept_delivery_cash": boolSchema(true)}, nil),
		"Product":                     obj(map[string]any{"id": str("uuid"), "sku": str(""), "name": str("Jack Daniels"), "description": str("Whiskey Tennessee Apple"), "category": str("destilados"), "image_url": str("https://storage.googleapis.com/bucket/product.png"), "price_cents": integer(9999), "cost_cents": integer(5999), "stock_quantity": integer(40), "min_stock_quantity": integer(10), "is_active": boolSchema(true), "created_at": str("2026-05-24T10:00:00Z"), "updated_at": str("2026-05-24T10:00:00Z")}, nil),
		"ProductList":                 arr("Product"),
		"CreateProductRequest":        obj(map[string]any{"sku": str(""), "name": str("Jack Daniels"), "description": str("Whiskey Tennessee Apple"), "category": str("destilados"), "image_url": str("https://storage.googleapis.com/bucket/product.png"), "price_cents": integer(9999), "cost_cents": integer(5999), "stock_quantity": integer(40), "min_stock_quantity": integer(10)}, []string{"name", "description", "category", "price_cents"}),
		"UpdateProductRequest":        obj(map[string]any{"sku": str(""), "name": str("Jack Daniels"), "description": str("Whiskey Tennessee Apple"), "category": str("destilados"), "image_url": str("https://storage.googleapis.com/bucket/product.png"), "price_cents": integer(9999), "cost_cents": integer(5999), "stock_quantity": integer(40), "min_stock_quantity": integer(10), "is_active": boolSchema(true)}, []string{"name", "description", "category", "price_cents", "is_active"}),
		"AvailabilityRequest":         obj(map[string]any{"on": boolSchema(true)}, []string{"on"}),
		"StockMovementRequest":        obj(map[string]any{"type": enum("entry", "sale", "adjustment", "loss"), "quantity": integer(12), "unit_cost_cents": integer(5000), "total_cost_cents": integer(60000), "notes": str("Entrada pelo painel")}, []string{"type", "quantity"}),
		"UploadResponse":              obj(map[string]any{"url": str("https://storage.googleapis.com/adega-produtos/products/file.png")}, []string{"url"}),
		"CustomerRequest":             obj(map[string]any{"name": str("Marina Cardoso"), "phone": str("(11) 98421-3340"), "email": str("cliente@email.com"), "document": str("00000000000")}, []string{"name", "phone"}),
		"AddressRequest":              obj(map[string]any{"street": str("R. Aspicuelta"), "number": str("420"), "complement": str("Apto 71"), "neighborhood": str("Vila Madalena"), "city": str("São Paulo"), "state": str("SP"), "zip_code": str("05432-000"), "latitude": number(-23.55), "longitude": number(-46.63)}, []string{"street", "number", "neighborhood", "city", "state", "zip_code"}),
		"OrderItemRequest":            obj(map[string]any{"product_id": str("uuid"), "quantity": integer(2)}, []string{"product_id", "quantity"}),
		"CreateOrderRequest":          obj(map[string]any{"customer": ref("CustomerRequest"), "address": ref("AddressRequest"), "payment_method": enum("pix", "credit_card", "debit_card", "cash"), "payment_mode": enum("online", "delivery"), "payment_token": str("token"), "provider": str("efi"), "installments": integer(1), "notes": str("Sem gelo"), "items": arr("OrderItemRequest")}, []string{"customer", "address", "items"}),
		"OrderItem":                   obj(map[string]any{"id": str("uuid"), "product_id": str("uuid"), "product_name": str("Jack Daniels"), "quantity": integer(2), "unit_price_cents": integer(9999), "total_cents": integer(19998)}, nil),
		"PaymentInfo":                 obj(map[string]any{"provider": str("efi"), "reference": str("provider-id"), "status": str("pending"), "qr_code": str("base64"), "copy_paste": str("pix copia e cola"), "payment_url": str("https://pagamento"), "payment_method": str("pix"), "payment_mode": str("online")}, nil),
		"Order":                       obj(map[string]any{"id": str("uuid"), "customer_id": str("uuid"), "customer_name": str("Marina Cardoso"), "customer_phone": str("(11) 98421-3340"), "delivery_address": str("R. Aspicuelta, 420, Vila Madalena, São Paulo"), "status": enum("created", "awaiting_payment", "paid", "separating", "out_for_delivery", "delivered", "canceled"), "subtotal_cents": integer(19998), "delivery_fee_cents": integer(790), "total_cents": integer(20788), "notes": str("Sem gelo"), "payment": ref("PaymentInfo"), "items": arr("OrderItem"), "created_at": str("2026-05-24T10:00:00Z"), "updated_at": str("2026-05-24T10:00:00Z")}, nil),
		"OrderList":                   arr("Order"),
		"UpdateStatusRequest":         obj(map[string]any{"status": enum("created", "awaiting_payment", "paid", "separating", "out_for_delivery", "delivered", "canceled")}, []string{"status"}),
		"ReportData":                  obj(map[string]any{"series": arr("SeriesPoint"), "payments": arr("PaymentSlice"), "top": arr("TopProduct"), "byCategory": arr("CategoryReport")}, nil),
		"SeriesPoint":                 obj(map[string]any{"label": str("24/05"), "value": number(1200.50), "orders": integer(12), "target": number(1500)}, nil),
		"PaymentSlice":                obj(map[string]any{"label": str("PIX"), "pct": integer(54), "color": str("var(--wine)")}, nil),
		"TopProduct":                  obj(map[string]any{"id": str("uuid"), "name": str("Jack Daniels"), "swatch": str("https://storage.googleapis.com/bucket/product.png"), "price": number(99.99), "sold": integer(20), "revenue": number(1999.80)}, nil),
		"CategoryReport":              obj(map[string]any{"id": str("destilados"), "label": str("destilados"), "sold": integer(20), "revenue": number(1999.80)}, nil),
		"Overview":                    obj(map[string]any{"entries_cents": integer(100000), "sales_cents": integer(250000), "cost_of_goods_cents": integer(120000), "gross_profit_cents": integer(130000), "net_profit_cents": integer(130000), "orders_count": integer(25)}, nil),
		"Person":                      obj(map[string]any{"id": str("uuid"), "name": str("João Silva"), "email": str("joao@adega.com"), "phone": str("(11) 99999-9999"), "role": str("driver"), "is_active": boolSchema(true), "must_change_password": boolSchema(true), "max_active_deliveries": integer(3)}, nil),
		"PersonList":                  arr("Person"),
		"CreatePersonRequest":         obj(map[string]any{"name": str("João Silva"), "email": str("joao@adega.com"), "phone": str("(11) 99999-9999"), "role": str("driver"), "max_active_deliveries": integer(3)}, []string{"name", "phone"}),
		"UpdatePersonRequest":         obj(map[string]any{"name": str("João Silva"), "email": str("joao@adega.com"), "phone": str("(11) 99999-9999"), "role": str("driver"), "is_active": boolSchema(true), "max_active_deliveries": integer(3)}, nil),
		"DriverLoginRequest":          obj(map[string]any{"email": str("joao@adega.com"), "password": str("senha-provisoria")}, []string{"email", "password"}),
		"DriverLoginResponse":         obj(map[string]any{"driver": ref("Person"), "must_change_password": boolSchema(true)}, nil),
		"DriverChangePasswordRequest": obj(map[string]any{"current_password": str("senha-provisoria"), "new_password": str("nova-senha")}, []string{"current_password", "new_password"}),
		"CreateDeliveryRequest":       obj(map[string]any{"driver_id": str("uuid")}, nil),
		"ClaimDeliveryRequest":        obj(map[string]any{"latitude": number(-23.55), "longitude": number(-46.63)}, []string{"latitude", "longitude"}),
		"LocationUpdateRequest":       obj(map[string]any{"status": str("out_for_delivery"), "latitude": number(-23.55), "longitude": number(-46.63), "estimated_arrival_at": str("2026-05-24T10:40:00Z")}, nil),
		"TrackingInfo":                obj(map[string]any{"id": str("uuid"), "order_id": str("uuid"), "driver_id": str("uuid"), "tracking_code": str("ABC12345"), "status": str("out_for_delivery"), "current_latitude": number(-23.55), "current_longitude": number(-46.63), "customer_name": str("Marina Cardoso"), "customer_phone": str("(11) 98421-3340"), "delivery_address": str("R. Aspicuelta, 420"), "driver_name": str("João Silva")}, nil),
		"TrackingList":                arr("TrackingInfo"),
	}
}

func obj(properties map[string]any, required []string) map[string]any {
	out := map[string]any{"type": "object", "properties": properties}
	if required != nil {
		out["required"] = required
	}
	return out
}

func arr(itemSchema string) map[string]any {
	return map[string]any{"type": "array", "items": ref(itemSchema)}
}

func str(example string) map[string]any {
	return map[string]any{"type": "string", "example": example}
}

func integer(example int) map[string]any {
	return map[string]any{"type": "integer", "example": example}
}

func number(example float64) map[string]any {
	return map[string]any{"type": "number", "example": example}
}

func boolSchema(example bool) map[string]any {
	return map[string]any{"type": "boolean", "example": example}
}

func enum(values ...string) map[string]any {
	return map[string]any{"type": "string", "enum": values, "example": values[0]}
}
