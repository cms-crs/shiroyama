DROP TRIGGER IF EXISTS set_updated_at ON team_members;

DROP FUNCTION IF EXISTS update_updated_at_column;

ALTER TABLE team_members
    DROP COLUMN IF EXISTS created_at,
    DROP COLUMN IF EXISTS updated_at;