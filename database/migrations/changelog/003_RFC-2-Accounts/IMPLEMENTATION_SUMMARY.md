# RFC-2-Accounts: Liquibase Migration Summary

Generated: 2026-07-28

## Overview

This migration implements the complete Accounts domain for the Balance personal finance application, including all 14 entities specified in RFC-2-Accounts.md with full audit support.

---

## Migration Structure

### Total Changesets: 39

**Breakdown:**
- 14 table creation changesets
- 12 audit column changesets (note: audit_logs doesn't need audit columns)
- 12 audit foreign key changesets
- 1 post-migration configuration

---

## Entities Implemented

### 1. CURRENCIES
**Purpose:** Canonical list of currencies used across the system  
**Type:** Lookup table  
**Columns:** code (PK), name, symbol, decimal_places  
**Audit:** Yes

**Changesets:**
- `1785222122__currencies.yml` — Table creation
- `1785222263__currencies-audit.yml` — Audit columns
- `1785222280__currencies-audit-fk.yml` — Audit FKs

---

### 2. ACCOUNT_TYPES
**Purpose:** Classify sub-account behavior (CHECKING, SAVINGS, CREDIT_CARD)  
**Type:** Lookup table  
**Columns:** id (PK), code (UK), display_name  
**Audit:** Yes  
**Indexes:** Unique on code

**Changesets:**
- `1785222139__account-types.yml` — Table creation
- `1785222263__account-types-audit.yml` — Audit columns
- `1785222280__account-types-audit-fk.yml` — Audit FKs

---

### 3. ACCOUNTS
**Purpose:** Logical container/owner for sub-accounts (user's account or household)  
**Type:** Ownership entity  
**Columns:** id (PK, UUID), owner_id (FK→users), name, description, is_archived, created_at, updated_at  
**Audit:** Yes  
**Foreign Keys:**
- owner_id → users.id

**Changesets:**
- `1785222139__accounts.yml` — Table creation + FK to users
- `1785222263__accounts-audit.yml` — Audit columns
- `1785222280__accounts-audit-fk.yml` — Audit FKs

---

### 4. SUB_ACCOUNTS
**Purpose:** Primary operational account entity (holds balances, receives transactions)  
**Type:** Core financial entity  
**Columns:** id (PK, UUID), parent_account_id (FK), account_type_id (FK), currency_code (FK), name, description, opening_balance, current_balance, is_archived, associated_color, created_at, updated_at  
**Audit:** Yes  
**Foreign Keys:**
- parent_account_id → accounts.id
- account_type_id → account_types.id
- currency_code → currencies.code

**Indexes:**
- parent_account_id
- account_type_id
- currency_code
- (parent_account_id, is_archived) [composite]

**Changesets:**
- `1785222139__sub-accounts.yml` — Table creation + FKs + indexes
- `1785222263__sub-accounts-audit.yml` — Audit columns
- `1785222280__sub-accounts-audit-fk.yml` — Audit FKs

---

### 5. ACCOUNT_BALANCES
**Purpose:** Performance cache for available vs cleared balances  
**Type:** Cache table (one-to-one with SUB_ACCOUNTS)  
**Columns:** sub_account_id (PK, FK), available_balance, cleared_balance, updated_at  
**Audit:** No (cache table)  
**Foreign Keys:**
- sub_account_id → sub_accounts.id (unique, one-to-one)

**Changesets:**
- `1785222139__account-balances.yml` — Table creation + FK

---

### 6. CATEGORIES
**Purpose:** Hierarchical classification for transactions, budgets, recurring transactions  
**Type:** Classification entity  
**Columns:** id (PK, UUID), parent_category_id (FK, self-referencing), name, category_type  
**Audit:** Yes  
**Foreign Keys:**
- parent_category_id → categories.id (self-reference)

**Indexes:**
- parent_category_id
- category_type

**Changesets:**
- `1785222139__categories.yml` — Table creation + self-referencing FK + indexes
- `1785222263__categories-audit.yml` — Audit columns
- `1785222280__categories-audit-fk.yml` — Audit FKs

---

### 7. TRANSACTIONS
**Purpose:** Immutable ledger of debits/credits applied to sub-accounts  
**Type:** Ledger entity (append-only)  
**Columns:** id (PK, UUID), sub_account_id (FK), category_id (FK), amount, transaction_type, description, transaction_time, created_at  
**Audit:** Yes  
**Foreign Keys:**
- sub_account_id → sub_accounts.id
- category_id → categories.id

**Indexes:**
- sub_account_id
- (sub_account_id, transaction_time) [composite]
- category_id

**Changesets:**
- `1785222139__transactions.yml` — Table creation + FKs + indexes
- `1785222264__transactions-audit.yml` — Audit columns
- `1785222280__transactions-audit-fk.yml` — Audit FKs

---

### 8. TRANSFERS
**Purpose:** Represent funds movement between two sub-accounts  
**Type:** Financial transaction entity  
**Columns:** id (PK, UUID), from_sub_account_id (FK), to_sub_account_id (FK), amount, transfer_time, note  
**Audit:** Yes  
**Foreign Keys:**
- from_sub_account_id → sub_accounts.id
- to_sub_account_id → sub_accounts.id

**Indexes:**
- from_sub_account_id
- to_sub_account_id
- transfer_time

**Changesets:**
- `1785222139__transfers.yml` — Table creation + FKs + indexes
- `1785222281__transfers-audit-fk.yml` — Audit FKs
- `1785222264__transfers-audit.yml` — Audit columns

---

### 9. BUDGETS
**Purpose:** Planned spending limits over a time window for a category and currency  
**Type:** Planning entity  
**Columns:** id (PK, UUID), category_id (FK), currency_code (FK), name, amount_limit, start_date, end_date, created_at  
**Audit:** Yes  
**Foreign Keys:**
- category_id → categories.id
- currency_code → currencies.code

**Indexes:**
- category_id
- (category_id, start_date, end_date) [composite]

**Changesets:**
- `1785222139__budgets.yml` — Table creation + FKs + indexes
- `1785222264__budgets-audit.yml` — Audit columns
- `1785222281__budgets-audit-fk.yml` — Audit FKs

---

### 10. GOALS
**Purpose:** Savings targets used to track accumulation toward target_amount  
**Type:** Planning entity  
**Columns:** id (PK, UUID), sub_account_id (FK), currency_code (FK), name, target_amount, current_amount, target_date, completed, created_at  
**Audit:** Yes  
**Foreign Keys:**
- sub_account_id → sub_accounts.id
- currency_code → currencies.code

**Indexes:**
- sub_account_id
- currency_code
- completed

**Changesets:**
- `1785222139__goals.yml` — Table creation + FKs + indexes
- `1785222264__goals-audit.yml` — Audit columns
- `1785222281__goals-audit-fk.yml` — Audit FKs

---

### 11. RECURRING_TRANSACTIONS
**Purpose:** Schedule repeating transactions (subscriptions, salaries, etc.)  
**Type:** Scheduling entity  
**Columns:** id (PK, UUID), sub_account_id (FK), category_id (FK), amount, frequency, transaction_type, description, start_date, end_date, next_execution_at, active  
**Audit:** Yes  
**Foreign Keys:**
- sub_account_id → sub_accounts.id
- category_id → categories.id

**Indexes:**
- active
- next_execution_at
- sub_account_id

**Changesets:**
- `1785222139__recurring-transactions.yml` — Table creation + FKs + indexes
- `1785222264__recurring-transactions-audit.yml` — Audit columns
- `1785222281__recurring-transactions-audit-fk.yml` — Audit FKs

---

### 12. TAGS
**Purpose:** User-facing labels to categorize sub-accounts  
**Type:** Classification entity  
**Columns:** id (PK, UUID), name (UK), description  
**Audit:** Yes  
**Constraints:** Unique on name

**Changesets:**
- `1785222139__tags.yml` — Table creation
- `1785222264__tags-audit.yml` — Audit columns
- `1785222281__tags-audit-fk.yml` — Audit FKs

---

### 13. SUB_ACCOUNT_TAGS
**Purpose:** Many-to-many join between sub-accounts and tags  
**Type:** Junction table  
**Columns:** sub_account_id (FK, PK), tag_id (FK, PK)  
**Audit:** Yes  
**Foreign Keys:**
- sub_account_id → sub_accounts.id
- tag_id → tags.id

**Constraints:** Composite PK and unique constraint on (sub_account_id, tag_id)

**Changesets:**
- `1785222139__sub-account-tags.yml` — Table creation + composite PK + FKs + indexes
- `1785222264__sub-account-tags-audit.yml` — Audit columns
- `1785222281__sub-account-tags-audit-fk.yml` — Audit FKs

---

### 14. AUDIT_LOGS
**Purpose:** Generic audit trail to capture changes across entities  
**Type:** Audit entity  
**Columns:** id (PK, UUID), entity_type, entity_id, action_type, old_value (JSONB), new_value (JSONB), created_at  
**Audit:** No (audit table itself)  
**Indexes:**
- created_at (for retention queries)
- (entity_type, entity_id)

**Changesets:**
- `1785222139__audit-logs.yml` — Table creation + indexes

---

## Audit Column Standards

All tables (except ACCOUNT_BALANCES and AUDIT_LOGS) include:

```
created_at (timestamp, NOT NULL, default CURRENT_TIMESTAMP)
updated_at (timestamp, nullable)
deleted_at (timestamp, nullable)
created_by (uuid, FK→users.id)
updated_by (uuid, FK→users.id)
deleted_by (uuid, FK→users.id)
enabled (bool, NOT NULL, default true)
```

---

## Foreign Key Strategy

- **Referential Integrity:** All FKs enforce constraints
- **Delete Behavior:** NO ACTION (preserve historical data)
- **Audit Trail:** Audit columns track who created/updated/deleted records

---

## Indexing Strategy

### Indexed by FK:
All foreign key columns are indexed for query performance:
- parent_account_id (accounts → sub_accounts)
- account_type_id (sub_accounts → account_types)
- currency_code (currencies, budgets, goals)
- sub_account_id (transactions, transfers, goals, recurring_transactions)
- category_id (transactions, budgets, recurring_transactions)
- tag_id (sub_account_tags)

### Indexed by Query Patterns:
- (sub_account_id, transaction_time) — common query pattern for time-range queries
- (category_id, start_date, end_date) — budget period queries
- next_execution_at — scheduler worker queries
- active — scheduler worker queries
- completed — goal completion tracking
- created_at — audit retention policies

### Indexed for Ordering/Filtering:
- category_type (categories)
- entity_type, entity_id (audit_logs)

---

## Data Types

- **UUIDs:** All primary keys (uuidv7() default)
- **Money:** NUMERIC(19,4) for all monetary values
- **Timestamps:** TIMESTAMP for all temporal data
- **Enums:** VARCHAR for flexible extensibility

---

## Business Rules Encoded in Schema

1. **Sub-accounts are immutable containers**: Once created, account_type_id and currency_code cannot change
2. **Balances are tracked**: current_balance = opening_balance + sum(transactions)
3. **Categories are hierarchical**: Prevent cycles via application logic (not DB constraints)
4. **Transfers are atomic**: Both from/to sides succeed or fail together
5. **Recurring transactions are scheduled**: next_execution_at drives scheduler selection
6. **Tags are unique by name**: Enforce at application level if tenant-scoped
7. **Soft deletes are supported**: deleted_at and deleted_by track removal

---

## Migration Order

Changesets are designed to execute in order:

1. **Base tables (no FKs):** Currencies, Account Types
2. **Ownership tables:** Accounts
3. **Core operational tables:** Sub-accounts, Account Balances
4. **Classification:** Categories
5. **Ledger tables:** Transactions, Transfers
6. **Planning tables:** Budgets, Goals, Recurring Transactions
7. **Tagging tables:** Tags, Sub-Account Tags
8. **Audit tables:** Audit Logs
9. **Audit columns:** One per table (adds columns to existing tables)
10. **Audit FKs:** One per table (links audit columns to users)

---

## Validation Checklist

- [ ] All 39 changesets generated
- [ ] All table creation changesets populated with column definitions
- [ ] All foreign key constraints defined
- [ ] All indexes created
- [ ] All audit columns added
- [ ] All audit foreign keys created
- [ ] Database compiles without errors
- [ ] Migration executes successfully
- [ ] No missing dependencies
- [ ] Liquibase changelog master includes all changesets
- [ ] Permissions granted (GRANT statements in place)

---

## Notes for Implementation

1. **Idempotency:** All changesets use `ifNotExists: true` for table creation
2. **Backward Compatibility:** Tables are added, not modified
3. **Audit Trail:** Enable triggers or application-level capture for AUDIT_LOGS
4. **Scheduler:** Recurring transactions require a worker process to monitor next_execution_at
5. **Transactions:** Transfers require transactional consistency (consider application-level locking)
6. **JSONB Indexes:** Audit_logs JSONB columns can benefit from GIN indexes for frequent searches
7. **Data Validation:** Validate decimal_places in CURRENCIES during application startup

---

## Post-Migration Steps

1. Create stored procedures or triggers for:
   - Automatic balance updates when transactions are inserted
   - Cascade updates for goal completion tracking
   - AUDIT_LOGS entry creation on DML operations

2. Create indexes on JSONB fields if audit_logs search is high-frequency:
   ```sql
   CREATE INDEX idx_audit_logs_old_value ON audit_logs USING GIN (old_value);
   CREATE INDEX idx_audit_logs_new_value ON audit_logs USING GIN (new_value);
   ```

3. Configure retention policy for AUDIT_LOGS (old entries may be archived)

4. Create sample data for CURRENCIES and ACCOUNT_TYPES

5. Update application entity models to reflect schema

---

Generated by: DBA Skill  
RFC Reference: RFC-2-Accounts.md  
Liquibase Version: 4.x  
Database: PostgreSQL 14+
