-- 日次出席記録テーブルの作成
CREATE TABLE IF NOT EXISTS daily_attendances (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    attendance_date DATE NOT NULL,
    total_duration_minutes INTEGER NOT NULL DEFAULT 0,
    check_in_count INTEGER NOT NULL DEFAULT 0,
    first_check_in_at TIME,
    last_check_out_at TIME,
    points INTEGER NOT NULL DEFAULT 0,
    is_holiday BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, attendance_date)
);
-- インデックスの作成
CREATE INDEX idx_daily_attendances_user_id ON daily_attendances(user_id);
CREATE INDEX idx_daily_attendances_date ON daily_attendances(attendance_date);
CREATE INDEX idx_daily_attendances_user_date ON daily_attendances(user_id, attendance_date);
-- コメント
COMMENT ON TABLE daily_attendances IS '日次出席記録';
COMMENT ON COLUMN daily_attendances.attendance_date IS '出席日';
COMMENT ON COLUMN daily_attendances.total_duration_minutes IS '合計滞在時間（分）';
COMMENT ON COLUMN daily_attendances.check_in_count IS 'チェックイン回数';
COMMENT ON COLUMN daily_attendances.points IS '獲得ポイント（1 or 0）';