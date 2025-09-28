CREATE TABLE IF NOT EXISTS accounts (
  id BIGSERIAL PRIMARY KEY,
  document_number TEXT NOT NULL UNIQUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS operation_types (
  id SMALLINT PRIMARY KEY,
  description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions (
  id BIGSERIAL PRIMARY KEY,
  account_id BIGINT NOT NULL REFERENCES accounts(id),
  operation_type_id SMALLINT NOT NULL REFERENCES operation_types(id),
  amount NUMERIC(18,2) NOT NULL,
  event_date TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT transactions_amount_sign CHECK (
    (operation_type_id IN (1,2,3) AND amount < 0) OR
    (operation_type_id = 4 AND amount > 0)
  )
);

CREATE INDEX IF NOT EXISTS idx_transactions_account_id ON transactions(account_id);
