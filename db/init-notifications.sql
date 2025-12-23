-- Banco de dados próprio do serviço de notificações
-- Autonomia de dados: apenas este serviço acessa este banco

CREATE TABLE IF NOT EXISTS notifications (
  id BIGSERIAL PRIMARY KEY,
  payment_id BIGINT NOT NULL,
  type TEXT NOT NULL,
  recipient TEXT NOT NULL,
  message TEXT NOT NULL,
  status TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
  -- NOTE: payment_id é apenas uma referência (sem FK, pois o pagamento está em outro serviço/banco)
);

-- NOTE: Cada microsserviço tem seu próprio banco de dados
-- Isso garante autonomia e evita acoplamento
