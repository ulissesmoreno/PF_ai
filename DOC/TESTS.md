# TESTS

This file registers all test results executed by the AI agent and guides the developer. **It must be filled without fail after each test cycle to ensure complete traceability and avoid loss of validations.** Reminder: Functionality and security tests are mandatory before any code (TDD), and the file is updated after each delivery.

## Updated on: YYYY-MM-DD HH:MM:SS

## 1. Executive Summary
- Test cycle: [E.g.: MVP Phase 1]
- Environment: [E.g.: Local with Docker]
- General notes: [E.g.: All tests passed; focus on security.]

## 2. Executed Test Registry

### Test 1: [ID] - Test name
- Type: Unit / Integration / Acceptance / Security / Regression
- Objective: [Clear description]
- Preconditions: [E.g.: Environment configured]
- Step by step:
  1. [Step 1]
  2. [Step 2]
  3. [Step 3]
- Expected result: [Description]
- Obtained result: [Actual description]
- Status: Passed / Failed / Pending
- Notes: [Additional notes]
- Corrective action: [If failed, what to do]

---

### Test 2: [ID] - Test name
- Type: Unit / Integration / Acceptance / Security / Regression
- Objective: [Clear description]
- Preconditions: [E.g.: Environment configured]
- Step by step:
  1. [Step 1]
  2. [Step 2]
  3. [Step 3]
- Expected result: [Description]
- Obtained result: [Actual description]
- Status: Passed / Failed / Pending
- Notes: [Additional notes]
- Corrective action: [If failed, what to do]

## 3. AI Agent Results
- Tests executed automatically by the agent: [List]
- Hypotheses tested: [E.g.: JWT validation]
- Conclusions: [Summary]
- Recommended adjustments: [E.g.: Improve coverage]

## 4. Developer Guide — Step-by-Step Test Execution

> This guide must be followed **at the end of each ROADMAP stage**. No stage advances without all steps below being executed and recorded in this file.

### 4.1 Environment Preparation

1. Confirm Docker is active: `docker ps` (should list project containers).
2. Bring up the database and dependencies: `docker-compose up -d`.
3. Check environment variables: refer to `ENV_SETUP.md` and ensure `.env` is configured (never committed).
4. Confirm installed versions:
   - Java: `java -version` (expected: 21+)
   - Python: `python --version` (expected: 3.12+)
   - Node.js: `node -version` (expected: as per `package.json`)
5. Read the stage acceptance criteria in `PLAN.md` before running any test.

---

### 4.2 Unit Tests

**Objective:** validate the isolated logic of each class/function without external dependencies.

#### Java Backend
```bash
./mvnw test
# or
mvn test -pl <module-name>
```
- Expected result: `BUILD SUCCESS` with 0 failures.
- Check coverage: report generated at `target/site/jacoco/index.html`.
- Minimum acceptable coverage: **80% per domain class**.

#### Python Backend
```bash
pytest tests/unit/ -v --cov=src --cov-report=term-missing
```
- Expected result: all tests `PASSED`.
- Minimum acceptable coverage: **80%**.

#### Frontend (Angular)
```bash
npm test -- --watch=false --code-coverage
```
- Report at `coverage/index.html`.
- Minimum acceptable coverage: **70% for critical components**.

---

### 4.3 Integration Tests

**Objective:** validate communication between layers (API ↔ database, domain ↔ adapters).

1. Ensure the database is running: `docker-compose up -d db`.
2. Run pending migrations (if any): `./mvnw flyway:migrate` or equivalent script.

#### Java Backend
```bash
./mvnw verify -P integration-tests
# or with Spring profile:
./mvnw test -Dspring.profiles.active=test
```

#### Python Backend
```bash
pytest tests/integration/ -v
```

- Expected result: all tests `PASSED`.
- On connection failure, check `DB_URL` in `.env` and confirm the container is active.

---

### 4.4 Security Tests

**Objective:** ensure no route is exposed without authentication and sensitive data is protected.

1. **JWT — Route without token must return 401:**
   ```bash
   curl -i -X GET http://localhost:8080/api/v1/<protected-route>
   # Expected: HTTP/1.1 401 Unauthorized
   ```

2. **JWT — Route with valid token must return 200:**
   ```bash
   TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{"username":"user","password":"pass"}' | jq -r '.token')

   curl -i -H "Authorization: Bearer $TOKEN" \
     http://localhost:8080/api/v1/<protected-route>
   # Expected: HTTP/1.1 200 OK
   ```

3. **Logs free of sensitive data:** inspect logs and confirm passwords, tokens, and payloads do not appear in plain text.
   ```bash
   docker logs <backend-container> | grep -i "password\|secret\|token"
   # Expected: no occurrences with real values
   ```

4. **Input sanitization:** attempt to send malicious data (e.g.: SQL injection, XSS) in input fields and confirm the API returns an error without processing.

---

### 4.5 E2E Tests (End-to-End)

**Objective:** simulate the complete flow of a real user, from frontend to database.

1. Bring up all services: `docker-compose up`.
2. Access the frontend: `http://localhost:4200` (or configured port).

#### Minimum flow to validate per stage:
- [ ] Login with valid credentials → correct redirect.
- [ ] Login with invalid credentials → clear error message.
- [ ] Execution of the stage's main feature → expected result displayed.
- [ ] Logout → session ended, access to protected routes blocked.

#### E2E Automation (when available):
```bash
# Playwright / Cypress (Angular)
npx playwright test
# or
npx cypress run
```

---

### 4.6 Regression Tests

**Objective:** ensure previous deliveries continue working after new implementations.

1. Run the **entire** unit and integration test suite (sections 4.2 and 4.3).
2. Re-run the E2E flow for features from previous stages (section 4.5).
3. Compare results against records in `## 2. Executed Test Registry`.
4. Any regression must be registered as a blocker in `STATE.md` before advancing.

---

### 4.7 Approval Criteria (Definition of Done — Tests)

| Criterion | Minimum Acceptable |
| :--- | :--- |
| Unit coverage (domain) | ≥ 80% |
| Integration tests | 100% passing |
| Routes without JWT return 401 | 100% |
| Logs free of sensitive data | 100% |
| Main E2E flow | No critical errors |
| Regression on previous items | 0 new failures |

---

### 4.8 On Failure

1. **Do not advance** to the next ROADMAP stage.
2. Record the failure in this file (`## 2. Executed Test Registry`) with:
   - Type, ID, and description of the test.
   - Obtained vs. expected result.
   - Relevant logs (without sensitive data).
   - Planned corrective action.
3. Register as a blocker in `STATE.md`.
4. If the cause is a requirement ambiguity, open an entry in `QUESTIONS.md`.
5. Fix, re-run the test, and update the status to `Passed` before proceeding.

---

### 4.9 Documenting Results

After executing all tests:
1. Fill in `## 2. Executed Test Registry` for each executed test.
2. Update `## 5. Metrics and Coverage` with current percentages.
3. Update `## 6. Roadmap Integration` with the stage status.
4. Record the delivery in `STATE.md` (section "What Was Completed") with timestamp.
5. Update `ROADMAP.md` marking the stage as `[x]` only after user approval.

## 5. Metrics and Coverage
- **Test Coverage:** [E.g.: 85% unit, 70% integration]
- **Average Execution Time:** [E.g.: 5 min]
- **Identified Errors:** [Quantity and types]
- **Reference:** Consult STATE.md for general progress.

## 6. Roadmap Integration
- **Current Stage:** [Link to ROADMAP.md, e.g.: Phase 1 MVP]
- **Tests per Phase:** [Mandatory validations before advancing]

## 7. Version History (Immutable)
- **[YYYY-MM-DD HH:MM] - Responsible:** Initial template version.
- [Add new entries here without deleting.]
1. Prepare environment with listed dependencies.
2. Review acceptance criteria in `PLAN.md`.
3. Execute the tests listed in this section.
4. Compare obtained result with expected result.
5. Record failures, evidence, and corrective actions.
6. Update status and attach logs whenever possible.

## 5. Quality Checklist
- [ ] All acceptance criteria were validated.
- [ ] All business rules were tested.
- [ ] All security validations were executed.
- [ ] Results were documented with timestamp and status.
- [ ] Corrective actions were registered when necessary.

---

## 8. Phase 0 Validation Results

### [2026-05-05 15:16] TEST-DOC-0001 - Onboarding Placeholder Scan
- **Type:** Documentation / Regression
- **Objective:** Confirm required onboarding files no longer contain onboarding placeholder markers.
- **Preconditions:** Phase 0 onboarding files filled.
- **Step by step:**
  1. Scanned `README.md`, `DOC/PROJECT.md`, `DOC/DESIGN.md`, `DOC/ROADMAP.md`, `DOC/PLAN.md`, `DOC/ENV_SETUP.md`, and `DOC/ARCHITECTURE.md`.
  2. Used placeholder patterns for bracketed onboarding markers, project-name templates, and post-onboarding text.
  3. Reviewed command result.
- **Expected result:** No placeholder occurrences in required onboarding files.
- **Obtained result:** No occurrences found.
- **Status:** Passed.
- **Notes:** Documentation is ready for Phase 0 closure review.
- **Corrective action:** None.

### [2026-05-05 15:16] TEST-SEC-0001 - Documentation Secret Scan
- **Type:** Security / Documentation
- **Objective:** Confirm onboarding documentation does not contain real secrets.
- **Preconditions:** Phase 0 onboarding files filled.
- **Step by step:**
  1. Scanned onboarding files for secret-related terms.
  2. Reviewed each occurrence manually.
  3. Verified whether occurrences were variable names/instructions or real credentials.
- **Expected result:** No real credentials, tokens, passwords, or API keys in documentation.
- **Obtained result:** Only environment variable names and secret-handling instructions were found; no real secret values were present.
- **Status:** Passed.
- **Notes:** `DOC/ENV_SETUP.md` intentionally documents variable names such as `PF_AI_DB_PASSWORD`, `PF_AI_AUTH_SECRET`, and `PF_AI_API_KEY`.
- **Corrective action:** None.

### [2026-05-05 15:16] TEST-ARCH-0001 - Go and Hexagonal Feasibility
- **Type:** Architecture
- **Objective:** Confirm Go remains compatible with PF_ai's Hexagonal Architecture requirement.
- **Preconditions:** [HUMAN] asked about Go and Hexagonal Architecture in `NEW-INSTRUCTIONS.md`.
- **Step by step:**
  1. Reviewed `DOC/ARCHITECTURE.md`.
  2. Registered the architectural answer in `QUESTIONS.md`.
  3. Confirmed domain/application/infrastructure boundaries are documented for Go.
- **Expected result:** Go architecture keeps domain logic isolated from adapters.
- **Obtained result:** Go is approved for Hexagonal Architecture with pure domain packages, application ports/use cases, and infrastructure adapters.
- **Status:** Passed.
- **Notes:** CTO must enforce this boundary before Phase 1 implementation.
- **Corrective action:** None.

## 9. Phase 0 Closure Review Snapshot
- **QA:** Approved for documentation closure.
- **SECURITY:** Approved with no real secrets found in onboarding files.
- **CODE_REVIEWER:** Approved for documentation consistency.
- **TECH_LEAD:** Consolidated status: Approved for [HUMAN] Stage Closure Gate review.
- **Remaining blocker:** [HUMAN:Ulisses] must explicitly approve Phase 0 closure and roadmap advancement.
