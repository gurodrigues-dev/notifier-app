BEGIN;

ALTER TABLE channels
DROP CONSTRAINT IF EXISTS unique_platform_target_group;

COMMIT;
