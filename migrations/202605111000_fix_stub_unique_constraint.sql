-- Fix: Replace the absolute UNIQUE constraint with a partial unique index.
-- Stub agents (identity_registry='' AND agent_id='') are NOT unique on (chain_id, identity_registry, agent_id)
-- because multiple stub agents on the same chain will all have (chain_id, '', '').
-- Only agents with actual identity (non-empty identity_registry OR non-empty agent_id) are unique.
--
-- Steps:
-- 1. Drop the bad constraint (created by 202605080001)
-- 2. Delete duplicate stubs (keep lowest uid per wallet+chain, since stubs are identified by wallet)
-- 3. Create a partial unique index for non-stub agents only

-- Step 1: Drop the problematic unique constraint
ALTER TABLE agents DROP CONSTRAINT IF EXISTS uq_agents_chain_registry_agent;

-- Step 2: For stub agents (empty identity_registry AND empty agent_id),
-- keep only the lowest uid per (chain_id, agent_wallet) — remove all others.
-- These are duplicate stubs created by concurrent FindOrCreateAgentByWallet calls.
DELETE FROM agents a
WHERE (a.identity_registry = '' AND a.agent_id = '')
  AND a.uid > (
    SELECT MIN(a2.uid)
    FROM agents a2
    WHERE a2.chain_id = a.chain_id
      AND a2.agent_wallet = a.agent_wallet
  );

-- Step 3: Create partial unique index for non-stub agents only.
-- This means: for any (chain_id, identity_registry, agent_id) combo where
-- at least one of identity_registry or agent_id is non-empty, enforce uniqueness.
CREATE UNIQUE INDEX IF NOT EXISTS uq_agents_identity_nonstub
ON agents (chain_id, identity_registry, agent_id)
WHERE identity_registry <> '' OR agent_id <> '';
