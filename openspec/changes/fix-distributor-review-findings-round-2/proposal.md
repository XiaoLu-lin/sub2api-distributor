# fix-distributor-review-findings-round-2 Proposal

## Why

The archived `fix-distributor-invite-and-settlement-integrity` change improved
invite-link correctness, distributor enablement, and settlement output
stability, but a follow-up code review identified three integrity gaps that
still affect real behavior:

1. Enabling a distributor writes `distributor_profiles` before ensuring the
   matching `user_affiliates` row exists, so a failed identity insert can leave
   a user able to log in without a usable invite code.
2. Operator payout mutation endpoints currently return `400` for unexpected
   backend failures, which leaks internal error text instead of the intended
   generic Chinese 5xx message.
3. The demo seed script is not a full reset. It preserves historical
   distributor rows and affiliate transfer ledger rows, so acceptance data can
   drift across repeated runs.

## What Changes

- Make distributor activation and affiliate identity provisioning atomic.
- Correct operator payout handler status mapping so unexpected failures are
  treated as server errors and hidden from the frontend.
- Make demo seed data reset all distributor demo state relevant to acceptance,
  including historical seeded distributor rows and demo transfer ledger data.
- Add regression tests for the new transactional and error-handling behavior.

## Impact

- Distributor users will no longer enter a partially enabled state.
- Operator payout APIs will produce safer and more accurate HTTP responses.
- Local demo and acceptance runs will become repeatable and deterministic.
