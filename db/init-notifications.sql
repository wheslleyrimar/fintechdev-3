-- Banco de dados próprio do serviço de notificações
-- Autonomia de dados: apenas este serviço acessa este banco

CREATE TABLE IF NOT EXISTS notifications (
  id BIGSERIAL PRIMARY KEY,
  type TEXT NOT NULL,
  recipient TEXT NOT NULL,
  message TEXT NOT NULL,
  status TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- NOTE: Cada microsserviço tem seu próprio banco de dados
-- Isso garante autonomia e evita acoplamento
