-- +goose Up
CREATE TABLE auth_codes (
  id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  code INTEGER NOT NULL DEFAULT (floor(random() * 900000 + 100000)),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (id, user_id)
);

-- CREATE OR REPLACE FUNCTION trigger_set_updated_at()
-- RETURNS TRIGGER AS $$
-- BEGIN
--   NEW.updated_at = NOW();
--   RETURN NEW;
-- END;
-- $$ LANGUAGE plpgsql;

CREATE TRIGGER set_updated_at
BEFORE UPDATE ON auth_codes
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_updated_at();



-- +goose Down
DROP TRIGGER IF EXISTS set_updated_at ON auth_codes;
DROP TABLE auth_codes;