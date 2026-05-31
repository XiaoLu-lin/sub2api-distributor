## ADDED Requirements

### Requirement: Distributor Enablement Must Be Atomic

When the system enables or re-activates a distributor profile, it must not
persist an active distributor profile unless the corresponding source-system
affiliate identity also exists after the same operation completes.

#### Scenario: affiliate identity creation fails

- **WHEN** an operator updates a distributor profile to `active`
- **AND** the source-system affiliate identity cannot be created
- **THEN** the request fails
- **AND** the distributor profile is not left newly active because of that
  failed request

#### Scenario: affiliate identity creation succeeds

- **WHEN** an operator updates a distributor profile to `active`
- **THEN** the distributor profile is persisted
- **AND** the user can immediately fetch invite metadata afterward

### Requirement: Operator Payout Mutation Errors Must Use Correct HTTP Semantics

Unexpected backend failures during operator payout mutations must be reported as
server errors rather than client-validation failures.

#### Scenario: operator mark-paid hits unexpected backend error

- **WHEN** `/api/ops/withdrawals/:id/mark-paid` encounters an unexpected backend
  failure
- **THEN** the response status is `500`
- **AND** the response message is the generic server-safe Chinese message

#### Scenario: operator cancel hits invalid state transition

- **WHEN** `/api/ops/withdrawals/:id/cancel` is called for a request that cannot
  transition anymore
- **THEN** the response status is `400`
- **AND** the response message stays user-readable

### Requirement: Demo Seed Must Reset Acceptance State Deterministically

Running the distributor demo seed repeatedly must restore the same demo and
acceptance baseline relevant to this project.

#### Scenario: repeated seed after prior demo interactions

- **WHEN** `scripts/seed_demo_data.sh` runs after previous withdrawals,
  transfers, or extra distributor demo rows were created
- **THEN** the seeded distributor summary matches the documented baseline
- **AND** acceptance scripts do not inherit stale demo transfer or historical
  distributor state
