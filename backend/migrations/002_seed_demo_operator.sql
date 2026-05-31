-- Demo operator account for local acceptance testing only.
-- Password: Distributor123!

DO $$
BEGIN
  IF EXISTS (SELECT 1 FROM users WHERE email = 'operator_demo@local.dev' AND deleted_at IS NULL) THEN
    UPDATE users
    SET password_hash = '$2a$10$hdYzPNjbRzzY3tuckj.A2.K2lBBYPwd3hhXvvbMb2szuIPXsWURPC',
        role = 'admin',
        status = 'active',
        username = 'operator_demo',
        notes = 'sub2api-distributor local acceptance operator',
        updated_at = NOW()
    WHERE email = 'operator_demo@local.dev' AND deleted_at IS NULL;
  ELSE
    INSERT INTO users (
      email,
      password_hash,
      role,
      balance,
      concurrency,
      status,
      username,
      notes,
      created_at,
      updated_at
    )
    VALUES (
      'operator_demo@local.dev',
      '$2a$10$hdYzPNjbRzzY3tuckj.A2.K2lBBYPwd3hhXvvbMb2szuIPXsWURPC',
      'admin',
      0,
      1,
      'active',
      'operator_demo',
      'sub2api-distributor local acceptance operator',
      NOW(),
      NOW()
    );
  END IF;
END $$;
