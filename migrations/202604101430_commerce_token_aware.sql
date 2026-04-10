-- ERC-8183 Commerce Token-Aware Budget: add token fields and USD volume fields.
-- Spec: docs/specs/erc8183-commerce-token-aware.md §2-§3

DO $$
BEGIN
  -- commerce_actions: token-aware budget fields (§2.1)
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'commerce_actions' AND column_name = 'payment_token'
  ) THEN
    ALTER TABLE commerce_actions ADD COLUMN payment_token varchar(255) DEFAULT '';
  END IF;

  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'commerce_actions' AND column_name = 'token_symbol'
  ) THEN
    ALTER TABLE commerce_actions ADD COLUMN token_symbol varchar(32) DEFAULT '';
  END IF;

  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'commerce_actions' AND column_name = 'budget_usd'
  ) THEN
    ALTER TABLE commerce_actions ADD COLUMN budget_usd numeric(36, 8) DEFAULT 0;
  END IF;

  -- Indexes for token-aware queries (§3)
  IF NOT EXISTS (
    SELECT 1 FROM pg_indexes WHERE indexname = 'idx_ca_payment_token'
  ) THEN
    CREATE INDEX idx_ca_payment_token ON commerce_actions(payment_token);
  END IF;

  IF NOT EXISTS (
    SELECT 1 FROM pg_indexes WHERE indexname = 'idx_ca_token_symbol'
  ) THEN
    CREATE INDEX idx_ca_token_symbol ON commerce_actions(token_symbol);
  END IF;
END $$;

DO $$
BEGIN
  -- commerce_scores: USD volume fields (§2.2)
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'commerce_scores' AND column_name = 'total_volume_usd'
  ) THEN
    ALTER TABLE commerce_scores ADD COLUMN total_volume_usd numeric(36, 8) DEFAULT 0;
  END IF;

  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'commerce_scores' AND column_name = 'weighted_volume_usd_sum'
  ) THEN
    ALTER TABLE commerce_scores ADD COLUMN weighted_volume_usd_sum numeric(36, 8) DEFAULT 0;
  END IF;
END $$;

DO $$
BEGIN
  -- commerce_scores_global: USD volume fields (§2.2)
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'commerce_scores_global' AND column_name = 'total_volume_usd'
  ) THEN
    ALTER TABLE commerce_scores_global ADD COLUMN total_volume_usd numeric(36, 8) DEFAULT 0;
  END IF;

  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'commerce_scores_global' AND column_name = 'weighted_volume_usd_sum'
  ) THEN
    ALTER TABLE commerce_scores_global ADD COLUMN weighted_volume_usd_sum numeric(36, 8) DEFAULT 0;
  END IF;
END $$;