#!/usr/bin/env bash
set -euo pipefail

DB_NAME="${DB_NAME:-sub2api}"

psql -v ON_ERROR_STOP=1 -d "$DB_NAME" -f /Users/lhl/Desktop/code/sub2api-distributor/backend/migrations/001_distributor_tables.sql >/dev/null
psql -v ON_ERROR_STOP=1 -d "$DB_NAME" -f /Users/lhl/Desktop/code/sub2api-distributor/backend/migrations/002_seed_demo_operator.sql >/dev/null

psql -v ON_ERROR_STOP=1 -d "$DB_NAME" <<'SQL'
DO $$
DECLARE
  inviter_id BIGINT;
  invitee_id BIGINT;
BEGIN
  INSERT INTO users (
    email, password_hash, role, balance, concurrency, status, username, notes, created_at, updated_at
  )
  VALUES (
    'dist_demo@example.com',
    '$2a$10$f0MXKoCFb/6e6Pu0mcD0IefnNezTQZTXIaCmv5CIgVQLmf0QRQ6ae',
    'user',
    0,
    5,
    'active',
    'dist_demo',
    'sub2api-distributor demo distributor',
    NOW(),
    NOW()
  )
  ON CONFLICT (email) WHERE deleted_at IS NULL DO UPDATE
  SET password_hash = EXCLUDED.password_hash,
      role = EXCLUDED.role,
      status = EXCLUDED.status,
      username = EXCLUDED.username,
      notes = EXCLUDED.notes,
      updated_at = NOW()
  RETURNING id INTO inviter_id;

  INSERT INTO users (
    email, password_hash, role, balance, concurrency, status, username, notes, created_at, updated_at
  )
  VALUES (
    'invitee_demo@example.com',
    '$2a$10$f0MXKoCFb/6e6Pu0mcD0IefnNezTQZTXIaCmv5CIgVQLmf0QRQ6ae',
    'user',
    0,
    5,
    'active',
    'invitee_demo',
    'sub2api-distributor demo invitee',
    NOW(),
    NOW()
  )
  ON CONFLICT (email) WHERE deleted_at IS NULL DO UPDATE
  SET password_hash = EXCLUDED.password_hash,
      role = EXCLUDED.role,
      status = EXCLUDED.status,
      username = EXCLUDED.username,
      notes = EXCLUDED.notes,
      updated_at = NOW()
  RETURNING id INTO invitee_id;

  INSERT INTO user_affiliates (
    user_id, aff_code, inviter_id, aff_count, aff_quota, aff_history_quota, aff_rebate_rate_percent, aff_code_custom, aff_frozen_quota, created_at, updated_at
  )
  VALUES (
    inviter_id, 'DISTDEMO2026', NULL, 1, 20.4, 20.4, NULL, true, 0, NOW(), NOW()
  )
  ON CONFLICT (user_id) DO UPDATE
  SET aff_code = EXCLUDED.aff_code,
      inviter_id = EXCLUDED.inviter_id,
      aff_count = EXCLUDED.aff_count,
      aff_quota = EXCLUDED.aff_quota,
      aff_history_quota = EXCLUDED.aff_history_quota,
      aff_rebate_rate_percent = EXCLUDED.aff_rebate_rate_percent,
      aff_code_custom = EXCLUDED.aff_code_custom,
      aff_frozen_quota = EXCLUDED.aff_frozen_quota,
      updated_at = NOW();

  INSERT INTO user_affiliates (
    user_id, aff_code, inviter_id, aff_count, aff_quota, aff_history_quota, aff_rebate_rate_percent, aff_code_custom, aff_frozen_quota, created_at, updated_at
  )
  VALUES (
    invitee_id, 'INVITEEDEMO2026', inviter_id, 0, 0, 0, NULL, true, 0, NOW(), NOW()
  )
  ON CONFLICT (user_id) DO UPDATE
  SET aff_code = EXCLUDED.aff_code,
      inviter_id = EXCLUDED.inviter_id,
      aff_count = EXCLUDED.aff_count,
      aff_quota = EXCLUDED.aff_quota,
      aff_history_quota = EXCLUDED.aff_history_quota,
      aff_rebate_rate_percent = EXCLUDED.aff_rebate_rate_percent,
      aff_code_custom = EXCLUDED.aff_code_custom,
      aff_frozen_quota = EXCLUDED.aff_frozen_quota,
      updated_at = NOW();

  DELETE FROM user_affiliate_ledger
  WHERE user_id = inviter_id
    AND source_user_id = invitee_id
    AND action IN ('accrue', 'transfer');

  INSERT INTO user_affiliate_ledger (
    user_id, action, amount, source_user_id, source_order_id, frozen_until,
    balance_after, aff_quota_after, aff_frozen_quota_after, aff_history_quota_after,
    created_at, updated_at
  )
  VALUES
    (inviter_id, 'accrue', 8.00, invitee_id, NULL, NOW() - INTERVAL '10 day', NULL, 8.00, 0, 8.00, NOW() - INTERVAL '9 day', NOW() - INTERVAL '9 day'),
    (inviter_id, 'accrue', 4.00, invitee_id, NULL, NOW() - INTERVAL '8 day', NULL, 12.00, 0, 12.00, NOW() - INTERVAL '7 day', NOW() - INTERVAL '7 day'),
    (inviter_id, 'accrue', 8.40, invitee_id, NULL, NOW() - INTERVAL '6 day', NULL, 20.40, 0, 20.40, NOW() - INTERVAL '5 day', NOW() - INTERVAL '5 day');
END $$;

DELETE FROM distributor_profiles
WHERE user_id IN (
  SELECT id
  FROM users
  WHERE email = 'legacy_review_demo@example.com'
);

WITH demo_users AS (
  SELECT
    MAX(CASE WHEN email = 'dist_demo@example.com' THEN id END) AS inviter_id,
    MAX(CASE WHEN email = 'invitee_demo@example.com' THEN id END) AS invitee_id
  FROM users
  WHERE email IN ('dist_demo@example.com', 'invitee_demo@example.com')
)
INSERT INTO distributor_profiles (
  user_id, status, display_name, settlement_channel, settlement_account_name,
  settlement_account_no, settlement_account_extra, notes, created_at, updated_at
)
SELECT inviter_id, 'active', '演示分销商 A', 'alipay', '邀请人A', 'dist-demo@demo', '{"nickname":"邀请人A"}'::jsonb, '用于分销商端演示', NOW(), NOW()
FROM demo_users
ON CONFLICT (user_id) DO UPDATE
SET status = EXCLUDED.status,
    display_name = EXCLUDED.display_name,
    settlement_channel = EXCLUDED.settlement_channel,
    settlement_account_name = EXCLUDED.settlement_account_name,
    settlement_account_no = EXCLUDED.settlement_account_no,
    settlement_account_extra = EXCLUDED.settlement_account_extra,
    notes = EXCLUDED.notes,
    updated_at = NOW();

WITH demo_users AS (
  SELECT
    MAX(CASE WHEN email = 'dist_demo@example.com' THEN id END) AS inviter_id,
    MAX(CASE WHEN email = 'invitee_demo@example.com' THEN id END) AS invitee_id
  FROM users
  WHERE email IN ('dist_demo@example.com', 'invitee_demo@example.com')
)
DELETE FROM distributor_withdraw_events
WHERE request_id IN (
  SELECT id
  FROM distributor_withdraw_requests
  WHERE user_id IN ((SELECT inviter_id FROM demo_users), (SELECT invitee_id FROM demo_users))
);

WITH demo_users AS (
  SELECT
    MAX(CASE WHEN email = 'dist_demo@example.com' THEN id END) AS inviter_id,
    MAX(CASE WHEN email = 'invitee_demo@example.com' THEN id END) AS invitee_id
  FROM users
  WHERE email IN ('dist_demo@example.com', 'invitee_demo@example.com')
)
DELETE FROM distributor_withdraw_requests
WHERE user_id IN ((SELECT inviter_id FROM demo_users), (SELECT invitee_id FROM demo_users));

WITH demo_users AS (
  SELECT
    MAX(CASE WHEN email = 'dist_demo@example.com' THEN id END) AS inviter_id,
    MAX(CASE WHEN email = 'invitee_demo@example.com' THEN id END) AS invitee_id
  FROM users
  WHERE email IN ('dist_demo@example.com', 'invitee_demo@example.com')
),
operator_user AS (
  SELECT id AS operator_user_id
  FROM users
  WHERE email = 'operator_demo@local.dev' AND deleted_at IS NULL
  LIMIT 1
),
inserted AS (
  INSERT INTO distributor_withdraw_requests (
    request_no, user_id, amount, status,
    snapshot_total_earned, snapshot_internal_transferred_amount,
    snapshot_paying_before, snapshot_paid_before,
    snapshot_withdrawable_before, snapshot_withdrawable_after,
    applicant_remark, paid_at, paid_channel, paid_reference_no, paid_remark,
    created_at, updated_at
  )
  SELECT 'DW-DEMO-4-PAYING', inviter_id, 8.00, 'paying', 20.4, 0, 0, 0, 20.4, 12.4, '首笔线下打款处理中', NULL, '', '', '', NOW() - INTERVAL '2 day', NOW() - INTERVAL '2 day'
  FROM demo_users
  UNION ALL
  SELECT 'DW-DEMO-4-PAID', inviter_id, 4.00, 'paid', 20.4, 0, 8.0, 0, 12.4, 8.4, '第二笔已线下打款', NOW() - INTERVAL '1 day', 'alipay', 'ALI-DEMO-0001', '运营已完成打款', NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day'
  FROM demo_users
  RETURNING id, request_no, user_id, status, amount
)
INSERT INTO distributor_withdraw_events (request_id, action, operator_user_id, detail, created_at)
SELECT id, 'create', user_id, jsonb_build_object('amount', amount, 'seed', true),
  CASE request_no
    WHEN 'DW-DEMO-4-PAYING' THEN NOW() - INTERVAL '2 day'
    WHEN 'DW-DEMO-4-PAID' THEN NOW() - INTERVAL '1 day'
    ELSE NOW() - INTERVAL '12 hour'
  END
FROM inserted
UNION ALL
SELECT id, 'mark_paid', operator_user_id, jsonb_build_object('paid_channel', 'alipay', 'paid_reference_no', 'ALI-DEMO-0001', 'seed', true), NOW() - INTERVAL '1 day'
FROM inserted
CROSS JOIN operator_user
WHERE request_no = 'DW-DEMO-4-PAID';
SQL

echo "Demo data seeded into database: $DB_NAME"
