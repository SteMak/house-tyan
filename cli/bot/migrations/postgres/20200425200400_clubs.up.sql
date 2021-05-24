CREATE TABLE clubs (
    id serial PRIMARY KEY,
    inserted_at timestamp without time zone NULL,
    updated_at timestamp without time zone NULL,
    owner_id varchar(20) NOT NULL UNIQUE,
    role_id varchar(20) UNIQUE,
    role_color varchar(6) DEFAULT 'ffffff',
    role_mentionable boolean NOT NULL DEFAULT false,
    channel_id varchar(20) UNIQUE,
    card_message_id varchar(20),
    title varchar(32) NOT NULL UNIQUE,
    description varchar(1024) NULL,
    symbol varchar(32) NOT NULL UNIQUE,
    icon_url varchar(128) NULL,
    xp bigint NOT NULL DEFAULT 0,
    expired_at timestamp without time zone NULL,
    verified boolean NOT NULL DEFAULT false
);

CREATE TRIGGER on_insert
    BEFORE INSERT ON clubs
    FOR EACH ROW
    EXECUTE PROCEDURE insert_time ();

CREATE TRIGGER on_update
    BEFORE UPDATE ON clubs
    FOR EACH ROW
    EXECUTE PROCEDURE update_time ();

