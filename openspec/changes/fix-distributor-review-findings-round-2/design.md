# fix-distributor-review-findings-round-2 Design

## Scope

This change addresses three concrete review findings inside the distributor
backend and demo tooling:

- atomic distributor enablement
- correct operator payout error semantics
- deterministic demo reset behavior

It does not change API contracts, add pagination, or modify the main `sub2api`
affiliate behavior outside the shared tables already used by this project.

## Design Decisions

### 1. Atomic Distributor Enablement

`UpsertProfile` currently persists `distributor_profiles` and only afterward
calls `ensureAffiliateIdentity`. If the second step fails, the first step is
already committed.

The fix is to move both steps into a single database transaction:

- begin tx
- upsert `distributor_profiles`
- if target status is `active`, ensure `user_affiliates` inside the same tx
- commit only after both succeed

`ensureAffiliateIdentity` will be refactored to accept a query/exec interface
compatible with both `*sql.DB` and `*sql.Tx`, following the main system's
affiliate repository pattern.

### 2. Operator Payout Error Mapping

`MarkPaid` and operator `Cancel` currently flatten most failures into `400`,
which causes internal error strings to bypass the generic 5xx sanitizer.

The fix is to centralize status mapping:

- `ErrWithdrawalNotFound` => `404`
- `ErrInvalidWithdrawalTransition` => `400`
- unexpected errors => `500`

No frontend API shape changes are required because the existing response body
already carries `message`.

### 3. Deterministic Demo Reset

`seed_demo_data.sh` needs to reset all demo-affecting data, not just a subset.

The script will additionally:

- delete affiliate `transfer` ledger rows for the demo inviter/invitee chain
- remove non-demo distributor profiles that previous review/verification work
  left behind
- keep the operator seed intact

The goal is not to wipe the entire main system database, only the rows that this
project seeds or depends on for acceptance.

## Risks And Mitigations

- Transactional refactor risk: use focused tests around activation failure and
  success paths.
- Seed overreach risk: delete only targeted demo/test rows by email or note
  markers, not all distributor rows globally.
- Handler regression risk: add tests for known and unexpected operator payout
  failures.
