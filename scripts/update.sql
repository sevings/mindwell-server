ALTER TABLE user_log
ADD COLUMN uid2 Bigint NOT NULL DEFAULT 0;

ALTER TABLE user_log
ALTER COLUMN uid2 DROP DEFAULT;

DROP VIEW user_log_view;

CREATE VIEW "mindwell"."user_log_view" AS
SELECT name, ip, to_hex(device) AS device, to_hex(app) AS app, to_hex(uid) AS uid, to_hex(uid2) AS uid2,
       to_char(at, 'YYYY.MM.DD HH24:MI:SS') AS at, first AS f
FROM mindwell.user_log
ORDER BY at DESC;
