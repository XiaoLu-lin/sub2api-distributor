# distributor-integrity-baseline Specification

## ADDED Requirements

### Requirement: Invite links must point to the real main-system registration entry

The distributor portal SHALL expose invite links that lead users into the main system registration flow instead of pointing back to the distributor console itself.

#### Scenario: Distributor copies invite link

- **WHEN** a distributor opens the invite page and copies the invite link
- **THEN** the generated URL MUST target the configured main-system registration base URL
- **AND** the URL MUST include the distributor invite code through the `aff` query parameter

#### Scenario: Main-system registration base URL is missing

- **WHEN** the frontend cannot resolve the configured main-system registration base URL
- **THEN** it MUST avoid generating a misleading local `/register` link
- **AND** it MUST present a clear fallback state so operators know configuration is incomplete

### Requirement: Enabling a distributor must produce a complete distributor identity

The operator enablement flow SHALL ensure that an active distributor can immediately access invite metadata without requiring manual database patching.

#### Scenario: Operator enables a historical user as distributor

- **WHEN** an operator creates or enables an active `distributor_profiles` row for a user whose `user_affiliates` record does not yet exist
- **THEN** the backend MUST create or ensure the corresponding affiliate identity during the same enablement flow

#### Scenario: Enabled distributor loads invite metadata

- **WHEN** a newly enabled distributor requests `/api/portal/invite-meta`
- **THEN** the endpoint MUST return a valid invite code instead of a missing-record error

### Requirement: Settlement arithmetic must be precision-safe

The distributor settlement backend SHALL compute rebate totals, withdrawal balances, and withdrawal validations without using floating-point arithmetic as the source of truth.

#### Scenario: Multiple rebate and withdrawal amounts are accumulated

- **WHEN** the backend aggregates earned, transferred, paying, and paid amounts
- **THEN** the intermediate arithmetic MUST use decimal-safe calculations
- **AND** the resulting withdrawable amount MUST match exact decimal expectations

#### Scenario: Withdrawal validation compares amount to available balance

- **WHEN** a distributor submits a withdrawal request near the current withdrawable balance boundary
- **THEN** the validation MUST compare values with decimal-safe precision
- **AND** it MUST not reject or allow the request because of floating-point rounding noise

### Requirement: Demo and acceptance assets must reflect the supported flow

The demo seed and acceptance tooling SHALL model the intended business roles and default credentials used for local verification.

#### Scenario: Demo data is seeded

- **WHEN** `scripts/seed_demo_data.sh` is executed
- **THEN** it MUST seed at least one distributor and one invited user
- **AND** the invited user MUST NOT also be marked as a distributor unless the script explicitly documents that variant

#### Scenario: Acceptance script is executed

- **WHEN** `scripts/api_acceptance.sh` runs with its default credentials
- **THEN** the defaults MUST match the current recommended seeded demo accounts
- **AND** the verified flow MUST align with the documented distributor and operator journey
