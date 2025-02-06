-- Drop indexes
DROP INDEX IF EXISTS idx_rule_end_date;
DROP INDEX IF EXISTS idx_rule_start_date;
DROP INDEX IF EXISTS idx_plan_period_end;
DROP INDEX IF EXISTS idx_plan_period_start;
DROP INDEX IF EXISTS idx_plan_end_date;
DROP INDEX IF EXISTS idx_plan_start_date;
DROP INDEX IF EXISTS idx_tx_ord;

-- Drop tables
DROP TABLE IF EXISTS rule;
DROP TABLE IF EXISTS plan_period;
DROP TABLE IF EXISTS plan;
DROP TABLE IF EXISTS tx_cat;
DROP TABLE IF EXISTS cat;
DROP TABLE IF EXISTS tx;
DROP TABLE IF EXISTS acc;
