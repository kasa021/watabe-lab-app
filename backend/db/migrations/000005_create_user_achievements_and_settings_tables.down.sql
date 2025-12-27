-- システム設定テーブルの削除
DROP TABLE IF EXISTS settings;
-- ユーザー称号テーブルの削除
DROP INDEX IF EXISTS idx_user_achievements_achievement_id;
DROP INDEX IF EXISTS idx_user_achievements_user_id;
DROP TABLE IF EXISTS user_achievements;