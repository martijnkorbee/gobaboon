ALTER TABLE sessions DROP INDEX sessions_expiry_idx;
DROP TABLE IF EXISTS sessions CASCADE;
