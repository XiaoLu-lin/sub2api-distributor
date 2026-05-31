## 1. OpenSpec And Documentation Baseline

- [x] 1.1 Write the proposal, design, and capability spec for the project quality baseline change.
- [x] 1.2 Rewrite `README.md`, `TESTING.md`, `docs/distributor-system-overview.md`, and `frontend/README.md` so ports, startup instructions, accounts, feature coverage, and verification commands stay consistent.

## 2. Frontend Verification Baseline

- [x] 2.1 Add a lightweight frontend test runner and script that fits the existing Vite + TypeScript stack.
- [x] 2.2 Add automated frontend tests for formatting helpers, session storage helpers, navigation definitions, and base HTTP request behavior.

## 3. Backend Verification And Readability

- [x] 3.1 Add backend tests for authentication token parsing and configuration environment loading while preserving existing logic tests.
- [x] 3.2 Add clear comments to exported backend types, functions, constants, and key transactional or routing logic.

## 4. Acceptance Coverage And Verification

- [x] 4.1 Update `scripts/api_acceptance.sh` so it covers the current invite metadata and operator user lookup endpoints.
- [x] 4.2 Run backend tests, frontend tests, frontend build, and the API acceptance flow; then record the resulting project state in the updated documentation.
