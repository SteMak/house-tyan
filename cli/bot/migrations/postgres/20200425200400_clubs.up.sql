CREATE TABLE clubs (
    id serial PRIMARY KEY,
    inserted_at timestamp without time zone NULL,
    updated_at timestamp without time zone NULL,
    owner_id varchar(18) NOT NULL UNIQUE,
    role_id varchar(18) UNIQUE,
    channel_id varchar(18) UNIQUE,
    title varchar(255) NOT NULL UNIQUE,
    description text NULL,
    symbol varchar(60) NOT NULL UNIQUE,
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

