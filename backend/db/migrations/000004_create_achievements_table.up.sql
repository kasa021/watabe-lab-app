-- 称号マスタテーブルの作成
CREATE TABLE IF NOT EXISTS achievements (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    icon_url VARCHAR(255),
    category VARCHAR(50),
    condition_type VARCHAR(50) NOT NULL,
    condition_value JSONB,
    points_reward INTEGER DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- インデックスの作成
CREATE INDEX idx_achievements_code ON achievements(code);
CREATE INDEX idx_achievements_category ON achievements(category);
-- コメント
COMMENT ON TABLE achievements IS '称号マスタ';
COMMENT ON COLUMN achievements.code IS '称号コード（例: early_bird）';
COMMENT ON COLUMN achievements.name IS '称号名（例: 早起き博士）';
COMMENT ON COLUMN achievements.category IS 'カテゴリ（attendance/time/streak/special）';
COMMENT ON COLUMN achievements.condition_type IS '条件タイプ';
COMMENT ON COLUMN achievements.condition_value IS '条件の詳細（JSON形式）';