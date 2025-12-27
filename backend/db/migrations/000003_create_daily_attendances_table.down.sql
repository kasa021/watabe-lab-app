-- 日次出席記録テーブルの削除
DROP INDEX IF EXISTS idx_daily_attendances_user_date;
DROP INDEX IF EXISTS idx_daily_attendances_date;
DROP INDEX IF EXISTS idx_daily_attendances_user_id;
DROP TABLE IF EXISTS daily_attendances;