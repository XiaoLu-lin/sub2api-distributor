# Distributor Integrity Fixes Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fix invite link correctness, distributor enablement completeness, settlement arithmetic precision, and demo/acceptance alignment in `sub2api-distributor`.

**Architecture:** The backend will gain a focused affiliate-identity helper inside the distributor domain so active distributor enablement can ensure `user_affiliates` rows without importing the main system service layer. Settlement calculations will move to decimal-backed arithmetic internally while preserving the current API shape. The frontend will stop guessing the registration origin and instead render links from an explicit main-app base URL with a safe fallback state.

**Tech Stack:** Go, Gin, PostgreSQL, Vue 3, TypeScript, Vite, Vitest, shell scripts, OpenSpec

---

### Task 1: Add failing tests for invite-link configuration behavior

**Files:**
- Create: None
- Modify: `/Users/lhl/Desktop/code/sub2api-distributor/frontend/src/utils/format.test.ts`
- Modify: `/Users/lhl/Desktop/code/sub2api-distributor/frontend/src/views/portal/PortalInviteesView.vue`
- Test: `/Users/lhl/Desktop/code/sub2api-distributor/frontend/src/utils/format.test.ts`

- [ ] **Step 1: Write the failing tests**

Add tests that codify:
- a configured main-app base URL produces `https://main.example.com/register?aff=CODE`
- a missing base URL does not produce a fake local `/register` link

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/lhl/Desktop/code/sub2api-distributor/frontend && pnpm test -- --run src/utils/format.test.ts`
Expected: FAIL because invite-link builder does not exist yet

- [ ] **Step 3: Write minimal implementation**

Introduce a small helper for invite link building and update `PortalInviteesView.vue` to use it instead of `window.location.origin`.

- [ ] **Step 4: Run test to verify it passes**

Run: `cd /Users/lhl/Desktop/code/sub2api-distributor/frontend && pnpm test -- --run src/utils/format.test.ts`
Expected: PASS

### Task 2: Add failing backend tests for affiliate identity ensure-on-enable

**Files:**
- Modify: `/Users/lhl/Desktop/code/sub2api-distributor/backend/internal/distributor/logic_test.go`
- Modify: `/Users/lhl/Desktop/code/sub2api-distributor/backend/internal/distributor/service.go`
- Create: `/Users/lhl/Desktop/code/sub2api-distributor/backend/internal/distributor/affiliate_identity.go`
- Test: `/Users/lhl/Desktop/code/sub2api-distributor/backend/internal/distributor/logic_test.go`

- [ ] **Step 1: Write the failing tests**

Add tests for:
- deterministic affiliate code normalization / validation helpers if introduced
- ensure-active-profile path requiring affiliate identity creation semantics

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/lhl/Desktop/code/sub2api-distributor/backend && go test ./internal/distributor -run 'TestEnsure|TestAffiliate'`
Expected: FAIL because helper/logic is not implemented yet

- [ ] **Step 3: Write minimal implementation**

Add a focused helper that:
- checks for existing `user_affiliates`
- generates a unique invite code when missing
- is called by `UpsertProfile` only for `active` profiles

- [ ] **Step 4: Run test to verify it passes**

Run: `cd /Users/lhl/Desktop/code/sub2api-distributor/backend && go test ./internal/distributor -run 'TestEnsure|TestAffiliate'`
Expected: PASS

### Task 3: Add failing settlement arithmetic precision tests

**Files:**
- Modify: `/Users/lhl/Desktop/code/sub2api-distributor/backend/internal/distributor/logic_test.go`
- Modify: `/Users/lhl/Desktop/code/sub2api-distributor/backend/internal/distributor/logic.go`
- Modify: `/Users/lhl/Desktop/code/sub2api-distributor/backend/internal/distributor/types.go`
- Test: `/Users/lhl/Desktop/code/sub2api-distributor/backend/internal/distributor/logic_test.go`

- [ ] **Step 1: Write the failing tests**

Add tests that demonstrate decimal-sensitive cases such as:
- `0.1 + 0.2` style accumulation
- withdrawable amount exact subtraction across earned / paying / paid values
- boundary comparison where requested amount equals available balance

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/lhl/Desktop/code/sub2api-distributor/backend && go test ./internal/distributor -run 'TestComputeWithdrawableAmount|TestDecimal'`
Expected: FAIL because current implementation is float-based

- [ ] **Step 3: Write minimal implementation**

Switch internal arithmetic to `shopspring/decimal` in:
- summary calculations
- withdrawable computation
- amount validation comparisons

- [ ] **Step 4: Run test to verify it passes**

Run: `cd /Users/lhl/Desktop/code/sub2api-distributor/backend && go test ./internal/distributor -run 'TestComputeWithdrawableAmount|TestDecimal'`
Expected: PASS

### Task 4: Fix invite page fallback and operator/distributor behavior

**Files:**
- Modify: `/Users/lhl/Desktop/code/sub2api-distributor/frontend/src/views/portal/PortalInviteesView.vue`
- Modify: `/Users/lhl/Desktop/code/sub2api-distributor/frontend/src/types.ts`
- Modify: `/Users/lhl/Desktop/code/sub2api-distributor/frontend/src/api/portal.ts`
- Modify: `/Users/lhl/Desktop/code/sub2api-distributor/frontend/src/utils/format.ts`
- Test: `/Users/lhl/Desktop/code/sub2api-distributor/frontend/src/utils/format.test.ts`

- [ ] **Step 1: Write the failing test**

Add a test or assertion for the fallback text state when main-app base URL is absent.

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/lhl/Desktop/code/sub2api-distributor/frontend && pnpm test -- --run src/utils/format.test.ts`
Expected: FAIL because fallback behavior is not implemented

- [ ] **Step 3: Write minimal implementation**

Update the invite page to:
- show a real configured link when available
- show clear incomplete-configuration copy when unavailable
- keep copy buttons from copying fake links

- [ ] **Step 4: Run test to verify it passes**

Run: `cd /Users/lhl/Desktop/code/sub2api-distributor/frontend && pnpm test -- --run src/utils/format.test.ts`
Expected: PASS

### Task 5: Align seed and acceptance scripts with the supported flow

**Files:**
- Modify: `/Users/lhl/Desktop/code/sub2api-distributor/scripts/seed_demo_data.sh`
- Modify: `/Users/lhl/Desktop/code/sub2api-distributor/scripts/api_acceptance.sh`
- Modify: `/Users/lhl/Desktop/code/sub2api-distributor/README.md`
- Modify: `/Users/lhl/Desktop/code/sub2api-distributor/TESTING.md`
- Modify: `/Users/lhl/Desktop/code/sub2api-distributor/docs/distributor-system-overview.md`
- Test: acceptance shell execution

- [ ] **Step 1: Write the failing verification expectation**

Record the intended baseline:
- one distributor demo account
- one invitee demo account
- acceptance script defaults point to the distributor demo account

- [ ] **Step 2: Run the current script/data flow to observe mismatch**

Run: `bash /Users/lhl/Desktop/code/sub2api-distributor/scripts/api_acceptance.sh`
Expected: FAIL or exercise the wrong accounts against current intended flow

- [ ] **Step 3: Write minimal implementation**

Update scripts and docs so:
- invitee is not seeded as distributor
- acceptance defaults match current recommended demo accounts

- [ ] **Step 4: Run script to verify it now targets the right flow**

Run: `bash /Users/lhl/Desktop/code/sub2api-distributor/scripts/seed_demo_data.sh && bash /Users/lhl/Desktop/code/sub2api-distributor/scripts/api_acceptance.sh`
Expected: PASS with report artifact written

### Task 6: Full verification

**Files:**
- Modify: `/Users/lhl/Desktop/code/sub2api-distributor/openspec/changes/fix-distributor-invite-and-settlement-integrity/tasks.md`

- [ ] **Step 1: Run backend tests**

Run: `cd /Users/lhl/Desktop/code/sub2api-distributor/backend && go test ./...`
Expected: PASS

- [ ] **Step 2: Run frontend tests**

Run: `cd /Users/lhl/Desktop/code/sub2api-distributor/frontend && pnpm test -- --run`
Expected: PASS

- [ ] **Step 3: Run frontend build**

Run: `cd /Users/lhl/Desktop/code/sub2api-distributor/frontend && pnpm build`
Expected: PASS

- [ ] **Step 4: Run acceptance flow**

Run: `bash /Users/lhl/Desktop/code/sub2api-distributor/scripts/seed_demo_data.sh && bash /Users/lhl/Desktop/code/sub2api-distributor/scripts/api_acceptance.sh`
Expected: PASS and generate a new report under `/Users/lhl/Desktop/code/sub2api-distributor/test-results/`

- [ ] **Step 5: Mark OpenSpec tasks complete**

Update `/Users/lhl/Desktop/code/sub2api-distributor/openspec/changes/fix-distributor-invite-and-settlement-integrity/tasks.md` checkboxes to done after verification succeeds.
