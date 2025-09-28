INSERT INTO operation_types (id, description) VALUES
  (1, 'CASH PURCHASE'),
  (2, 'INSTALLMENT PURCHASE'),
  (3, 'WITHDRAWAL'),
  (4, 'PAYMENT')
ON CONFLICT (id) DO NOTHING;
