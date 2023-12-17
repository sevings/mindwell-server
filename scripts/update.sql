ALTER TABLE "mindwell"."users"
ADD COLUMN "send_wishes" Boolean NOT NULL DEFAULT TRUE;

CREATE TABLE "mindwell"."wish_states" (
    "id" Integer UNIQUE NOT NULL,
    "state" Text NOT NULL
);

INSERT INTO "mindwell"."wish_states"(id, state) VALUES (0, 'new');
INSERT INTO "mindwell"."wish_states"(id, state) VALUES (1, 'sent');
INSERT INTO "mindwell"."wish_states"(id, state) VALUES (2, 'declined');
INSERT INTO "mindwell"."wish_states"(id, state) VALUES (3, 'complained');
INSERT INTO "mindwell"."wish_states"(id, state) VALUES (4, 'thanked');

CREATE TABLE "mindwell"."wishes" (
    "id" Serial NOT NULL,
    "from_id" Integer NOT NULL,
    "to_id" Integer NOT NULL,
    "content" Text DEFAULT '' NOT NULL,
    "state" Integer DEFAULT 0 NOT NULL,
    "created_at" Timestamp With Time Zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT "unique_wish_id" PRIMARY KEY( "id" ),
    CONSTRAINT "wish_sender" FOREIGN KEY("from_id") REFERENCES "mindwell"."users"("id"),
    CONSTRAINT "wish_receiver" FOREIGN KEY("to_id") REFERENCES "mindwell"."users"("id"),
    CONSTRAINT "enum_wish_state" FOREIGN KEY("state") REFERENCES "mindwell"."wish_states"("id")
);

CREATE INDEX "index_wish_id" ON "mindwell"."wishes" USING btree( "id" );
CREATE INDEX "index_wish_from_id" ON "mindwell"."wishes" USING btree( "from_id" );
CREATE INDEX "index_wish_to_id" ON "mindwell"."wishes" USING btree( "to_id" );

INSERT INTO "mindwell"."complain_type" VALUES(5, 'wish');

INSERT INTO "mindwell"."notification_type" VALUES(10, 'wish_created');
INSERT INTO "mindwell"."notification_type" VALUES(11, 'wish_received');
