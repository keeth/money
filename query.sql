-- name: GetAccByXid :one
SELECT * FROM acc WHERE xid = ? LIMIT 1;

-- name: GetAccs :many
SELECT * FROM acc;

-- name: CreateAcc :exec
INSERT INTO acc (
    xid, 
    name, 
    is_active
) VALUES (
    ?, 
    ?, 
    1
) 
RETURNING id;

-- name: GetTxByAccAndXid :one
SELECT * FROM tx WHERE acc_id = ? AND xid = ? LIMIT 1;

-- name: GetTxs :many
SELECT sqlc.embed(tx), sqlc.embed(acc), sqlc.embed(cat)
FROM tx
JOIN acc ON tx.acc_id = acc.id
LEFT JOIN tx_cat ON tx.id = tx_cat.tx_id
LEFT JOIN cat ON tx_cat.cat_id = cat.id
WHERE ord < ?
ORDER BY ord DESC
LIMIT ?;

-- name: CreateTx :exec
INSERT INTO tx (
    xid, 
    date, 
    orig_date, 
    desc, 
    orig_desc, 
    amount, 
    orig_amount, 
    acc_id,
    ord
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?
) RETURNING id;

-- name: UpdateTx :exec
UPDATE tx SET
    date = ?,
    desc = ?,
    amount = ?,
    ord = ?
WHERE id = ?;

-- name: CreateOrUpdateTx :exec
INSERT INTO tx (
    xid, 
    date, 
    orig_date, 
    desc, 
    orig_desc, 
    amount, 
    orig_amount, 
    acc_id,
    ord
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?
) 
ON CONFLICT (xid, acc_id) DO UPDATE SET
    date = EXCLUDED.date,
    desc = EXCLUDED.desc,
    amount = EXCLUDED.amount,
    ord = EXCLUDED.ord
RETURNING id;

-- name: GetCats :many
SELECT * FROM cat WHERE is_active = 1 ORDER BY name;

-- name: CreateCat :exec
INSERT INTO cat (name, kind, is_active) VALUES (?, ?, 1);

-- name: UpdateCat :exec
UPDATE cat SET name = ? WHERE id = ?;

-- name: DeactivateCat :exec
UPDATE cat SET is_active = 0 WHERE id = ?;

-- name: GetPlans :many
SELECT * 
FROM plan 
WHERE (start_date IS NULL OR start_date >= ?) 
    AND (end_date IS NULL OR end_date < ?)
ORDER BY start_date DESC;

-- name: CreatePlan :exec
INSERT INTO plan (
    start_date, 
    end_date, 
    amount_expr, 
    cat_id,
    period
) VALUES (
    ?,
    ?,
    ?,
    ?,
    ?
);

-- name: UpdatePlan :exec
UPDATE plan SET start_date = ?, end_date = ?, amount_expr = ? WHERE id = ?;

-- name: GetPlanPeriodsByPlan :many
SELECT * 
FROM plan_period 
WHERE plan_id = ? 
    AND period_start >= ? 
    AND period_end < ?;

-- name: GetPlanPeriods :many
SELECT sqlc.embed(plan_period), sqlc.embed(plan)
FROM plan_period 
JOIN plan ON plan_period.plan_id = plan.id
WHERE period_start >= ? AND period_end < ?;

-- name: CreatePlanPeriod :exec
INSERT INTO plan_period (
    plan_id, 
    period_start, 
    period_end, 
    amount
) VALUES (
    ?, ?, ?, ?
);

-- name: UpdatePlanPeriod :exec
UPDATE plan_period SET amount = ? WHERE id = ?;

-- name: DeletePlanPeriods :exec
DELETE FROM plan_period 
WHERE plan_id = ?
    AND period_start >= ?
    AND period_end < ?;

-- name: GetRules :many
SELECT * FROM rule 
WHERE (start_date IS NULL OR start_date >= ?) 
    AND (end_date IS NULL OR end_date < ?)
ORDER BY ord;

-- name: CreateRule :exec
INSERT INTO rule (
    start_date, 
    end_date, 
    test_expr, 
    cat_id,
    amount_expr,
    desc_expr,
    date_expr,
    ord
) VALUES (
    ?, 
    ?, 
    ?, 
    ?, 
    ?,
    ?,
    ?,
    ?
);

-- name: UpdateRule :exec
UPDATE rule SET 
    start_date = ?,
    end_date = ?,
    test_expr = ?,
    cat_id = ?,
    amount_expr = ?,
    desc_expr = ?,
    date_expr = ?,
    ord = ?
WHERE id = ?;

-- name: DeleteRule :exec
DELETE FROM rule WHERE id = ?;