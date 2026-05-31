CREATE TABLE IF NOT EXISTS distributor_profiles (
  user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  status VARCHAR(20) NOT NULL DEFAULT 'active',
  display_name VARCHAR(100) NOT NULL DEFAULT '',
  settlement_channel VARCHAR(20) NOT NULL DEFAULT '',
  settlement_account_name VARCHAR(100) NOT NULL DEFAULT '',
  settlement_account_no VARCHAR(255) NOT NULL DEFAULT '',
  settlement_account_extra JSONB NOT NULL DEFAULT '{}'::jsonb,
  notes TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS distributor_withdraw_requests (
  id BIGSERIAL PRIMARY KEY,
  request_no VARCHAR(40) NOT NULL UNIQUE,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  amount NUMERIC(20,8) NOT NULL,
  status VARCHAR(20) NOT NULL,
  snapshot_total_earned NUMERIC(20,8) NOT NULL DEFAULT 0,
  snapshot_internal_transferred_amount NUMERIC(20,8) NOT NULL DEFAULT 0,
  snapshot_paying_before NUMERIC(20,8) NOT NULL DEFAULT 0,
  snapshot_paid_before NUMERIC(20,8) NOT NULL DEFAULT 0,
  snapshot_withdrawable_before NUMERIC(20,8) NOT NULL DEFAULT 0,
  snapshot_withdrawable_after NUMERIC(20,8) NOT NULL DEFAULT 0,
  applicant_remark TEXT NOT NULL DEFAULT '',
  paid_at TIMESTAMPTZ NULL,
  paid_channel VARCHAR(20) NOT NULL DEFAULT '',
  paid_reference_no VARCHAR(100) NOT NULL DEFAULT '',
  paid_remark TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_distributor_withdraw_requests_user_status
  ON distributor_withdraw_requests (user_id, status, created_at DESC);

CREATE TABLE IF NOT EXISTS distributor_withdraw_events (
  id BIGSERIAL PRIMARY KEY,
  request_id BIGINT NOT NULL REFERENCES distributor_withdraw_requests(id) ON DELETE CASCADE,
  action VARCHAR(30) NOT NULL,
  operator_user_id BIGINT NULL REFERENCES users(id) ON DELETE SET NULL,
  detail JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_distributor_withdraw_events_request_id
  ON distributor_withdraw_events (request_id, created_at DESC);
