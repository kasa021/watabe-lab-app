-- 称号マスタテーブルの削除
DROP INDEX IF EXISTS idx_achievements_category;
DROP INDEX IF EXISTS idx_achievements_code;
DROP TABLE IF EXISTS achievements;