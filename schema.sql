CREATE TABLE acc (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    xid TEXT NOT NULL UNIQUE,
    is_active INTEGER NOT NULL CHECK (is_active IN (0, 1))
);

CREATE TABLE tx (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    xid TEXT NOT NULL,
    date TEXT NOT NULL,
    orig_date TEXT NOT NULL,
    desc TEXT NOT NULL,
    orig_desc TEXT NOT NULL,
    amount REAL NOT NULL,
    orig_amount REAL NOT NULL,
    acc_id INTEGER NOT NULL,
    ord TEXT NOT NULL,
    FOREIGN KEY (acc_id) REFERENCES acc(id) ON DELETE CASCADE,
    UNIQUE (acc_id, xid)
);

CREATE INDEX idx_tx_ord ON tx(ord);

CREATE TABLE cat (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    kind TEXT NOT NULL CHECK (kind IN ('income', 'expense', 'transfer')),
    is_active INTEGER NOT NULL CHECK (is_active IN (0, 1))
);

CREATE TABLE tx_cat (
    tx_id INTEGER NOT NULL,
    cat_id INTEGER NOT NULL,
    FOREIGN KEY (tx_id) REFERENCES tx(id) ON DELETE CASCADE,
    FOREIGN KEY (cat_id) REFERENCES cat(id) ON DELETE CASCADE,
    PRIMARY KEY (tx_id, cat_id)
);

CREATE TABLE plan (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
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
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    plan_id INTEGER NOT NULL,
    period_start TEXT NOT NULL,
    period_end TEXT NOT NULL,
    amount REAL NOT NULL,
    FOREIGN KEY (plan_id) REFERENCES plan(id) ON DELETE CASCADE,
    PRIMARY KEY (plan_id, period_start, period_end)
);

CREATE INDEX idx_plan_eval_period_start ON plan_eval(period_start);
CREATE INDEX idx_plan_eval_period_end ON plan_eval(period_end);

CREATE TABLE rule (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
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
