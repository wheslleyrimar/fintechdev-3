package api

import (
	"encoding/json"
	app "fintech-payments-service/application"
	"fintech-payments-service/domain"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PaymentsHandler struct {
	createUC *app.CreatePixPaymentUseCase
	repo     domain.PixPaymentRepository
}

type createPixRequest struct {
	Amount float64 `json:"amount"`
}

func NewPaymentsHandler(createUC *app.CreatePixPaymentUseCase, repo domain.PixPaymentRepository) *PaymentsHandler {
	return &PaymentsHandler{createUC: createUC, repo: repo}
}

func (h *PaymentsHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/pix", h.handlePayments)
	mux.HandleFunc("/pix/", h.handlePaymentByID)
	mux.HandleFunc("/pix/monitor/", h.monitorPayment)
	mux.HandleFunc("/monitor", h.monitorPage)
}

func (h *PaymentsHandler) handlePayments(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listAll(w, r)
	case http.MethodPost:
		h.create(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *PaymentsHandler) listAll(w http.ResponseWriter, r *http.Request) {
	log.Printf("INFO: Listing all PIX payments")

	paymentsList, err := h.repo.FindAll()
	if err != nil {
		log.Printf("ERROR: Failed to list payments: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if paymentsList == nil {
		paymentsList = []*domain.PixPayment{} // Retorna array vazio ao invés de null
	}

	log.Printf("INFO: Found %d payments", len(paymentsList))
	writeJSON(w, http.StatusOK, paymentsList)
}

func (h *PaymentsHandler) create(w http.ResponseWriter, r *http.Request) {
	var req createPixRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR: Failed to decode request: %v", err)
		http.Error(w, "invalid json: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validação do valor
	if req.Amount <= 0 {
		log.Printf("ERROR: Invalid amount: %.2f", req.Amount)
		http.Error(w, "amount must be greater than 0", http.StatusBadRequest)
		return
	}

	log.Printf("INFO: Creating PIX payment with amount: %.2f", req.Amount)

	payment, err := h.createUC.Execute(req.Amount)
	if err != nil {
		log.Printf("ERROR: Failed to create payment: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("INFO: Payment created successfully - ID: %d, Amount: %.2f, Status: %s",
		payment.ID, payment.Amount, payment.Status)

	writeJSON(w, http.StatusOK, payment)
}

func (h *PaymentsHandler) handlePaymentByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Extrair ID da URL: /pix/{id}
	path := strings.TrimPrefix(r.URL.Path, "/pix/")
	if path == "" {
		http.Error(w, "payment ID is required", http.StatusBadRequest)
		return
	}

	// Ignorar rotas de monitoramento
	if strings.HasPrefix(path, "monitor/") {
		return
	}

	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		log.Printf("ERROR: Invalid payment ID: %s", path)
		http.Error(w, "invalid payment ID", http.StatusBadRequest)
		return
	}

	log.Printf("INFO: Fetching payment with ID: %d", id)

	payment, err := h.repo.FindByID(id)
	if err != nil {
		log.Printf("ERROR: Failed to find payment %d: %v", id, err)
		http.Error(w, "payment not found", http.StatusNotFound)
		return
	}

	log.Printf("INFO: Payment found - ID: %d, Amount: %.2f, Status: %s",
		payment.ID, payment.Amount, payment.Status)

	writeJSON(w, http.StatusOK, payment)
}

func (h *PaymentsHandler) monitorPayment(w http.ResponseWriter, r *http.Request) {
	// Extrair ID da URL: /pix/monitor/{id}
	path := strings.TrimPrefix(r.URL.Path, "/pix/monitor/")
	if path == "" {
		http.Error(w, "payment ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		http.Error(w, "invalid payment ID", http.StatusBadRequest)
		return
	}

	// Verificar se o pagamento existe
	payment, err := h.repo.FindByID(id)
	if err != nil {
		http.Error(w, "payment not found", http.StatusNotFound)
		return
	}

	// Configurar SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Enviar status inicial imediatamente
	initialEvent := domain.PaymentEvent{
		PaymentID: payment.ID,
		Status:    payment.Status,
		Amount:    payment.Amount,
		Timestamp: time.Now(),
		Message:   "Status inicial do pagamento",
	}
	if err := sendSSEEvent(w, "initial", initialEvent); err != nil {
		log.Printf("ERROR: Failed to send initial event: %v", err)
		return
	}
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	// Inscrever no broadcaster
	broadcaster := GetBroadcaster()
	eventChan := broadcaster.Subscribe(id)
	defer broadcaster.Unsubscribe(id, eventChan)

	// Enviar eventos em tempo real
	for event := range eventChan {
		if err := sendSSEEvent(w, "status_change", event); err != nil {
			log.Printf("ERROR: Failed to send SSE event: %v", err)
			return
		}
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}
}

func sendSSEEvent(w http.ResponseWriter, eventType string, data domain.PaymentEvent) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte("event: " + eventType + "\n"))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte("data: " + string(jsonData) + "\n\n"))
	return err
}

func (h *PaymentsHandler) monitorPage(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Monitor PIX - Tempo Real (Microsserviços)</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            border-radius: 12px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
            overflow: hidden;
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 30px;
            text-align: center;
        }
        .header h1 { font-size: 2.5em; margin-bottom: 10px; }
        .header p { opacity: 0.9; }
        .content {
            padding: 30px;
        }
        .input-group {
            margin-bottom: 30px;
            display: flex;
            gap: 10px;
        }
        input {
            flex: 1;
            padding: 15px;
            border: 2px solid #e0e0e0;
            border-radius: 8px;
            font-size: 16px;
            transition: border-color 0.3s;
        }
        input:focus {
            outline: none;
            border-color: #667eea;
        }
        button {
            padding: 15px 30px;
            background: #667eea;
            color: white;
            border: none;
            border-radius: 8px;
            font-size: 16px;
            cursor: pointer;
            transition: background 0.3s;
        }
        button:hover { background: #5568d3; }
        button:disabled {
            background: #ccc;
            cursor: not-allowed;
        }
        .status-box {
            background: #f8f9fa;
            border-radius: 8px;
            padding: 20px;
            margin-bottom: 20px;
        }
        .status-item {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 15px;
            margin: 10px 0;
            border-radius: 8px;
            background: white;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .status-badge {
            padding: 8px 16px;
            border-radius: 20px;
            font-weight: bold;
            font-size: 14px;
        }
        .status-CREATED { background: #ffc107; color: #000; }
        .status-AUTHORIZED { background: #17a2b8; color: white; }
        .status-SETTLED { background: #28a745; color: white; }
        .event-log {
            max-height: 400px;
            overflow-y: auto;
            background: #1e1e1e;
            color: #d4d4d4;
            padding: 20px;
            border-radius: 8px;
            font-family: 'Courier New', monospace;
            font-size: 14px;
        }
        .event-item {
            margin: 10px 0;
            padding: 10px;
            background: #2d2d2d;
            border-radius: 4px;
            border-left: 3px solid #667eea;
        }
        .event-time {
            color: #858585;
            font-size: 12px;
        }
        .no-connection {
            text-align: center;
            padding: 40px;
            color: #666;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔍 Monitor PIX em Tempo Real</h1>
            <p>Monitore mudanças de status de pagamentos PIX instantaneamente (Microsserviços)</p>
        </div>
        <div class="content">
            <div style="margin-bottom: 20px; padding: 20px; background: #f0f4ff; border-radius: 8px; border-left: 4px solid #667eea;">
                <h3 style="margin-bottom: 15px; color: #667eea;">🚀 Criar e Monitorar PIX</h3>
                <div class="input-group">
                    <input type="number" id="paymentAmount" placeholder="Valor do pagamento (ex: 123.45)" step="0.01" min="0.01" style="flex: 1;">
                    <button onclick="createAndMonitor()" style="background: #28a745;">Criar e Monitorar PIX</button>
                </div>
                <p style="margin-top: 10px; color: #666; font-size: 14px;">💡 Este botão cria um pagamento PIX e inicia o monitoramento automaticamente, permitindo ver todas as mudanças de status em tempo real!</p>
            </div>
            
            <div style="margin: 30px 0; text-align: center; color: #666;">
                <strong>OU</strong>
            </div>
            
            <div class="input-group">
                <input type="number" id="paymentId" placeholder="Digite o ID do pagamento existente" min="1">
                <button onclick="startMonitoring()">Monitorar Pagamento Existente</button>
                <button onclick="stopMonitoring()" id="stopBtn" disabled>Parar</button>
            </div>
            
            <div id="statusBox" class="status-box" style="display: none;">
                <h3>Status Atual</h3>
                <div class="status-item">
                    <div>
                        <strong>ID:</strong> <span id="currentId">-</span><br>
                        <strong>Valor:</strong> R$ <span id="currentAmount">-</span>
                    </div>
                    <div>
                        <span id="currentStatus" class="status-badge">-</span>
                    </div>
                </div>
            </div>

            <div id="eventLog" class="event-log" style="display: none;">
                <h3 style="color: white; margin-bottom: 15px;">📋 Log de Eventos</h3>
                <div id="events"></div>
            </div>

            <div id="noConnection" class="no-connection">
                <p>Digite um ID de pagamento e clique em "Iniciar Monitoramento"</p>
            </div>
        </div>
    </div>

    <script>
        let eventSource = null;
        let paymentId = null;

        async function createAndMonitor() {
            const amountInput = document.getElementById('paymentAmount').value;
            const amount = parseFloat(amountInput);
            
            if (isNaN(amount) || amount <= 0) {
                alert('Por favor, digite um valor válido para o pagamento (ex: 123.45)');
                return;
            }

            const createBtn = document.querySelector('button[onclick="createAndMonitor()"]');
            createBtn.disabled = true;
            createBtn.textContent = 'Criando pagamento...';

            document.getElementById('noConnection').style.display = 'none';
            document.getElementById('statusBox').style.display = 'block';
            document.getElementById('eventLog').style.display = 'block';
            document.getElementById('events').innerHTML = '';
            document.getElementById('paymentAmount').disabled = true;

            try {
                const response = await fetch('/pix', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ amount: amount })
                });

                if (!response.ok) {
                    // Tentar ler a mensagem de erro do corpo da resposta
                    let errorMessage = response.statusText;
                    try {
                        const errorData = await response.text();
                        if (errorData) {
                            errorMessage = errorData;
                        }
                    } catch (e) {
                        // Se não conseguir ler, usar statusText
                    }
                    throw new Error(errorMessage);
                }

                const payment = await response.json();
                const id = payment.id;

                document.getElementById('paymentId').value = id;
                startMonitoringWithId(id);
                addEvent('✅ Pagamento Criado', 'Pagamento PIX criado com sucesso (ID: ' + id + ', Status: CREATED)', new Date().toISOString());

            } catch (error) {
                console.error('Erro ao criar pagamento:', error);
                alert('Erro ao criar pagamento: ' + error.message);
                createBtn.disabled = false;
                createBtn.textContent = 'Criar e Monitorar PIX';
                document.getElementById('paymentAmount').disabled = false;
            }
        }

        function startMonitoring() {
            const id = document.getElementById('paymentId').value;
            if (!id) {
                alert('Por favor, digite um ID de pagamento');
                return;
            }
            startMonitoringWithId(id);
        }

        function startMonitoringWithId(id) {
            if (!id) {
                alert('ID de pagamento inválido');
                return;
            }

            paymentId = id;
            document.getElementById('paymentId').disabled = true;
            document.getElementById('stopBtn').disabled = false;
            document.querySelector('button[onclick="startMonitoring()"]').disabled = true;
            const createBtn = document.querySelector('button[onclick="createAndMonitor()"]');
            if (createBtn) createBtn.disabled = true;
            document.getElementById('noConnection').style.display = 'none';
            document.getElementById('statusBox').style.display = 'block';
            document.getElementById('eventLog').style.display = 'block';
            document.getElementById('events').innerHTML = '';

            eventSource = new EventSource('/pix/monitor/' + id);

            eventSource.addEventListener('initial', function(e) {
                try {
                    const event = JSON.parse(e.data);
                    updateStatus(event);
                    addEvent('🟢 Conectado', 'Conexão estabelecida - Status: ' + event.status, event.timestamp);
                } catch (err) {
                    console.error('Erro ao processar evento initial:', err);
                    addEvent('❌ Erro', 'Erro ao processar status inicial', new Date().toISOString());
                }
            });

            eventSource.addEventListener('status_change', function(e) {
                try {
                    const event = JSON.parse(e.data);
                    updateStatus(event);
                    addEvent('🔄 Mudança de Status', event.message + ' (Status: ' + event.status + ')', event.timestamp);
                } catch (err) {
                    console.error('Erro ao processar evento status_change:', err);
                    addEvent('❌ Erro', 'Erro ao processar mudança de status', new Date().toISOString());
                }
            });

            eventSource.onopen = function(e) {
                addEvent('🔗 Conectando', 'Estabelecendo conexão com o servidor...', new Date().toISOString());
            };

            eventSource.onerror = function(e) {
                console.error('SSE Error:', e);
                if (eventSource.readyState === EventSource.CLOSED) {
                    addEvent('❌ Erro', 'Conexão perdida ou pagamento não encontrado', new Date().toISOString());
                } else {
                    addEvent('⚠️ Aviso', 'Erro na conexão, tentando reconectar...', new Date().toISOString());
                }
            };
        }

        function stopMonitoring() {
            if (eventSource) {
                eventSource.close();
                eventSource = null;
            }
            document.getElementById('paymentId').disabled = false;
            document.getElementById('paymentAmount').disabled = false;
            document.getElementById('stopBtn').disabled = true;
            document.querySelector('button[onclick="startMonitoring()"]').disabled = false;
            const createBtn = document.querySelector('button[onclick="createAndMonitor()"]');
            if (createBtn) {
                createBtn.disabled = false;
                createBtn.textContent = 'Criar e Monitorar PIX';
            }
            addEvent('⏹️ Parado', 'Monitoramento interrompido', new Date().toISOString());
        }

        function updateStatus(event) {
            document.getElementById('currentId').textContent = event.payment_id;
            document.getElementById('currentAmount').textContent = event.amount.toFixed(2);
            const statusEl = document.getElementById('currentStatus');
            statusEl.textContent = event.status;
            statusEl.className = 'status-badge status-' + event.status;
        }

        function addEvent(type, message, timestamp) {
            const eventsDiv = document.getElementById('events');
            const eventDiv = document.createElement('div');
            eventDiv.className = 'event-item';
            
            const date = new Date(timestamp);
            const time = date.toLocaleTimeString('pt-BR') + '.' + date.getMilliseconds().toString().padStart(3, '0');
            eventDiv.innerHTML = 
                '<div style="color: #667eea; font-weight: bold;">' + type + '</div>' +
                '<div style="margin: 5px 0;">' + message + '</div>' +
                '<div class="event-time">' + time + '</div>';
            
            eventsDiv.insertBefore(eventDiv, eventsDiv.firstChild);
        }

        document.getElementById('paymentId').addEventListener('keypress', function(e) {
            if (e.key === 'Enter') {
                startMonitoring();
            }
        });
    </script>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(html))
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
