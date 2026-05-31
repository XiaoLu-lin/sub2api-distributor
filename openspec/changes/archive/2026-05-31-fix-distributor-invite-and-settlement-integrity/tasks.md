## 1. Invite Link Integrity

- [x] 1.1 Replace portal invite link generation so it points to the configured main-system registration base URL instead of the distributor console origin.
- [x] 1.2 Add a clear fallback/error state when the main-system registration base URL is missing in local or deployment configuration.

## 2. Distributor Enablement Integrity

- [x] 2.1 Add backend logic that ensures an active distributor profile also has a corresponding `user_affiliates` record with a valid invite code.
- [x] 2.2 Cover the enablement path with tests so newly enabled distributors can immediately fetch invite metadata.

## 3. Settlement Amount Integrity

- [x] 3.1 Refactor backend settlement calculations from `float64`-driven arithmetic to decimal-backed arithmetic while preserving the existing API contract.
- [x] 3.2 Add or update tests for withdrawable amount, paying amount, paid amount, and withdrawal validation using decimal-safe comparisons.

## 4. Demo And Acceptance Alignment

- [x] 4.1 Fix `scripts/seed_demo_data.sh` so the default invited user is not also seeded as a distributor.
- [x] 4.2 Update `scripts/api_acceptance.sh` and project docs so the default demo accounts match the current recommended end-to-end flow.

## 5. Verification

- [x] 5.1 Run backend tests covering auth, distributor logic, and any new helper logic.
- [x] 5.2 Run frontend tests and frontend build.
- [x] 5.3 Run the acceptance script against the corrected demo dataset and record the result.
