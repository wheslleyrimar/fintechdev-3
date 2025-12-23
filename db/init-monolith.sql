-- Banco de dados compartilhado no monólito
-- Todos os domínios (payments, notifications) usam o mesmo banco

CREATE TABLE IF NOT EXISTS pix_payments (
  id BIGSERIAL PRIMARY KEY,
  amount NUMERIC(18,2) NOT NULL,
  status TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS notifications (
  id BIGSERIAL PRIMARY KEY,
  type TEXT NOT NULL,
  recipient TEXT NOT NULL,
  message TEXT NOT NULL,
  status TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- NOTE: No monólito, ambos os domínios compartilham o mesmo banco
-- Isso quebra a autonomia de dados, mas é aceitável em um monólito inicial
