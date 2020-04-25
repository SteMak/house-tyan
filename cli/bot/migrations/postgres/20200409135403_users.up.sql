CREATE TABLE users (
    id varchar(18) NOT NULL PRIMARY KEY,
    inserted_at timestamp without time zone NULL,
    updated_at timestamp without time zone NULL,
    username varchar(32) NOT NULL DEFAULT '',
    discriminator varchar(4) NOT NULL DEFAULT '',
    xp bigint NOT NULL DEFAULT 0,
    balance bigint NOT NULL DEFAULT 0,
    club_id bigint NULL
);

CREATE TRIGGER on_insert
    BEFORE INSERT ON users
    FOR EACH ROW
    EXECUTE PROCEDURE insert_time ();

CREATE TRIGGER on_update
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE PROCEDURE update_time ();

