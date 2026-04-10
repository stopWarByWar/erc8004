-- migrations/202604111000_commerce_jobs.sql
-- Commerce Jobs 表：存储所有 Job 的实时状态快照

CREATE TABLE commerce_jobs (
    uid bigserial PRIMARY KEY,
    chain_id varchar(255) NOT NULL,
    commerce_contract varchar(255) NOT NULL,
    job_id bigint NOT NULL,
    client varchar(255) NOT NULL,
    provider varchar(255) DEFAULT '',
    evaluator varchar(255) DEFAULT '',
    description text DEFAULT '',
    budget numeric(36,8) DEFAULT 0,
    paid_amount numeric(36,8) DEFAULT 0,
    paid_amount_usd numeric(36,8) DEFAULT 0,
    payment_token varchar(255) DEFAULT '',
    token_symbol varchar(32) DEFAULT '',
    status varchar(32) NOT NULL DEFAULT 'open',
    hook_address varchar(255) DEFAULT '',
    expired_at bigint DEFAULT 0,
    submitted_at bigint DEFAULT 0,
    completed_at bigint DEFAULT 0,
    latest_action_uid bigint,
    latest_block_number bigint,
    latest_tx_hash varchar(255),
    updated_at bigint NOT NULL,
    CONSTRAINT uniq_commerce_job UNIQUE (chain_id, commerce_contract, job_id)
);

CREATE INDEX idx_cj_chain_contract ON commerce_jobs(chain_id, commerce_contract);
CREATE INDEX idx_cj_client ON commerce_jobs(client);
CREATE INDEX idx_cj_provider ON commerce_jobs(provider);
CREATE INDEX idx_cj_status ON commerce_jobs(status);
CREATE INDEX idx_cj_updated ON commerce_jobs(updated_at);
