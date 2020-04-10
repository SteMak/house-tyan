CREATE TABLE awards (
    author_id varchar(18) NOT NULL,
    reason text,
    status smallint NOT NULL DEFAULT 0,
    PRIMARY KEY (id)
)
INHERITS (
    base
);

CREATE TRIGGER tr_on_awards_insert
    BEFORE INSERT ON awards
    FOR EACH ROW
    EXECUTE PROCEDURE insert_time ();

CREATE TRIGGER tr_on_awards_update
    BEFORE UPDATE ON awards
    FOR EACH ROW
    EXECUTE PROCEDURE update_time ();

