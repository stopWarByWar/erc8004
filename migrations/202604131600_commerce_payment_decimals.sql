-- Add payment_decimals to commerce_jobs and commerce_actions.
-- Purpose: correctly scale token amounts (budget/payment/fees) in API.
-- Default behavior: do NOT backfill historical rows; missing decimals treated as 18 by application.

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'commerce_jobs' AND column_name = 'payment_decimals'
  ) THEN
    ALTER TABLE commerce_jobs ADD COLUMN payment_decimals smallint DEFAULT 18;
  END IF;
END $$;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'commerce_actions' AND column_name = 'payment_decimals'
  ) THEN
    ALTER TABLE commerce_actions ADD COLUMN payment_decimals smallint DEFAULT 18;
  END IF;
END $$;

