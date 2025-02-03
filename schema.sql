CREATE TABLE acc (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    xid TEXT NOT NULL UNIQUE,
);

CREATE TABLE tx (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    xid TEXT NOT NULL,
    date TEXT NOT NULL,
    orig_date TEXT NOT NULL,
    description TEXT NOT NULL,
    amount REAL NOT NULL,
    orig_amount REAL NOT NULL,
    acc_id INTEGER NOT NULL,
    FOREIGN KEY (acc_id) REFERENCES acc(id),
    UNIQUE (acc_id, xid)
);

CREATE INDEX idx_tx_date ON tx(date);

CREATE TABLE cat (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    pattern TEXT NOT NULL,
);

CREATE TABLE tx_cat (
    tx_id INTEGER NOT NULL,
    cat_id INTEGER NOT NULL,
    FOREIGN KEY (tx_id) REFERENCES tx(id),
    FOREIGN KEY (cat_id) REFERENCES cat(id),
    PRIMARY KEY (tx_id, cat_id)
);

CREATE TABLE exp (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    period TEXT NOT NULL,
    pattern TEXT NOT NULL,
);

CREATE TABLE tx_exp (
    tx_id INTEGER NOT NULL,
    exp_id INTEGER NOT NULL,
    FOREIGN KEY (tx_id) REFERENCES tx(id),
    FOREIGN KEY (exp_id) REFERENCES exp(id),
    PRIMARY KEY (tx_id, exp_id)
);

CREATE TABLE plan (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    start_date TEXT NOT NULL,
    end_date TEXT NOT NULL,
);

CREATE INDEX idx_plan_start_date ON plan(start_date);
CREATE INDEX idx_plan_end_date ON plan(end_date);

CREATE TABLE plan_exp (
    plan_id INTEGER NOT NULL,
    exp_id INTEGER NOT NULL,
    amount REAL NOT NULL,
    FOREIGN KEY (plan_id) REFERENCES plan(id),
    FOREIGN KEY (exp_id) REFERENCES exp(id),
    PRIMARY KEY (plan_id, exp_id)
);
