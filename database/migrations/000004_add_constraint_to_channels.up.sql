BEGIN;

ALTER TABLE channels
ADD CONSTRAINT unique_platform_target_group UNIQUE (platform, target_id, "group");

COMMIT;
