DROP VIEW user_log_view;
DROP TABLE user_log;

CREATE TABLE "mindwell"."user_log" (
    "name" Text NOT NULL,
    "user_agent" Text NOT NULL,
    "ip" Inet NOT NULL,
    "device" Integer NOT NULL,
    "app" Bigint NOT NULL,
    "uid" Integer NOT NULL,
    "at" Timestamp With Time Zone NOT NULL,
    "first" Boolean NOT NULL
);

CREATE VIEW "mindwell"."user_log_view" AS
SELECT name, ip, to_hex(device) AS device, to_hex(app) AS app, to_hex(uid) AS uid, to_char(at, 'YYYY.MM.DD HH24:MI:SS') AS at, first
FROM user_log
ORDER BY at DESC;

CREATE INDEX "index_user_log_at" ON "mindwell"."user_log" USING btree( "at" );

CREATE INDEX "index_user_log_name" ON "mindwell"."user_log" USING btree( "name" );
