-- ユーザーテーブルの作成
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    email VARCHAR(255),
    role VARCHAR(20) NOT NULL DEFAULT 'student',
    is_presence_public BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_login_at TIMESTAMP,
    is_active BOOLEAN NOT NULL DEFAULT true
);
-- インデックスの作成
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_role ON users(role);
-- コメント
COMMENT ON TABLE users IS 'ユーザー情報';
COMMENT ON COLUMN users.username IS 'LDAPユーザー名';
COMMENT ON COLUMN users.display_name IS '表示名';
COMMENT ON COLUMN users.role IS '権限（student/teacher/admin）';
COMMENT ON COLUMN users.is_presence_public IS '在室状態の公開設定';