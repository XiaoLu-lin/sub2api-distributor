## ADDED Requirements

### Requirement: Project documentation remains consistent with the runnable system
The project SHALL provide consistent documentation for startup, local ports, demo data, verification commands, and current feature coverage across the root documentation set.

#### Scenario: Maintainer follows the setup guide
- **WHEN** a maintainer reads the root README and testing documentation
- **THEN** the documented backend port, recommended frontend port, startup commands, demo account strategy, and verification commands MUST match the current runnable project

#### Scenario: Frontend contributors open the frontend README
- **WHEN** a contributor reads `frontend/README.md`
- **THEN** the file MUST describe the distributor frontend’s actual scripts, purpose, and test/build commands instead of the default Vite template text

### Requirement: The frontend exposes a runnable automated test baseline
The frontend SHALL provide an automated test command that validates key non-visual behaviors without depending on manual browser interaction.

#### Scenario: Contributor runs frontend tests
- **WHEN** a contributor runs the documented frontend test command
- **THEN** the project MUST execute automated tests covering critical utilities or session/request helpers and report pass or fail status

#### Scenario: Frontend helper behavior regresses
- **WHEN** formatting, session storage, navigation configuration, or base request behavior changes incompatibly
- **THEN** at least one automated frontend test MUST fail and surface the regression

### Requirement: The backend exposes a documented and tested quality baseline
The backend SHALL document its exported API surface and validate critical pure logic through automated tests.

#### Scenario: Maintainer reads backend core packages
- **WHEN** a maintainer opens backend core files under `backend/internal`
- **THEN** exported types, exported functions, and key business flow helpers MUST have comments that explain their role and usage boundary

#### Scenario: Authentication or configuration logic regresses
- **WHEN** token parsing, environment loading, withdrawal logic, or CORS allowlisting changes incompatibly
- **THEN** automated backend tests MUST fail and surface the regression

### Requirement: Acceptance coverage tracks the live API surface
The local API acceptance script SHALL cover the currently supported portal and operator endpoints that are required for end-to-end verification.

#### Scenario: Maintainer runs the acceptance script
- **WHEN** `scripts/api_acceptance.sh` is executed against a seeded local environment
- **THEN** the script MUST verify the portal invite metadata endpoint and the operator user lookup endpoint in addition to the existing core flows

#### Scenario: Acceptance results are reviewed later
- **WHEN** a maintainer inspects the latest acceptance artifact
- **THEN** the generated report MUST reflect the expanded endpoint coverage so the verified scope is visible
