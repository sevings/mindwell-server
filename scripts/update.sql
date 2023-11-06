ALTER TABLE "mindwell"."adm"
ADD COLUMN "phone" Text NOT NULL DEFAULT '';

ALTER TABLE "mindwell"."adm"
ALTER COLUMN "phone" DROP DEFAULT;

CREATE TABLE "mindwell"."user_chat_privacy" (
    "id" Integer UNIQUE NOT NULL,
    "type" Text NOT NULL );

INSERT INTO "mindwell"."user_chat_privacy" VALUES(0, 'invited');
INSERT INTO "mindwell"."user_chat_privacy" VALUES(1, 'followers');
INSERT INTO "mindwell"."user_chat_privacy" VALUES(2, 'friends');
INSERT INTO "mindwell"."user_chat_privacy" VALUES(3, 'me');

ALTER TABLE "mindwell"."users"
ADD COLUMN "chat_privacy" Integer DEFAULT 0 NOT NULL;

ALTER TABLE "mindwell"."users"
ADD CONSTRAINT "enum_user_chat_privacy" FOREIGN KEY("chat_privacy") REFERENCES "mindwell"."user_chat_privacy"("id");

DROP TRIGGER can_send_ins ON mindwell.talkers;
DROP TRIGGER can_send_invited ON mindwell.users;
DROP TRIGGER can_send_related ON mindwell.relations;
DROP TRIGGER can_send_not_related ON mindwell.relations;
DROP TRIGGER can_send_new_msg ON mindwell.talkers;

DROP FUNCTION mindwell.set_can_send_ins;
DROP FUNCTION mindwell.set_can_send_invited;
DROP FUNCTION mindwell.set_can_send_related;
DROP FUNCTION mindwell.set_can_send_not_related;
DROP FUNCTION mindwell.set_can_send_new_msg;

DROP FUNCTION mindwell.is_invited;
DROP FUNCTION mindwell.is_partner_ignoring;

ALTER TABLE "mindwell"."talkers"
DROP COLUMN "can_send";
