CREATE TABLE awards (
    id bigserial PRIMARY KEY,
    inserted_at timestamp without time zone NULL,
    updated_at timestamp without time zone NULL,
    author_id varchar(18) NOT NULL,
    blank_mid varchar(18) UNIQUE,
    reason text,
    status smallint NOT NULL DEFAULT 0
);

CREATE TRIGGER on_insert
    BEFORE INSERT ON awards
    FOR EACH ROW
    EXECUTE PROCEDURE insert_time ();

CREATE TRIGGER on_update
    BEFORE UPDATE ON awards
    FOR EACH ROW
    EXECUTE PROCEDURE update_time ();

CREATE INDEX ids_awards_blank_mid ON awards (blank_mid);

