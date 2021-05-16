CREATE TABLE club_invites (
    club_id int NOT NULL,
    user_id varchar(18) NOT NULL,
    is_request boolean NOT NULL,
    inserted_at timestamp without time zone NULL,
    updated_at timestamp without time zone NULL,
    status int NOT NULL DEFAULT 0,
    PRIMARY KEY(club_id, user_id)
);

CREATE TRIGGER on_insert
    BEFORE INSERT ON club_invites
    FOR EACH ROW
    EXECUTE PROCEDURE insert_time ();

CREATE TRIGGER on_update
    BEFORE UPDATE ON club_invites
    FOR EACH ROW
    EXECUTE PROCEDURE update_time ();