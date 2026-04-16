-- Add USD budget + payment breakdown fields to commerce_jobs.
-- Meanings:
-- - budget: token amount set by BudgetSet
-- - budget_usd: USD value of budget at BudgetSet time
-- - paid_amount: token amount actually released to provider (net) at completion
-- - paid_amount_usd: USD value of paid_amount at completion time
-- - platform_fee_amount/platform_fee_usd: platform fee at completion
-- - evaluator_fee_amount/evaluator_fee_usd: evaluator fee at completion

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'commerce_jobs' AND column_name = 'budget_usd'
  ) THEN
    ALTER TABLE commerce_jobs ADD COLUMN budget_usd numeric(36, 8) DEFAULT 0;
  END IF;

  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'commerce_jobs' AND column_name = 'platform_fee_amount'
  ) THEN
    ALTER TABLE commerce_jobs ADD COLUMN platform_fee_amount numeric(36, 8) DEFAULT 0;
  END IF;

  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'commerce_jobs' AND column_name = 'platform_fee_usd'
  ) THEN
    ALTER TABLE commerce_jobs ADD COLUMN platform_fee_usd numeric(36, 8) DEFAULT 0;
  END IF;

  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'commerce_jobs' AND column_name = 'evaluator_fee_amount'
  ) THEN
    ALTER TABLE commerce_jobs ADD COLUMN evaluator_fee_amount numeric(36, 8) DEFAULT 0;
  END IF;

  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'commerce_jobs' AND column_name = 'evaluator_fee_usd'
  ) THEN
    ALTER TABLE commerce_jobs ADD COLUMN evaluator_fee_usd numeric(36, 8) DEFAULT 0;
  END IF;
END $$;

