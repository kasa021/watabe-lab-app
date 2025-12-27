-- チェックインログテーブルの作成
CREATE TABLE IF NOT EXISTS check_in_logs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    check_in_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    check_out_at TIMESTAMP,
    duration_minutes INTEGER,
    check_in_method VARCHAR(20),
    wifi_ssid VARCHAR(100),
    gps_latitude DECIMAL(10, 8),
    gps_longitude DECIMAL(11, 8),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- インデックスの作成
CREATE INDEX idx_check_in_logs_user_id ON check_in_logs(user_id);
CREATE INDEX idx_check_in_logs_check_in_at ON check_in_logs(check_in_at);
CREATE INDEX idx_check_in_logs_check_out_at ON check_in_logs(check_out_at);
CREATE INDEX idx_check_in_logs_current ON check_in_logs(user_id, check_out_at)
WHERE check_out_at IS NULL;
-- コメント
COMMENT ON TABLE check_in_logs IS 'チェックインログ';
COMMENT ON COLUMN check_in_logs.user_id IS 'ユーザーID';
COMMENT ON COLUMN check_in_logs.check_in_at IS 'チェックイン日時';
COMMENT ON COLUMN check_in_logs.check_out_at IS 'チェックアウト日時（NULL=在室中）';
COMMENT ON COLUMN check_in_logs.duration_minutes IS '滞在時間（分）';