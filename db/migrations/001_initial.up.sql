CREATE TABLE acc (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT current_timestamp,
    updated_at TEXT NOT NULL DEFAULT current_timestamp,
    name TEXT NOT NULL UNIQUE,
    xid TEXT NOT NULL UNIQUE,
    kind TEXT NOT NULL CHECK (kind IN ('bank', 'cc')),
    is_active INTEGER NOT NULL CHECK (is_active IN (0, 1))
);

CREATE TABLE tx (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT current_timestamp,
    updated_at TEXT NOT NULL DEFAULT current_timestamp,
    xid TEXT NOT NULL,
    date TEXT NOT NULL,
    desc TEXT NOT NULL,
    amount REAL NOT NULL,
    orig_date TEXT,
    orig_desc TEXT,
    orig_amount REAL,
    acc_id INTEGER NOT NULL,
    cat_id INTEGER, 
    ord TEXT NOT NULL,
    FOREIGN KEY (acc_id) REFERENCES acc(id) ON DELETE CASCADE,
    FOREIGN KEY (cat_id) REFERENCES cat(id) ON DELETE SET NULL,
    UNIQUE (acc_id, xid)
);

CREATE INDEX idx_tx_ord ON tx(ord);

CREATE TABLE cat (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT current_timestamp,
    updated_at TEXT NOT NULL DEFAULT current_timestamp,
    name TEXT NOT NULL UNIQUE,
    kind TEXT NOT NULL CHECK (kind IN ('income', 'expense', 'transfer')),
    is_active INTEGER NOT NULL CHECK (is_active IN (0, 1))
);

CREATE TABLE plan (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT current_timestamp,
    updated_at TEXT NOT NULL DEFAULT current_timestamp,
    start_date TEXT,
    end_date TEXT,
    cat_id INTEGER NOT NULL,
    amount_expr TEXT NOT NULL,
    period TEXT NOT NULL CHECK (period IN ('month', 'year')),
    FOREIGN KEY (cat_id) REFERENCES cat(id) ON DELETE CASCADE
);

CREATE INDEX idx_plan_start_date ON plan(start_date);
CREATE INDEX idx_plan_end_date ON plan(end_date);

CREATE TABLE plan_period (
    created_at TEXT NOT NULL DEFAULT current_timestamp,
    updated_at TEXT NOT NULL DEFAULT current_timestamp,
    plan_id INTEGER NOT NULL,
    period_start TEXT NOT NULL,
    period_end TEXT NOT NULL,
    amount REAL NOT NULL,
    FOREIGN KEY (plan_id) REFERENCES plan(id) ON DELETE CASCADE,
    PRIMARY KEY (plan_id, period_start, period_end)
);

CREATE INDEX idx_plan_period_start ON plan_period(period_start);
CREATE INDEX idx_plan_period_end ON plan_period(period_end);

CREATE TABLE rule (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL DEFAULT current_timestamp,
    updated_at TEXT NOT NULL DEFAULT current_timestamp,
    start_date TEXT,
    end_date TEXT,
    test_expr TEXT NOT NULL,
    cat_id INTEGER,
    amount_expr TEXT,
    desc_expr TEXT,
    date_expr TEXT,
    ord INTEGER NOT NULL,
    FOREIGN KEY (cat_id) REFERENCES cat(id) ON DELETE CASCADE
);

CREATE INDEX idx_rule_start_date ON rule(start_date);
CREATE INDEX idx_rule_end_date ON rule(end_date);
