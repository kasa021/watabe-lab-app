-- チェックインログテーブルの削除
DROP INDEX IF EXISTS idx_check_in_logs_current;
DROP INDEX IF EXISTS idx_check_in_logs_check_out_at;
DROP INDEX IF EXISTS idx_check_in_logs_check_in_at;
DROP INDEX IF EXISTS idx_check_in_logs_user_id;
DROP TABLE IF EXISTS check_in_logs;