CREATE TABLE "mindwell"."badges" (
    "id" Serial NOT NULL,
    "code" Text NOT NULL,
    "level" Integer DEFAULT 1 NOT NULL,
    "title" Text NOT NULL,
    "description" Text NOT NULL,
    "icon" Text NOT NULL,
    "user_count" Integer DEFAULT 0 NOT NULL,
    CONSTRAINT "unique_badge_id" PRIMARY KEY ( "id" ),
    CONSTRAINT "unique_badge_code_level" UNIQUE ( "code", "level" )
);

INSERT INTO mindwell.badges(code, title, description, icon)
VALUES('test1', 'Тест (1)', 'Это тестовый значок (1). Вы восхитительны!', 'test.webp');

INSERT INTO mindwell.badges(code, title, description, icon)
VALUES('test2', 'Тест (2)', 'Это тестовый значок (2). Вы восхитительны!', 'test.webp');

INSERT INTO mindwell.badges(code, title, description, icon)
VALUES('test3', 'Тест (3)', 'Это тестовый значок (3). Вы восхитительны!', 'test.webp');

INSERT INTO mindwell.badges(code, title, description, icon)
VALUES('test4', 'Тест (4)', 'Это тестовый значок (4). Вы восхитительны!', 'test.webp');

INSERT INTO mindwell.badges(code, title, description, icon)
VALUES('test5', 'Тест (5)', 'Это тестовый значок (5). Вы восхитительны!', 'test.webp');

CREATE TABLE "user_badges" (
    "user_id" Integer NOT NULL,
    "badge_id" Integer NOT NULL,
    "given_at" Timestamp With Time Zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT "unique_user_badge_badge_id" FOREIGN KEY ("badge_id") REFERENCES mindwell.badges( "id" ),
    CONSTRAINT "unique_user_badge_user_id" FOREIGN KEY ( "user_id" ) REFERENCES "mindwell"."users"("id"),
    CONSTRAINT "unique_user_badge" UNIQUE ( "user_id", "badge_id" )
);

ALTER TABLE users
ADD COLUMN "badges_count" Integer DEFAULT 0 NOT NULL;

CREATE OR REPLACE FUNCTION mindwell.inc_badges() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE mindwell.users
        SET badges_count = badges_count + 1
        WHERE id = NEW.user_id;

        UPDATE mindwell.badges
        SET user_count = user_count + 1
        WHERE id = NEW.badge_id;

        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION mindwell.dec_badges() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE mindwell.users
        SET badges_count = badges_count - 1
        WHERE id = OLD.user_id;

        UPDATE mindwell.badges
        SET user_count = user_count - 1
        WHERE id = OLD.badge_id;

        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER cnt_badges_inc
    AFTER INSERT ON mindwell.user_badges
    FOR EACH ROW EXECUTE PROCEDURE mindwell.inc_badges();

CREATE TRIGGER cnt_badges_dec
    AFTER DELETE ON mindwell.user_badges
    FOR EACH ROW EXECUTE PROCEDURE mindwell.dec_badges();

CREATE OR REPLACE FUNCTION mindwell.notify_user_badges() RETURNS TRIGGER AS $$
    BEGIN
            PERFORM pg_notify('user_badges', json_build_object(
                'user_id', NEW.user_id,
                'badge_id', NEW.badge_id
            )::text);
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER ntf_badges
    AFTER INSERT ON mindwell.user_badges
    FOR EACH ROW EXECUTE FUNCTION mindwell.notify_user_badges();

INSERT INTO "mindwell"."notification_type" VALUES(13, 'badge');

ALTER TABLE users
ADD COLUMN "email_badges" Boolean NOT NULL DEFAULT FALSE;

ALTER TABLE users
ADD COLUMN "telegram_badges" Boolean NOT NULL DEFAULT TRUE;
