CREATE SCHEMA "yummy";

-- CREATE TABLE "gender" ---------------------------------
CREATE TABLE "yummy"."gender" (
    "id" Integer UNIQUE NOT NULL,
    "type" Text NOT NULL );

INSERT INTO "yummy"."gender" VALUES(0, 'not set');
INSERT INTO "yummy"."gender" VALUES(1, 'male');
INSERT INTO "yummy"."gender" VALUES(2, 'female');
-- -------------------------------------------------------------



-- CREATE TABLE "user_privacy" ---------------------------------
CREATE TABLE "yummy"."user_privacy" (
    "id" Integer UNIQUE NOT NULL,
    "type" Text NOT NULL );

INSERT INTO "yummy"."user_privacy" VALUES(0, 'all');
INSERT INTO "yummy"."user_privacy" VALUES(1, 'registered');
INSERT INTO "yummy"."user_privacy" VALUES(2, 'followers');
-- -------------------------------------------------------------



-- CREATE TYPE "font_family" -----------------------------------
CREATE TABLE "yummy"."font_family" (
    "id" Integer UNIQUE NOT NULL,
    "type" Text NOT NULL );

INSERT INTO "yummy"."font_family" VALUES(0, 'Arial');
-- -------------------------------------------------------------



-- CREATE TYPE "alignment" -------------------------------------
CREATE TABLE "yummy"."alignment" (
    "id" Integer UNIQUE NOT NULL,
    "type" Text NOT NULL );

INSERT INTO "yummy"."alignment" VALUES(0, 'left');
INSERT INTO "yummy"."alignment" VALUES(1, 'right');
INSERT INTO "yummy"."alignment" VALUES(2, 'center');
INSERT INTO "yummy"."alignment" VALUES(3, 'justify');
-- -------------------------------------------------------------



-- CREATE TABLE "users" ----------------------------------------
CREATE TABLE "yummy"."users" ( 
	"id" Serial NOT NULL,
	"name" Text NOT NULL,
	"show_name" Text DEFAULT '' NOT NULL,
	"password_hash" Bytea NOT NULL,
	"gender" Integer DEFAULT 0 NOT NULL,
	"is_daylog" Boolean DEFAULT false NOT NULL,
	"privacy" Integer DEFAULT 0 NOT NULL,
	"title" Text DEFAULT '' NOT NULL,
	"last_seen_at" Timestamp With Time Zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"karma" Real DEFAULT 0 NOT NULL,
	"created_at" Timestamp With Time Zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "invited_by" Integer NOT NULL,
	"birthday" Date,
	"css" Text DEFAULT '' NOT NULL,
	"entries_count" Integer DEFAULT 0 NOT NULL,
	"followings_count" Integer DEFAULT 0 NOT NULL,
	"followers_count" Integer DEFAULT 0 NOT NULL,
	"comments_count" Integer DEFAULT 0 NOT NULL,
	"ignored_count" Integer DEFAULT 0 NOT NULL,
	"invited_count" Integer DEFAULT 0 NOT NULL,
	"favorites_count" Integer DEFAULT 0 NOT NULL,
	"tags_count" Integer DEFAULT 0 NOT NULL,
	"country" Text DEFAULT '' NOT NULL,
	"city" Text DEFAULT '' NOT NULL,
	"email" Text NOT NULL,
	"verified" Boolean DEFAULT false NOT NULL,
    "api_key" Text NOT NULL,
    "valid_thru" Timestamp With Time Zone DEFAULT CURRENT_TIMESTAMP + interval '6 months' NOT NULL,
	"avatar" Text DEFAULT '' NOT NULL,
	"font_family" Integer DEFAULT 0 NOT NULL,
	"font_size" SmallInt DEFAULT 100 NOT NULL,
	"text_alignment" Integer DEFAULT 0 NOT NULL,
	"name_color" Character( 7 ) NOT NULL,
	"avatar_color" Character( 7 )  NOT NULL,
	"text_color" Character( 7 ) DEFAULT '#000000' NOT NULL,
	"background_color" Character( 7 ) DEFAULT '#ffffff' NOT NULL,
	CONSTRAINT "unique_user_id" PRIMARY KEY( "id" ),
    CONSTRAINT "enum_user_gender" FOREIGN KEY("gender") REFERENCES "yummy"."gender"("id"),
    CONSTRAINT "enum_user_privacy" FOREIGN KEY("privacy") REFERENCES "yummy"."user_privacy"("id"),
    CONSTRAINT "enum_user_alignment" FOREIGN KEY("text_alignment") REFERENCES "yummy"."alignment"("id"),
    CONSTRAINT "enum_user_font_family" FOREIGN KEY("font_family") REFERENCES "yummy"."font_family"("id") );
 ;
-- -------------------------------------------------------------

-- CREATE INDEX "index_user_id" --------------------------------
CREATE UNIQUE INDEX "index_user_id" ON "yummy"."users" USING btree( "id" );
-- -------------------------------------------------------------

-- CREATE INDEX "index_user_name" ------------------------------
CREATE UNIQUE INDEX "index_user_name" ON "yummy"."users" USING btree( lower("name") );
-- -------------------------------------------------------------

-- CREATE INDEX "index_user_email" -----------------------------
CREATE UNIQUE INDEX "index_user_email" ON "yummy"."users" USING btree( lower("email") );
-- -------------------------------------------------------------

-- CREATE INDEX "index_token_user" -----------------------------
CREATE UNIQUE INDEX "index_user_key" ON "yummy"."users" USING btree( "api_key" );
-- -------------------------------------------------------------

-- CREATE INDEX "index_invited_by" -----------------------------
CREATE INDEX "index_invited_by" ON "yummy"."users" USING btree( "invited_by" );
-- -------------------------------------------------------------

CREATE OR REPLACE FUNCTION yummy.count_invited() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE yummy.users
        SET invited_count = invited_count + 1 
        WHERE id = NEW.invited_by;

        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER cnt_invited
    AFTER INSERT ON yummy.users
    FOR EACH ROW EXECUTE PROCEDURE yummy.count_invited();

INSERT INTO yummy.users
    (name, show_name, email, password_hash, api_key, invited_by, name_color, avatar_color)
    VALUES('HaveANiceDay', 'Хорошего дня!', '', '', '', 1, '#000000', '#ffffff');



CREATE VIEW yummy.short_users AS
SELECT id, name, show_name,
    now() - last_seen_at < interval '15 minutes' AS is_online,
    name_color, avatar_color, avatar
FROM yummy.users;



CREATE VIEW yummy.long_users AS
SELECT users.id,
    users.name,
    users.show_name,
    users.password_hash,
    gender.type AS gender,
    users.is_daylog,
    user_privacy.type AS privacy,
    users.title,
    users.last_seen_at,
    users.karma,
    users.created_at,
    users.invited_by,
    users.birthday,
    users.css,
    users.entries_count,
    users.followings_count,
    users.followers_count,
    users.comments_count,
    users.ignored_count,
    users.invited_count,
    users.favorites_count,
    users.tags_count,
    users.country,
    users.city,
    users.email,
    users.verified,
    users.api_key,
    users.valid_thru,
    users.avatar,
    font_family.type AS font_family,
    users.font_size,
    alignment.type AS text_alignment,
    users.name_color,
    users.avatar_color,
    users.text_color,
    users.background_color,
    now() - last_seen_at < interval '15 minutes' AS is_online,
    extract(year from age(birthday))::integer as "age",
    short_users.id AS invited_by_id,
    short_users.name AS invited_by_name,
    short_users.show_name AS invited_by_show_name,
    short_users.is_online AS invited_by_is_online,
    short_users.name_color AS invited_by_name_color,
    short_users.avatar_color AS invited_by_avatar_color,
    short_users.avatar AS invited_by_avatar
FROM yummy.users, yummy.short_users,
    yummy.gender, yummy.user_privacy, yummy.font_family, yummy.alignment
WHERE users.invited_by = short_users.id
    AND users.gender = gender.id
    AND users.privacy = user_privacy.id
    AND users.font_family = font_family.id
    AND users.text_alignment = alignment.id;

    

-- CREATE TABLE "invite_words" ---------------------------------
CREATE TABLE "yummy"."invite_words" (
    "id" Serial NOT NULL,
    "word" Text NOT NULL,
	CONSTRAINT "unique_word_id" PRIMARY KEY( "id" ),
	CONSTRAINT "unique_word" UNIQUE( "word" ) );
;
-- -------------------------------------------------------------

-- CREATE INDEX "index_invite_word_id" -------------------------
CREATE UNIQUE INDEX "index_invite_word_id" ON "yummy"."invite_words" USING btree( "id" );
-- -------------------------------------------------------------

-- CREATE INDEX "index_invite_word" ----------------------------
CREATE UNIQUE INDEX "index_invite_word" ON "yummy"."invite_words" USING btree( "word" );
-- -------------------------------------------------------------

INSERT INTO yummy.invite_words ("word") VALUES('acknown');
INSERT INTO yummy.invite_words ("word") VALUES('aery');
INSERT INTO yummy.invite_words ("word") VALUES('affectioned');
INSERT INTO yummy.invite_words ("word") VALUES('agnize');
INSERT INTO yummy.invite_words ("word") VALUES('ambition');
INSERT INTO yummy.invite_words ("word") VALUES('amerce');
INSERT INTO yummy.invite_words ("word") VALUES('anters');
INSERT INTO yummy.invite_words ("word") VALUES('argal');
INSERT INTO yummy.invite_words ("word") VALUES('arrant');
INSERT INTO yummy.invite_words ("word") VALUES('arras');
INSERT INTO yummy.invite_words ("word") VALUES('asquint');
INSERT INTO yummy.invite_words ("word") VALUES('atomies');
INSERT INTO yummy.invite_words ("word") VALUES('augurers');
INSERT INTO yummy.invite_words ("word") VALUES('bastinado');
INSERT INTO yummy.invite_words ("word") VALUES('batten');
INSERT INTO yummy.invite_words ("word") VALUES('bawbling');
INSERT INTO yummy.invite_words ("word") VALUES('bawcock');
INSERT INTO yummy.invite_words ("word") VALUES('bawd');
INSERT INTO yummy.invite_words ("word") VALUES('behoveful');
INSERT INTO yummy.invite_words ("word") VALUES('beldams');
INSERT INTO yummy.invite_words ("word") VALUES('belike');
INSERT INTO yummy.invite_words ("word") VALUES('berattle');
INSERT INTO yummy.invite_words ("word") VALUES('beshrew');
INSERT INTO yummy.invite_words ("word") VALUES('betid');
INSERT INTO yummy.invite_words ("word") VALUES('betimes');
INSERT INTO yummy.invite_words ("word") VALUES('betoken');
INSERT INTO yummy.invite_words ("word") VALUES('bewray');
INSERT INTO yummy.invite_words ("word") VALUES('biddy');
INSERT INTO yummy.invite_words ("word") VALUES('bilboes');
INSERT INTO yummy.invite_words ("word") VALUES('blasted');
INSERT INTO yummy.invite_words ("word") VALUES('blazon');
INSERT INTO yummy.invite_words ("word") VALUES('bodements');
INSERT INTO yummy.invite_words ("word") VALUES('bodkin');
INSERT INTO yummy.invite_words ("word") VALUES('bombard');
INSERT INTO yummy.invite_words ("word") VALUES('bootless');
INSERT INTO yummy.invite_words ("word") VALUES('bosky');
INSERT INTO yummy.invite_words ("word") VALUES('bowers');
INSERT INTO yummy.invite_words ("word") VALUES('brach');
INSERT INTO yummy.invite_words ("word") VALUES('brainsickly');
INSERT INTO yummy.invite_words ("word") VALUES('brock');
INSERT INTO yummy.invite_words ("word") VALUES('bruit');
INSERT INTO yummy.invite_words ("word") VALUES('buckler');
INSERT INTO yummy.invite_words ("word") VALUES('busky');
INSERT INTO yummy.invite_words ("word") VALUES('caitiff');
INSERT INTO yummy.invite_words ("word") VALUES('caliver');
INSERT INTO yummy.invite_words ("word") VALUES('callet');
INSERT INTO yummy.invite_words ("word") VALUES('cantons');
INSERT INTO yummy.invite_words ("word") VALUES('carded');
INSERT INTO yummy.invite_words ("word") VALUES('carrions');
INSERT INTO yummy.invite_words ("word") VALUES('cashiered');
INSERT INTO yummy.invite_words ("word") VALUES('casing');
INSERT INTO yummy.invite_words ("word") VALUES('catch');
INSERT INTO yummy.invite_words ("word") VALUES('caterwauling');
INSERT INTO yummy.invite_words ("word") VALUES('cautel');
INSERT INTO yummy.invite_words ("word") VALUES('cerecloth');
INSERT INTO yummy.invite_words ("word") VALUES('cerements');
INSERT INTO yummy.invite_words ("word") VALUES('certes');
INSERT INTO yummy.invite_words ("word") VALUES('champain');
INSERT INTO yummy.invite_words ("word") VALUES('chaps');
INSERT INTO yummy.invite_words ("word") VALUES('charactery');
INSERT INTO yummy.invite_words ("word") VALUES('chariest');
INSERT INTO yummy.invite_words ("word") VALUES('charmingly');
INSERT INTO yummy.invite_words ("word") VALUES('chinks');
INSERT INTO yummy.invite_words ("word") VALUES('chopt');
INSERT INTO yummy.invite_words ("word") VALUES('chough');
INSERT INTO yummy.invite_words ("word") VALUES('civet');
INSERT INTO yummy.invite_words ("word") VALUES('clepe');
INSERT INTO yummy.invite_words ("word") VALUES('climatures');
INSERT INTO yummy.invite_words ("word") VALUES('clodpole');
INSERT INTO yummy.invite_words ("word") VALUES('cobbler');
INSERT INTO yummy.invite_words ("word") VALUES('cockatrices');
INSERT INTO yummy.invite_words ("word") VALUES('collied');
INSERT INTO yummy.invite_words ("word") VALUES('collier');
INSERT INTO yummy.invite_words ("word") VALUES('colour');
INSERT INTO yummy.invite_words ("word") VALUES('compass');
INSERT INTO yummy.invite_words ("word") VALUES('comptible');
INSERT INTO yummy.invite_words ("word") VALUES('conceit');
INSERT INTO yummy.invite_words ("word") VALUES('condition');
INSERT INTO yummy.invite_words ("word") VALUES('continuate');
INSERT INTO yummy.invite_words ("word") VALUES('corky');
INSERT INTO yummy.invite_words ("word") VALUES('coronets');
INSERT INTO yummy.invite_words ("word") VALUES('corse');
INSERT INTO yummy.invite_words ("word") VALUES('coxcomb');
INSERT INTO yummy.invite_words ("word") VALUES('coystrill');
INSERT INTO yummy.invite_words ("word") VALUES('cozen');
INSERT INTO yummy.invite_words ("word") VALUES('cozier');
INSERT INTO yummy.invite_words ("word") VALUES('crisped');
INSERT INTO yummy.invite_words ("word") VALUES('crochets');
INSERT INTO yummy.invite_words ("word") VALUES('crossed');
INSERT INTO yummy.invite_words ("word") VALUES('crowner');
INSERT INTO yummy.invite_words ("word") VALUES('cubiculo');
INSERT INTO yummy.invite_words ("word") VALUES('cursy');
INSERT INTO yummy.invite_words ("word") VALUES('dallying');
INSERT INTO yummy.invite_words ("word") VALUES('dateless');
INSERT INTO yummy.invite_words ("word") VALUES('daws');
INSERT INTO yummy.invite_words ("word") VALUES('denotement');
INSERT INTO yummy.invite_words ("word") VALUES('dilate');
INSERT INTO yummy.invite_words ("word") VALUES('dissemble');
INSERT INTO yummy.invite_words ("word") VALUES('distaff');
INSERT INTO yummy.invite_words ("word") VALUES('distemperature');
INSERT INTO yummy.invite_words ("word") VALUES('doit');
INSERT INTO yummy.invite_words ("word") VALUES('doublet');
INSERT INTO yummy.invite_words ("word") VALUES('doves');
INSERT INTO yummy.invite_words ("word") VALUES('drabbing');
INSERT INTO yummy.invite_words ("word") VALUES('dram');
INSERT INTO yummy.invite_words ("word") VALUES('drossy');
INSERT INTO yummy.invite_words ("word") VALUES('dudgeon');
INSERT INTO yummy.invite_words ("word") VALUES('dunnest');
INSERT INTO yummy.invite_words ("word") VALUES('eanlings');
INSERT INTO yummy.invite_words ("word") VALUES('elflocks');
INSERT INTO yummy.invite_words ("word") VALUES('eliads');
INSERT INTO yummy.invite_words ("word") VALUES('encave');
INSERT INTO yummy.invite_words ("word") VALUES('enchafed');
INSERT INTO yummy.invite_words ("word") VALUES('endues');
INSERT INTO yummy.invite_words ("word") VALUES('engluts');
INSERT INTO yummy.invite_words ("word") VALUES('ensteeped');
INSERT INTO yummy.invite_words ("word") VALUES('envy');
INSERT INTO yummy.invite_words ("word") VALUES('enwheel');
INSERT INTO yummy.invite_words ("word") VALUES('erns');
INSERT INTO yummy.invite_words ("word") VALUES('extremities');
INSERT INTO yummy.invite_words ("word") VALUES('eyeless');
INSERT INTO yummy.invite_words ("word") VALUES('fable');
INSERT INTO yummy.invite_words ("word") VALUES('factious');
INSERT INTO yummy.invite_words ("word") VALUES('fadge');
INSERT INTO yummy.invite_words ("word") VALUES('fain');
INSERT INTO yummy.invite_words ("word") VALUES('fashion');
INSERT INTO yummy.invite_words ("word") VALUES('favour');
INSERT INTO yummy.invite_words ("word") VALUES('festinate');
INSERT INTO yummy.invite_words ("word") VALUES('fetches');
INSERT INTO yummy.invite_words ("word") VALUES('figures');
INSERT INTO yummy.invite_words ("word") VALUES('fleer');
INSERT INTO yummy.invite_words ("word") VALUES('fleering');
INSERT INTO yummy.invite_words ("word") VALUES('flote');
INSERT INTO yummy.invite_words ("word") VALUES('flowerets');
INSERT INTO yummy.invite_words ("word") VALUES('fobbed');
INSERT INTO yummy.invite_words ("word") VALUES('foison');
INSERT INTO yummy.invite_words ("word") VALUES('fopped');
INSERT INTO yummy.invite_words ("word") VALUES('fordid');
INSERT INTO yummy.invite_words ("word") VALUES('forks');
INSERT INTO yummy.invite_words ("word") VALUES('franklin');
INSERT INTO yummy.invite_words ("word") VALUES('frieze');
INSERT INTO yummy.invite_words ("word") VALUES('frippery');
INSERT INTO yummy.invite_words ("word") VALUES('fulsome');
INSERT INTO yummy.invite_words ("word") VALUES('fust');
INSERT INTO yummy.invite_words ("word") VALUES('fustian');
INSERT INTO yummy.invite_words ("word") VALUES('gage');
INSERT INTO yummy.invite_words ("word") VALUES('gaged');
INSERT INTO yummy.invite_words ("word") VALUES('gallow');
INSERT INTO yummy.invite_words ("word") VALUES('gamesome');
INSERT INTO yummy.invite_words ("word") VALUES('gaskins');
INSERT INTO yummy.invite_words ("word") VALUES('gasted');
INSERT INTO yummy.invite_words ("word") VALUES('gauntlet');
INSERT INTO yummy.invite_words ("word") VALUES('gentle');
INSERT INTO yummy.invite_words ("word") VALUES('glazed');
INSERT INTO yummy.invite_words ("word") VALUES('gleek');
INSERT INTO yummy.invite_words ("word") VALUES('goatish');
INSERT INTO yummy.invite_words ("word") VALUES('goodyears');
INSERT INTO yummy.invite_words ("word") VALUES('goose');
INSERT INTO yummy.invite_words ("word") VALUES('gouts');
INSERT INTO yummy.invite_words ("word") VALUES('gramercy');
INSERT INTO yummy.invite_words ("word") VALUES('grise');
INSERT INTO yummy.invite_words ("word") VALUES('grizzled');
INSERT INTO yummy.invite_words ("word") VALUES('groundings');
INSERT INTO yummy.invite_words ("word") VALUES('gudgeon');
INSERT INTO yummy.invite_words ("word") VALUES('gull');
INSERT INTO yummy.invite_words ("word") VALUES('guttered');
INSERT INTO yummy.invite_words ("word") VALUES('hams');
INSERT INTO yummy.invite_words ("word") VALUES('haply');
INSERT INTO yummy.invite_words ("word") VALUES('hardiment');
INSERT INTO yummy.invite_words ("word") VALUES('harpy');
INSERT INTO yummy.invite_words ("word") VALUES('hart');
INSERT INTO yummy.invite_words ("word") VALUES('heath');
INSERT INTO yummy.invite_words ("word") VALUES('hests');
INSERT INTO yummy.invite_words ("word") VALUES('hilding');
INSERT INTO yummy.invite_words ("word") VALUES('hinds');
INSERT INTO yummy.invite_words ("word") VALUES('holidam');
INSERT INTO yummy.invite_words ("word") VALUES('holp');
INSERT INTO yummy.invite_words ("word") VALUES('housewives');
INSERT INTO yummy.invite_words ("word") VALUES('humour');
INSERT INTO yummy.invite_words ("word") VALUES('hurlyburly');
INSERT INTO yummy.invite_words ("word") VALUES('husbandry');
INSERT INTO yummy.invite_words ("word") VALUES('ides');
INSERT INTO yummy.invite_words ("word") VALUES('import');
INSERT INTO yummy.invite_words ("word") VALUES('incarnadine');
INSERT INTO yummy.invite_words ("word") VALUES('indign');
INSERT INTO yummy.invite_words ("word") VALUES('ingraft');
INSERT INTO yummy.invite_words ("word") VALUES('ingrafted');
INSERT INTO yummy.invite_words ("word") VALUES('insuppressive');
INSERT INTO yummy.invite_words ("word") VALUES('intentively');
INSERT INTO yummy.invite_words ("word") VALUES('intermit');
INSERT INTO yummy.invite_words ("word") VALUES('jaunce');
INSERT INTO yummy.invite_words ("word") VALUES('jaundice');
INSERT INTO yummy.invite_words ("word") VALUES('jealous');
INSERT INTO yummy.invite_words ("word") VALUES('jointress');
INSERT INTO yummy.invite_words ("word") VALUES('jowls');
INSERT INTO yummy.invite_words ("word") VALUES('knapped');
INSERT INTO yummy.invite_words ("word") VALUES('ladybird');
INSERT INTO yummy.invite_words ("word") VALUES('leasing');
INSERT INTO yummy.invite_words ("word") VALUES('leman');
INSERT INTO yummy.invite_words ("word") VALUES('lethe');
INSERT INTO yummy.invite_words ("word") VALUES('lief');
INSERT INTO yummy.invite_words ("word") VALUES('liver');
INSERT INTO yummy.invite_words ("word") VALUES('livings');
INSERT INTO yummy.invite_words ("word") VALUES('loath');
INSERT INTO yummy.invite_words ("word") VALUES('loggerheads');
INSERT INTO yummy.invite_words ("word") VALUES('lown');
INSERT INTO yummy.invite_words ("word") VALUES('magnificoes');
INSERT INTO yummy.invite_words ("word") VALUES('maidenhead');
INSERT INTO yummy.invite_words ("word") VALUES('malapert');
INSERT INTO yummy.invite_words ("word") VALUES('marchpane');
INSERT INTO yummy.invite_words ("word") VALUES('marry');
INSERT INTO yummy.invite_words ("word") VALUES('masterless');
INSERT INTO yummy.invite_words ("word") VALUES('maugre');
INSERT INTO yummy.invite_words ("word") VALUES('mazzard');
INSERT INTO yummy.invite_words ("word") VALUES('meet');
INSERT INTO yummy.invite_words ("word") VALUES('meetest');
INSERT INTO yummy.invite_words ("word") VALUES('meiny');
INSERT INTO yummy.invite_words ("word") VALUES('meshes');
INSERT INTO yummy.invite_words ("word") VALUES('micher');
INSERT INTO yummy.invite_words ("word") VALUES('minion');
INSERT INTO yummy.invite_words ("word") VALUES('misprision');
INSERT INTO yummy.invite_words ("word") VALUES('moo');
INSERT INTO yummy.invite_words ("word") VALUES('mooncalf');
INSERT INTO yummy.invite_words ("word") VALUES('mountebanks');
INSERT INTO yummy.invite_words ("word") VALUES('mushrumps');
INSERT INTO yummy.invite_words ("word") VALUES('mute');
INSERT INTO yummy.invite_words ("word") VALUES('naughty');
INSERT INTO yummy.invite_words ("word") VALUES('nonce');
INSERT INTO yummy.invite_words ("word") VALUES('nuncle');
INSERT INTO yummy.invite_words ("word") VALUES('occulted');
INSERT INTO yummy.invite_words ("word") VALUES('ordinary');
INSERT INTO yummy.invite_words ("word") VALUES('othergates');
INSERT INTO yummy.invite_words ("word") VALUES('overname');
INSERT INTO yummy.invite_words ("word") VALUES('paddock');
INSERT INTO yummy.invite_words ("word") VALUES('palmy');
INSERT INTO yummy.invite_words ("word") VALUES('palter');
INSERT INTO yummy.invite_words ("word") VALUES('parle');
INSERT INTO yummy.invite_words ("word") VALUES('patch');
INSERT INTO yummy.invite_words ("word") VALUES('paunch');
INSERT INTO yummy.invite_words ("word") VALUES('pearl');
INSERT INTO yummy.invite_words ("word") VALUES('peize');
INSERT INTO yummy.invite_words ("word") VALUES('pennyworths');
INSERT INTO yummy.invite_words ("word") VALUES('perdy');
INSERT INTO yummy.invite_words ("word") VALUES('pignuts');
INSERT INTO yummy.invite_words ("word") VALUES('portance');
INSERT INTO yummy.invite_words ("word") VALUES('possets');
INSERT INTO yummy.invite_words ("word") VALUES('posy');
INSERT INTO yummy.invite_words ("word") VALUES('praetor');
INSERT INTO yummy.invite_words ("word") VALUES('prate');
INSERT INTO yummy.invite_words ("word") VALUES('prick');
INSERT INTO yummy.invite_words ("word") VALUES('primy');
INSERT INTO yummy.invite_words ("word") VALUES('princox');
INSERT INTO yummy.invite_words ("word") VALUES('prithee');
INSERT INTO yummy.invite_words ("word") VALUES('prodigies');
INSERT INTO yummy.invite_words ("word") VALUES('proper');
INSERT INTO yummy.invite_words ("word") VALUES('prorogued');
INSERT INTO yummy.invite_words ("word") VALUES('pudder');
INSERT INTO yummy.invite_words ("word") VALUES('puddled');
INSERT INTO yummy.invite_words ("word") VALUES('puling');
INSERT INTO yummy.invite_words ("word") VALUES('purblind');
INSERT INTO yummy.invite_words ("word") VALUES('pursy');
INSERT INTO yummy.invite_words ("word") VALUES('quailing');
INSERT INTO yummy.invite_words ("word") VALUES('quaint');
INSERT INTO yummy.invite_words ("word") VALUES('quiddities');
INSERT INTO yummy.invite_words ("word") VALUES('quilets');
INSERT INTO yummy.invite_words ("word") VALUES('quillets');
INSERT INTO yummy.invite_words ("word") VALUES('reference');
INSERT INTO yummy.invite_words ("word") VALUES('instrument');
INSERT INTO yummy.invite_words ("word") VALUES('ranker');
INSERT INTO yummy.invite_words ("word") VALUES('rated');
INSERT INTO yummy.invite_words ("word") VALUES('razes');
INSERT INTO yummy.invite_words ("word") VALUES('receiving');
INSERT INTO yummy.invite_words ("word") VALUES('reechy');
INSERT INTO yummy.invite_words ("word") VALUES('reeking');
INSERT INTO yummy.invite_words ("word") VALUES('remembrances');
INSERT INTO yummy.invite_words ("word") VALUES('rheumy');
INSERT INTO yummy.invite_words ("word") VALUES('rive');
INSERT INTO yummy.invite_words ("word") VALUES('robustious');
INSERT INTO yummy.invite_words ("word") VALUES('romage');
INSERT INTO yummy.invite_words ("word") VALUES('ronyon');
INSERT INTO yummy.invite_words ("word") VALUES('rouse');
INSERT INTO yummy.invite_words ("word") VALUES('sallies');
INSERT INTO yummy.invite_words ("word") VALUES('saws');
INSERT INTO yummy.invite_words ("word") VALUES('scanted');
INSERT INTO yummy.invite_words ("word") VALUES('scarfed');
INSERT INTO yummy.invite_words ("word") VALUES('scrimers');
INSERT INTO yummy.invite_words ("word") VALUES('scutcheon');
INSERT INTO yummy.invite_words ("word") VALUES('seel');
INSERT INTO yummy.invite_words ("word") VALUES('sennet');
INSERT INTO yummy.invite_words ("word") VALUES('sequestration');
INSERT INTO yummy.invite_words ("word") VALUES('shent');
INSERT INTO yummy.invite_words ("word") VALUES('shoon');
INSERT INTO yummy.invite_words ("word") VALUES('shoughs');
INSERT INTO yummy.invite_words ("word") VALUES('shrift');
INSERT INTO yummy.invite_words ("word") VALUES('sleave');
INSERT INTO yummy.invite_words ("word") VALUES('slubber');
INSERT INTO yummy.invite_words ("word") VALUES('smilets');
INSERT INTO yummy.invite_words ("word") VALUES('sonties');
INSERT INTO yummy.invite_words ("word") VALUES('sooth');
INSERT INTO yummy.invite_words ("word") VALUES('sounded');
INSERT INTO yummy.invite_words ("word") VALUES('spleen');
INSERT INTO yummy.invite_words ("word") VALUES('splenetive');
INSERT INTO yummy.invite_words ("word") VALUES('spongy');
INSERT INTO yummy.invite_words ("word") VALUES('springe');
INSERT INTO yummy.invite_words ("word") VALUES('steads');
INSERT INTO yummy.invite_words ("word") VALUES('still');
INSERT INTO yummy.invite_words ("word") VALUES('stoup');
INSERT INTO yummy.invite_words ("word") VALUES('stronds');
INSERT INTO yummy.invite_words ("word") VALUES('suit');
INSERT INTO yummy.invite_words ("word") VALUES('swoopstake');
INSERT INTO yummy.invite_words ("word") VALUES('swounded');
INSERT INTO yummy.invite_words ("word") VALUES('tabor');
INSERT INTO yummy.invite_words ("word") VALUES('taper');
INSERT INTO yummy.invite_words ("word") VALUES('teen');
INSERT INTO yummy.invite_words ("word") VALUES('tenders');
INSERT INTO yummy.invite_words ("word") VALUES('termagant');
INSERT INTO yummy.invite_words ("word") VALUES('tetchy');
INSERT INTO yummy.invite_words ("word") VALUES('tinkers');
INSERT INTO yummy.invite_words ("word") VALUES('topgallant');
INSERT INTO yummy.invite_words ("word") VALUES('traffic');
INSERT INTO yummy.invite_words ("word") VALUES('traject');
INSERT INTO yummy.invite_words ("word") VALUES('trencher');
INSERT INTO yummy.invite_words ("word") VALUES('trimmed');
INSERT INTO yummy.invite_words ("word") VALUES('tristful');
INSERT INTO yummy.invite_words ("word") VALUES('trowest');
INSERT INTO yummy.invite_words ("word") VALUES('truncheon');
INSERT INTO yummy.invite_words ("word") VALUES('unbend');
INSERT INTO yummy.invite_words ("word") VALUES('unbitted');
INSERT INTO yummy.invite_words ("word") VALUES('unbound');
INSERT INTO yummy.invite_words ("word") VALUES('unbraced');
INSERT INTO yummy.invite_words ("word") VALUES('unbruised');
INSERT INTO yummy.invite_words ("word") VALUES('undone');
INSERT INTO yummy.invite_words ("word") VALUES('ungently');
INSERT INTO yummy.invite_words ("word") VALUES('unhoused');
INSERT INTO yummy.invite_words ("word") VALUES('unmake');
INSERT INTO yummy.invite_words ("word") VALUES('unprevailing');
INSERT INTO yummy.invite_words ("word") VALUES('unprovide');
INSERT INTO yummy.invite_words ("word") VALUES('unreclaimed');
INSERT INTO yummy.invite_words ("word") VALUES('unstuffed');
INSERT INTO yummy.invite_words ("word") VALUES('untaught');
INSERT INTO yummy.invite_words ("word") VALUES('untented');
INSERT INTO yummy.invite_words ("word") VALUES('unthrifty');
INSERT INTO yummy.invite_words ("word") VALUES('unyoke');
INSERT INTO yummy.invite_words ("word") VALUES('usance');
INSERT INTO yummy.invite_words ("word") VALUES('vailing');
INSERT INTO yummy.invite_words ("word") VALUES('varlets');
INSERT INTO yummy.invite_words ("word") VALUES('verdure');
INSERT INTO yummy.invite_words ("word") VALUES('villanies');
INSERT INTO yummy.invite_words ("word") VALUES('vizards');
INSERT INTO yummy.invite_words ("word") VALUES('wafter');
INSERT INTO yummy.invite_words ("word") VALUES('welkin');
INSERT INTO yummy.invite_words ("word") VALUES('weraday');
INSERT INTO yummy.invite_words ("word") VALUES('whoreson');
INSERT INTO yummy.invite_words ("word") VALUES('wilt');
INSERT INTO yummy.invite_words ("word") VALUES('windlasses');
INSERT INTO yummy.invite_words ("word") VALUES('yarely');
INSERT INTO yummy.invite_words ("word") VALUES('yerked');
INSERT INTO yummy.invite_words ("word") VALUES('yoeman');
INSERT INTO yummy.invite_words ("word") VALUES('younker');



-- CREATE TABLE "invites" ------------------------------------
CREATE TABLE "yummy"."invites" (
    "id" Serial NOT NULL,
    "referrer_id" Integer NOT NULL,
    "word1" Integer NOT NULL,
    "word2" Integer NOT NULL,
    "word3" Integer NOT NULL,
    CONSTRAINT "unique_invite_id" PRIMARY KEY( "id" ),
    CONSTRAINT "invite_word1" FOREIGN KEY("word1") REFERENCES "yummy"."invite_words"("id"),
    CONSTRAINT "invite_word2" FOREIGN KEY("word2") REFERENCES "yummy"."invite_words"("id"),
    CONSTRAINT "invite_word3" FOREIGN KEY("word3") REFERENCES "yummy"."invite_words"("id") );
;
-- -------------------------------------------------------------

-- CREATE INDEX "index_referrer_id" -------------------------
CREATE INDEX "index_referrer_id" ON "yummy"."invites" USING btree( "referrer_id" );
-- -------------------------------------------------------------

INSERT INTO yummy.invites (referrer_id, word1, word2, word3) VALUES(1, 1, 1, 1);
INSERT INTO yummy.invites (referrer_id, word1, word2, word3) VALUES(1, 1, 1, 1);
INSERT INTO yummy.invites (referrer_id, word1, word2, word3) VALUES(1, 1, 1, 1);

CREATE VIEW yummy.unwrapped_invites AS
SELECT invites.id AS id, 
    users.id AS user_id,
    lower(users.name) AS name, 
    one.word AS word1, 
    two.word AS word2, 
    three.word AS word3
FROM yummy.invites, yummy.users,
    yummy.invite_words AS one, 
    yummy.invite_words AS two, 
    yummy.invite_words AS three
WHERE invites.referrer_id = users.id
    AND invites.word1 = one.id 
    AND invites.word2 = two.id 
    AND invites.word3 = three.id;



-- CREATE TABLE "relation" -------------------------------------
CREATE TABLE "yummy"."relation" (
    "id" Integer UNIQUE NOT NULL,
    "type" Text NOT NULL );

INSERT INTO "yummy"."relation" VALUES(0, 'followed');
INSERT INTO "yummy"."relation" VALUES(1, 'requested');
INSERT INTO "yummy"."relation" VALUES(2, 'cancelled');
INSERT INTO "yummy"."relation" VALUES(3, 'ignored');
-- -------------------------------------------------------------



-- CREATE TABLE "relations" ------------------------------------
CREATE TABLE "yummy"."relations" ( 
	"from_id" Integer NOT NULL,
	"to_id" Integer NOT NULL,
	"type" Integer NOT NULL,
	"changed_at" Timestamp With Time Zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
	CONSTRAINT "unique_relation" PRIMARY KEY ("from_id" , "to_id"),
    CONSTRAINT "unique_from_relation" FOREIGN KEY ("from_id") REFERENCES "yummy"."users"("id"),
    CONSTRAINT "unique_to_relation" FOREIGN KEY ("to_id") REFERENCES "yummy"."users"("id"),
    CONSTRAINT "enum_relation_type" FOREIGN KEY("type") REFERENCES "yummy"."relation"("id") );
 ;
-- -------------------------------------------------------------

-- CREATE INDEX "index_related_to_users" -----------------------
CREATE INDEX "index_related_to_users" ON "yummy"."relations" USING btree( "to_id" );
-- -------------------------------------------------------------

-- CREATE INDEX "index_related_from_users" ---------------------
CREATE INDEX "index_related_from_users" ON "yummy"."relations" USING btree( "from_id" );
-- -------------------------------------------------------------

CREATE OR REPLACE FUNCTION yummy.count_relations_ins() RETURNS TRIGGER AS $$
    BEGIN
        IF (NEW."type" = (SELECT id FROM yummy.relation WHERE "type" = 'followed')) THEN
            UPDATE yummy.users
            SET followers_count = followers_count + 1
            WHERE id = NEW.to_id;
            UPDATE yummy.users
            SET followings_count = followings_count + 1
            WHERE id = NEW.from_id;
        ELSIF (NEW."type" = (SELECT id FROM yummy.relation WHERE "type" = 'ignored')) THEN
            UPDATE yummy.users
            SET ignored_count = ignored_count + 1
            WHERE id = NEW.from_id;
        END IF;
        
        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER cnt_relations_ins
    AFTER INSERT OR UPDATE ON yummy.relations
    FOR EACH ROW EXECUTE PROCEDURE yummy.count_relations_ins();

CREATE OR REPLACE FUNCTION yummy.count_relations_del() RETURNS TRIGGER AS $$
    BEGIN
        IF (OLD."type" = (SELECT id FROM yummy.relation WHERE "type" = 'followed')) THEN
            UPDATE yummy.users
            SET followers_count = followers_count - 1
            WHERE id = OLD.to_id;
            UPDATE yummy.users
            SET followings_count = followings_count - 1
            WHERE id = OLD.from_id;
        ELSIF (OLD."type" = (SELECT id FROM yummy.relation WHERE "type" = 'ignored')) THEN
            UPDATE users
            SET ignored_count = ignored_count - 1
            WHERE id = OLD.from_id;
        END IF;
    
        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER cnt_relations_del
    AFTER UPDATE OR DELETE ON yummy.relations
    FOR EACH ROW EXECUTE PROCEDURE yummy.count_relations_del();



-- CREATE TABLE "entry_privacy" --------------------------------
CREATE TABLE "yummy"."entry_privacy" (
    "id" Integer UNIQUE NOT NULL,
    "type" Text NOT NULL );

INSERT INTO "yummy"."entry_privacy" VALUES(0, 'all');
INSERT INTO "yummy"."entry_privacy" VALUES(1, 'some');
INSERT INTO "yummy"."entry_privacy" VALUES(2, 'me');
INSERT INTO "yummy"."entry_privacy" VALUES(3, 'anonymous');
-- -------------------------------------------------------------



-- CREATE TABLE "entries" --------------------------------------
CREATE TABLE "yummy"."entries" ( 
	"id" Serial NOT NULL,
	"created_at" Timestamp With Time Zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"author_id" Integer NOT NULL,
	"rating" Integer DEFAULT 0 NOT NULL,
	"title" Text DEFAULT '' NOT NULL,
	"content" Text NOT NULL,
	"word_count" Integer NOT NULL,
	"visible_for" Integer NOT NULL,
	"is_votable" Boolean NOT NULL,
	"comments_count" Integer DEFAULT 0 NOT NULL,
	CONSTRAINT "unique_entry_id" PRIMARY KEY( "id" ),
    CONSTRAINT "entry_user_id" FOREIGN KEY("author_id") REFERENCES "yummy"."users"("id"),
    CONSTRAINT "enum_entry_privacy" FOREIGN KEY("visible_for") REFERENCES "yummy"."entry_privacy"("id") );
 ;
-- -------------------------------------------------------------

-- CREATE INDEX "index_entry_id" -------------------------------
CREATE INDEX "index_entry_id" ON "yummy"."entries" USING btree( "id" );
-- -------------------------------------------------------------

-- CREATE INDEX "index_entry_date" -----------------------------
CREATE INDEX "index_entry_date" ON "yummy"."entries" USING btree( "created_at" );
-- -------------------------------------------------------------

-- CREATE INDEX "index_entry_users_id" -------------------------
CREATE INDEX "index_entry_users_id" ON "yummy"."entries" USING btree( "author_id" );
-- -------------------------------------------------------------

-- CREATE INDEX "index_entry_rating" ---------------------------
CREATE INDEX "index_entry_rating" ON "yummy"."entries" USING btree( "rating" );
-- -------------------------------------------------------------

-- CREATE INDEX "index_entry_word_count" -----------------------
CREATE INDEX "index_entry_word_count" ON "yummy"."entries" USING btree( "word_count" );
-- -------------------------------------------------------------

CREATE OR REPLACE FUNCTION yummy.inc_tlog_entries() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE yummy.users
        SET entries_count = entries_count + 1
        WHERE id = NEW.author_id;
        
        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION yummy.dec_tlog_entries() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE yummy.users
        SET entries_count = entries_count - 1
        WHERE id = OLD.author_id;
        
        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER cnt_tlog_entries_ins
    AFTER INSERT ON yummy.entries
    FOR EACH ROW 
    WHEN (NEW.visible_for = 0) -- visible_for = all
    EXECUTE PROCEDURE yummy.inc_tlog_entries();

CREATE TRIGGER cnt_tlog_entries_upd_inc
    AFTER UPDATE ON yummy.entries
    FOR EACH ROW 
    WHEN (OLD.visible_for <> 0 AND NEW.visible_for = 0)
    EXECUTE PROCEDURE yummy.inc_tlog_entries();

CREATE TRIGGER cnt_tlog_entries_upd_dec
    AFTER UPDATE ON yummy.entries
    FOR EACH ROW 
    WHEN (OLD.visible_for = 0 AND NEW.visible_for <> 0)
    EXECUTE PROCEDURE yummy.dec_tlog_entries();

CREATE TRIGGER cnt_tlog_entries_del
    AFTER DELETE ON yummy.entries
    FOR EACH ROW 
    WHEN (OLD.visible_for = 0)
    EXECUTE PROCEDURE yummy.dec_tlog_entries();

CREATE VIEW yummy.feed AS
SELECT entries.id, entries.created_at, rating, 
    entries.title, content, word_count,
    entry_privacy.type AS entry_privacy,
    is_votable, entries.comments_count,
    long_users.id AS author_id,
    long_users.name AS author_name, 
    long_users.show_name AS author_show_name,
    long_users.is_online AS author_is_online,
    long_users.name_color AS author_name_color, 
    long_users.avatar_color AS author_avatar_color, 
    long_users.avatar AS author_avatar,
    long_users.privacy AS author_privacy
FROM yummy.long_users, yummy.entries, yummy.entry_privacy
WHERE long_users.id = entries.author_id 
    AND entry_privacy.id = entries.visible_for
ORDER BY entries.created_at DESC;



-- CREATE TABLE "tags" -----------------------------------------
CREATE TABLE "yummy"."tags" (
    "id" Serial NOT NULL,
    "tag" Text NOT NULL,
    CONSTRAINT "unique_tag_id" PRIMARY KEY( "id" ) );
;
-- -------------------------------------------------------------

-- CREATE INDEX "index_tag" ------------------------------------
CREATE UNIQUE INDEX "index_tag" ON "yummy"."tags" USING btree( "tag" ) ;
-- -------------------------------------------------------------



-- CREATE TABLE "entry_tags" -----------------------------------
CREATE TABLE "yummy"."entry_tags" (
    "entry_id" Integer NOT NULL,
    "tag_id" Integer NOT NULL,
    CONSTRAINT "entry_tags_entry" FOREIGN KEY("entry_id") REFERENCES "yummy"."entries"("id"),
    CONSTRAINT "entry_tags_tag" FOREIGN KEY("tag_id") REFERENCES "yummy"."tags"("id"),
    CONSTRAINT "unique_entry_tag" UNIQUE("entry_id", "tag_id") );
;
-- -------------------------------------------------------------

-- CREATE INDEX "index_tag" ------------------------------------
CREATE INDEX "index_entry_tags_entry" ON "yummy"."entry_tags" USING btree( "entry_id" );
-- -------------------------------------------------------------

-- CREATE INDEX "index_tag" ------------------------------------
CREATE INDEX "index_entry_tags_tag" ON "yummy"."entry_tags" USING btree( "tag_id" );
-- -------------------------------------------------------------

CREATE OR REPLACE FUNCTION yummy.count_tags() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE yummy.users
        SET tags_count = counts.cnt 
        FROM yummy.users,
        (
            SELECT DISTINCT author_id as id
            FROM yummy.entries, changes
            WHERE entries.id = changes.entry_id
        ) AS authors,
        (
            SELECT author_id, COUNT(tag_id) as cnt
            FROM yummy.entries, yummy.entry_tags, authors
            WHERE authors.id = entries.author_id AND entries.id = entry_tags.entry_id
            GROUP BY author_id
        ) AS counts
        WHERE authors.id = users.id AND counts.author_id = users.id;
        
        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER cnt_tags_ins
    AFTER INSERT ON yummy.entry_tag
    REFERENCING NEW TABLE as changes
    FOR EACH STATEMENT EXECUTE PROCEDURE yummy.count_tags();

CREATE TRIGGER cnt_tags_del
    AFTER DELETE ON yummy.entry_tags
    REFERENCING OLD TABLE as changes
    FOR EACH STATEMENT EXECUTE PROCEDURE yummy.count_tags();



-- CREATE TABLE "favorites" ------------------------------------
CREATE TABLE "yummy"."favorites" ( 
	"user_id" Integer NOT NULL,
	"entry_id" Integer NOT NULL,
    CONSTRAINT "favorite_user_id" FOREIGN KEY("user_id") REFERENCES "yummy"."users"("id"),
    CONSTRAINT "favorite_entry_id" FOREIGN KEY("entry_id") REFERENCES "yummy"."entries"("id"),
    CONSTRAINT "unique_user_favorite" UNIQUE("user_id", "entry_id") );
 ;
-- -------------------------------------------------------------

-- CREATE INDEX "index_favorite_entries" -----------------------
CREATE INDEX "index_favorite_entries" ON "yummy"."favorites" USING btree( "entry_id" );
-- -------------------------------------------------------------

-- CREATE INDEX "index_favorite_users" -------------------------
CREATE INDEX "index_favorite_users" ON "yummy"."favorites" USING btree( "user_id" );
-- -------------------------------------------------------------

CREATE OR REPLACE FUNCTION yummy.inc_favorites() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE yummy.users
        SET favorites_count = favorites_count + 1
        WHERE id = NEW.user_id;
        
        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION yummy.dec_favorites() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE yummy.users
        SET favorites_count = favorites_count - 1
        WHERE id = OLD.user_id;
        
        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER cnt_favorites_inc
    AFTER INSERT ON yummy.favorites
    FOR EACH ROW EXECUTE PROCEDURE yummy.inc_favorites();

CREATE TRIGGER cnt_favorites_dec
    AFTER DELETE ON yummy.favorites
    FOR EACH ROW EXECUTE PROCEDURE yummy.dec_favorites();



-- CREATE TABLE "watching" -------------------------------------
CREATE TABLE "yummy"."watching" ( 
	"user_id" Integer NOT NULL,
	"entry_id" Integer NOT NULL,
    CONSTRAINT "watching_user_id" FOREIGN KEY("user_id") REFERENCES "yummy"."users"("id"),
    CONSTRAINT "watching_entry_id" FOREIGN KEY("entry_id") REFERENCES "yummy"."entries"("id"),
    CONSTRAINT "unique_user_watching" UNIQUE("user_id", "entry_id") );
 ;
-- -------------------------------------------------------------

-- CREATE INDEX "index_watching_entries" -----------------------
CREATE INDEX "index_watching_entries" ON "yummy"."watching" USING btree( "entry_id" );
-- -------------------------------------------------------------

-- CREATE INDEX "index_watching_users" -------------------------
CREATE INDEX "index_watching_users" ON "yummy"."watching" USING btree( "user_id" );
-- -------------------------------------------------------------



-- CREATE TABLE "entry_votes" ----------------------------------
CREATE TABLE "yummy"."entry_votes" ( 
	"user_id" Integer NOT NULL,
	"entry_id" Integer NOT NULL,
	"positive" Boolean NOT NULL,
	"taken" Boolean DEFAULT TRUE NOT NULL,
    CONSTRAINT "entry_vote_user_id" FOREIGN KEY("user_id") REFERENCES "yummy"."users"("id"),
    CONSTRAINT "entry_vote_entry_id" FOREIGN KEY("entry_id") REFERENCES "yummy"."entries"("id"),
    CONSTRAINT "unique_entry_vote" UNIQUE("user_id", "entry_id") );
 ;
-- -------------------------------------------------------------

-- CREATE INDEX "index_voted_entries" --------------------------
CREATE INDEX "index_voted_entries" ON "yummy"."entry_votes" USING btree( "entry_id" );
-- -------------------------------------------------------------

-- CREATE INDEX "index_voted_users" ----------------------------
CREATE INDEX "index_voted_users" ON "yummy"."entry_votes" USING btree( "user_id" );
-- -------------------------------------------------------------

CREATE OR REPLACE FUNCTION yummy.inc_entry_votes() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE yummy.entries
        SET rating = rating + 1
        WHERE id = NEW.entry_id;
        
        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION yummy.dec_entry_votes() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE yummy.entries
        SET rating = rating - 1
        WHERE id = OLD.entry_id;
        
        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION yummy.inc_entry_votes2() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE yummy.entries
        SET rating = rating + 2
        WHERE id = NEW.entry_id;
        
        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION yummy.dec_entry_votes2() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE yummy.entries
        SET rating = rating - 2
        WHERE id = OLD.entry_id;
        
        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER cnt_entry_votes_ins_inc
    AFTER INSERT ON yummy.entry_votes
    FOR EACH ROW 
    WHEN (NEW."positive" = true)
    EXECUTE PROCEDURE yummy.inc_entry_votes();

CREATE TRIGGER cnt_entry_votes_ins_dec
    AFTER INSERT ON yummy.entry_votes
    FOR EACH ROW 
    WHEN (NEW."positive" = false)
    EXECUTE PROCEDURE yummy.dec_entry_votes();

CREATE TRIGGER cnt_entry_votes_upd_inc
    AFTER UPDATE ON yummy.entry_votes
    FOR EACH ROW 
    WHEN (NEW."positive" = true AND NEW."positive" <> OLD."positive")
    EXECUTE PROCEDURE yummy.inc_entry_votes2();

CREATE TRIGGER cnt_entry_votes_upd_dec
    AFTER UPDATE ON yummy.entry_votes
    FOR EACH ROW 
    WHEN (NEW."positive" = false AND NEW."positive" <> OLD."positive")
    EXECUTE PROCEDURE yummy.dec_entry_votes2();

CREATE TRIGGER cnt_entry_votes_del_dec
    AFTER DELETE ON yummy.entry_votes
    FOR EACH ROW 
    WHEN (OLD."positive" = true)
    EXECUTE PROCEDURE yummy.dec_entry_votes();

CREATE TRIGGER cnt_entry_votes_del_inc
    AFTER DELETE ON yummy.entry_votes
    FOR EACH ROW 
    WHEN (OLD."positive" = false)
    EXECUTE PROCEDURE yummy.inc_entry_votes();



-- CREATE TABLE "entries_privacy" ------------------------------
CREATE TABLE "yummy"."entries_privacy" ( 
	"user_id" Integer NOT NULL,
	"entry_id" Integer NOT NULL,
    CONSTRAINT "entries_privacy_user_id" FOREIGN KEY("user_id") REFERENCES "yummy"."users"("id"),
    CONSTRAINT "entries_privacy_entry_id" FOREIGN KEY("entry_id") REFERENCES "yummy"."entries"("id"),
    CONSTRAINT "unique_entry_privacy" UNIQUE("user_id", "entry_id") );
 ;
-- -------------------------------------------------------------

-- CREATE INDEX "index_private_entries" ------------------------
CREATE INDEX "index_private_entries" ON "yummy"."entries_privacy" USING btree( "entry_id" );
-- -------------------------------------------------------------

-- CREATE INDEX "index_private_users" --------------------------
CREATE INDEX "index_private_users" ON "yummy"."entries_privacy" USING btree( "user_id" );
-- -------------------------------------------------------------



-- CREATE TABLE "comments" -------------------------------------
CREATE TABLE "yummy"."comments" ( 
	"id" Serial NOT NULL,
	"author_id" Integer NOT NULL,
	"entry_id" Integer NOT NULL,
	"created_at" Timestamp With Time Zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"content" Text NOT NULL,
	"rating" Integer DEFAULT 0 NOT NULL,
	CONSTRAINT "unique_comment_id" PRIMARY KEY( "id" ),
    CONSTRAINT "comment_user_id" FOREIGN KEY("author_id") REFERENCES "yummy"."users"("id"),
    CONSTRAINT "comment_entry_id" FOREIGN KEY("entry_id") REFERENCES "yummy"."entries"("id") );
 ;
-- -------------------------------------------------------------

-- CREATE INDEX "index_entry_id" -------------------------------
CREATE INDEX "index_comment_entry_id" ON "yummy"."comments" USING btree( "entry_id" );
-- -------------------------------------------------------------

CREATE OR REPLACE FUNCTION yummy.inc_comments() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE yummy.users
        SET comments_count = comments_count + 1
        WHERE id = NEW.author_id;
        
        UPDATE yummy.entries
        SET comments_count = comments_count + 1
        WHERE id = NEW.entry_id;
        
        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER cnt_comments_inc
    AFTER INSERT ON yummy.comments
    FOR EACH ROW EXECUTE PROCEDURE yummy.inc_comments();

CREATE OR REPLACE FUNCTION yummy.dec_comments() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE yummy.users
        SET comments_count = comments_count - 1
        WHERE id = OLD.author_id;
        
        UPDATE yummy.entries
        SET comments_count = comments_count - 1
        WHERE id = OLD.entry_id;
        
        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER cnt_comments_dec
    AFTER DELETE ON yummy.comments
    FOR EACH ROW EXECUTE PROCEDURE yummy.dec_comments();



-- CREATE TABLE "comment_votes" --------------------------------
CREATE TABLE "yummy"."comment_votes" ( 
	"user_id" Integer NOT NULL,
	"comment_id" Integer NOT NULL,
	"positive" Boolean NOT NULL,
	"taken" Boolean DEFAULT TRUE NOT NULL,
    CONSTRAINT "comment_vote_user_id" FOREIGN KEY("user_id") REFERENCES "yummy"."users"("id"),
    CONSTRAINT "comment_vote_comment_id" FOREIGN KEY("comment_id") REFERENCES "yummy"."comments"("id"),
    CONSTRAINT "unique_comment_vote" UNIQUE("user_id", "comment_id") );
 ;
-- -------------------------------------------------------------

-- CREATE INDEX "index_voted_comments" -------------------------
CREATE INDEX "index_voted_comments" ON "yummy"."comment_votes" USING btree( "comment_id" );
-- -------------------------------------------------------------

-- CREATE INDEX "index_comment_voted_users" --------------------
CREATE INDEX "index_comment_voted_users" ON "yummy"."comment_votes" USING btree( "user_id" );
-- -------------------------------------------------------------

CREATE OR REPLACE FUNCTION yummy.inc_comment_votes() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE yummy.comments
        SET rating = rating + 1
        WHERE id = NEW.comment_id;
        
        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION yummy.dec_comment_votes() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE yummy.comments
        SET rating = rating - 1
        WHERE id = OLD.comment_id;
        
        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION yummy.inc_comment_votes2() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE yummy.comments
        SET rating = rating + 2
        WHERE id = NEW.comment_id;
        
        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION yummy.dec_comment_votes2() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE yummy.comments
        SET rating = rating - 2
        WHERE id = OLD.comment_id;
        
        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER cnt_comment_votes_ins_inc
    AFTER INSERT ON yummy.comment_votes
    FOR EACH ROW 
    WHEN (NEW."positive" = true)
    EXECUTE PROCEDURE yummy.inc_comment_votes();

CREATE TRIGGER cnt_comment_votes_ins_dec
    AFTER INSERT ON yummy.comment_votes
    FOR EACH ROW 
    WHEN (NEW."positive" = false)
    EXECUTE PROCEDURE yummy.dec_comment_votes();

CREATE TRIGGER cnt_comment_votes_upd_inc
    AFTER UPDATE ON yummy.comment_votes
    FOR EACH ROW 
    WHEN (NEW."positive" = true)
    EXECUTE PROCEDURE yummy.inc_comment_votes2();

CREATE TRIGGER cnt_comment_votes_upd_dec
    AFTER UPDATE ON yummy.comment_votes
    FOR EACH ROW 
    WHEN (NEW."positive" = false)
    EXECUTE PROCEDURE yummy.dec_comment_votes2();

CREATE TRIGGER cnt_comment_votes_del_dec
    AFTER DELETE ON yummy.comment_votes
    FOR EACH ROW 
    WHEN (OLD."positive" = true)
    EXECUTE PROCEDURE yummy.dec_comment_votes();

CREATE TRIGGER cnt_comment_votes_del_inc
    AFTER DELETE ON yummy.comment_votes
    FOR EACH ROW 
    WHEN (OLD."positive" = false)
    EXECUTE PROCEDURE yummy.inc_comment_votes();

    

-- -- CREATE TABLE "images" ---------------------------------------
-- CREATE TABLE "yummy"."images" (
-- 	"id" Serial NOT NULL,
-- 	"entry_id" Integer NOT NULL,
-- 	"url" Text NOT NULL,
-- 	"sorting" SmallInt NOT NULL,
-- 	CONSTRAINT "unique_image_id" PRIMARY KEY( "id" ),
--     CONSTRAINT "image_entry_id" FOREIGN KEY("entry_id") REFERENCES "yummy"."entries"("id") );
--  ;
-- -- -------------------------------------------------------------
--
-- -- CREATE INDEX "index_image_id" -------------------------------
-- CREATE INDEX "index_image_id" ON "yummy"."images" USING btree( "id" );
-- -- -------------------------------------------------------------
--
-- -- CREATE INDEX "index_image_entry" ----------------------------
-- CREATE INDEX "index_image_entry" ON "yummy"."images" USING btree( "entry_id" );
-- -- -------------------------------------------------------------
--
--
--
-- -- CREATE TABLE "media" ----------------------------------------
-- CREATE TABLE "yummy"."media" (
-- 	"id" Serial NOT NULL,
-- 	"duration" Integer NOT NULL,
-- 	"icon" Text NOT NULL,
-- 	"preview" Text NOT NULL,
-- 	"title" Text NOT NULL,
-- 	"url" Text NOT NULL,
-- 	"entry_id" Integer NOT NULL,
-- 	CONSTRAINT "unique_media_id" PRIMARY KEY( "id" ),
--     CONSTRAINT "media_entry_id" FOREIGN KEY("entry_id") REFERENCES "yummy"."entries"("id") );
--  ;
-- -- -------------------------------------------------------------
--
-- -- CREATE INDEX "index_media_id" -------------------------------
-- CREATE INDEX "index_media_id" ON "yummy"."media" USING btree( "id" );
-- -- -------------------------------------------------------------
--
-- -- CREATE INDEX "index_media_entry" ----------------------------
-- CREATE INDEX "index_media_entry" ON "yummy"."media" USING btree( "entry_id" );
-- -- -------------------------------------------------------------
--
--
--
-- -- CREATE TABLE "chats" ----------------------------------------
-- CREATE TABLE "yummy"."chats" (
-- 	"id" Serial NOT NULL,
-- 	"messages_count" Integer DEFAULT 0 NOT NULL,
-- 	"avatar" Text DEFAULT '' NOT NULL,
-- 	CONSTRAINT "unique_chat_id" PRIMARY KEY( "id" ) );
--  ;
-- -- -------------------------------------------------------------
--
-- -- CREATE INDEX "index_chat_id" --------------------------------
-- CREATE INDEX "index_chat_id" ON "yummy"."chats" USING btree( "id" );
-- -- -------------------------------------------------------------
--
--
--
-- -- CREATE TABLE "messages" -------------------------------------
-- CREATE TABLE "yummy"."messages" (
-- 	"id" Serial NOT NULL,
-- 	"chat_id" Integer NOT NULL,
-- 	"author_id" Integer NOT NULL,
-- 	"created_at" Timestamp With Time Zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
-- 	"content" Text NOT NULL,
-- 	"reply_to" Integer,
-- 	CONSTRAINT "unique_message_id" PRIMARY KEY( "id" ),
--     CONSTRAINT "message_user_id" FOREIGN KEY("author_id") REFERENCES "yummy"."users"("id"),
--     CONSTRAINT "message_chat_id" FOREIGN KEY("chat_id") REFERENCES "yummy"."chats"("id"),
--     CONSTRAINT "message_reply_to" FOREIGN KEY("reply_to") REFERENCES "yummy"."comments"("id") );
--  ;
-- -- -------------------------------------------------------------
--
-- -- CREATE INDEX "index_message_id" -----------------------------
-- CREATE INDEX "index_message_id" ON "yummy"."messages" USING btree( "id" );
-- -- -------------------------------------------------------------
--
-- -- CREATE INDEX "index_message_chat" ---------------------------
-- CREATE INDEX "index_message_chat" ON "yummy"."messages" USING btree( "chat_id" );
-- -- -------------------------------------------------------------
--
-- CREATE OR REPLACE FUNCTION yummy.count_messages() RETURNS TRIGGER AS $$
--     DECLARE
--         delta   integer;
--         chat_id integer;
--
--     BEGIN
--         IF (TG_OP = 'INSERT') THEN
--             delta = 1;
--             chat_id = NEW.chat_id;
--         ELSIF (TG_OP = 'DELETE') THEN
--             delta = -1;
--             chat_id = OLD.chat_id;
--         END IF;
--
--         UPDATE yummy.chats
--         SET messages_count = messages_count + delta
--         WHERE id = chat_id;
--
--         RETURN NULL;
--     END;
-- $$ LANGUAGE plpgsql;
--
-- CREATE TRIGGER cnt_messages
--     AFTER INSERT OR DELETE ON yummy.messages
--     FOR EACH ROW EXECUTE PROCEDURE yummy.count_messages();
--
--
--
-- -- CREATE TABLE "talker_status" --------------------------------
-- CREATE TABLE "yummy"."talker_status" (
--     "id" Integer NOT NULL,
--     "type" Text NOT NULL );
--
-- INSERT INTO "yummy"."talker_status" VALUES(0, "creator");
-- INSERT INTO "yummy"."talker_status" VALUES(1, "banned");
-- INSERT INTO "yummy"."talker_status" VALUES(2, "normal");
-- INSERT INTO "yummy"."talker_status" VALUES(3, "left");
-- INSERT INTO "yummy"."talker_status" VALUES(4, "admin");
-- -- -------------------------------------------------------------
--
--
--
-- -- CREATE TABLE "talking" --------------------------------------
-- CREATE TABLE "yummy"."talking" (
-- 	"chat_id" Integer NOT NULL,
-- 	"last_read" Integer,
-- 	"user_id" Integer NOT NULL,
-- 	"unread_count" Text NOT NULL,
-- 	"status" Integer NOT NULL,
-- 	"not_disturb" Boolean DEFAULT false NOT NULL,
--     CONSTRAINT "talking_user_id" FOREIGN KEY("user_id") REFERENCES "yummy"."users"("id"),
--     CONSTRAINT "talking_chat_id" FOREIGN KEY("chat_id") REFERENCES "yummy"."chats"("id"),
--     CONSTRAINT "enum_talking_status" FOREIGN KEY("status") REFERENCES "yummy"."talker_status"("id"));
--  ;
-- -- -------------------------------------------------------------
--
-- -- CREATE INDEX "index_talking_chat" ---------------------------
-- CREATE INDEX "index_talking_chat" ON "yummy"."talking" USING btree( "chat_id" );
-- -- -------------------------------------------------------------
--
-- -- CREATE INDEX "index_talking_user" ---------------------------
-- CREATE INDEX "index_talking_user" ON "yummy"."talking" USING btree( "user_id" )
--     WHERE "status" NOT IN (
--         SELECT "id" from "talker_status"
--         WHERE "type" = "banned" OR "type" = "left");
-- -- -------------------------------------------------------------
--
-- CREATE OR REPLACE FUNCTION yummy.count_unread() RETURNS TRIGGER AS $$
--     BEGIN
--         IF (TG_OP = 'INSERT') THEN
--             UPDATE yummy.talking
--             SET unread_count = unread_count + 1
--             WHERE talking.chat_id = NEW.chat_id AND talking.user_id <> NEW.user_id
--         ELSIF (TG_OP = 'DELETE') THEN
--             UPDATE yummy.talking
--             SET unread_count = unread_count -1
--             WHERE talking.chat_id = OLD.chat_id AND talking.user_id <> OLD.user_id
--                 AND (last_read = NULL OR last_read < OLD.id)
--         END IF;
--
--         RETURN NULL;
--     END;
-- $$ LANGUAGE plpgsql;
--
-- CREATE TRIGGER cnt_unread
--     AFTER INSERT OR DELETE ON yummy.messages
--     FOR EACH ROW EXECUTE PROCEDURE yummy.count_unread();
