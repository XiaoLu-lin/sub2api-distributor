## 1. Atomic Distributor Enablement

- [x] 1.1 Refactor distributor profile upsert so `distributor_profiles` and
      `user_affiliates` creation run inside one transaction.
- [x] 1.2 Add regression tests that prove failed affiliate-identity creation
      does not leave a newly active distributor behind.

## 2. Operator Error Semantics

- [x] 2.1 Fix operator payout mutation handlers so unexpected errors return 500,
      while known validation/state errors keep their current client-error codes.
- [x] 2.2 Add handler tests covering both expected and unexpected payout errors.

## 3. Demo Reset Determinism

- [x] 3.1 Update `scripts/seed_demo_data.sh` to clear demo transfer rows and
      stale demo distributor profiles before reseeding.
- [x] 3.2 Update docs if the reset scope or seed assumptions change.

## 4. Verification

- [x] 4.1 Run backend tests.
- [x] 4.2 Run frontend tests and build if touched.
- [x] 4.3 Rerun demo seed plus API acceptance and confirm baseline output.
