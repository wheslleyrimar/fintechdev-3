-- Banco de dados próprio do serviço de pagamentos
-- Autonomia de dados: apenas este serviço acessa este banco

CREATE TABLE IF NOT EXISTS pix_payments (
  id BIGSERIAL PRIMARY KEY,
  amount NUMERIC(18,2) NOT NULL,
  status TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- NOTE: Cada microsserviço tem seu próprio banco de dados
-- Isso garante autonomia e evita acoplamento
