CREATE TABLE club_invites (
    club_id int NOT NULL,
    user_id varchar(18) NOT NULL,
    is_request boolean NOT NULL,
    message_id varchar(18) NOT NULL,
    inserted_at timestamp without time zone NULL,
    updated_at timestamp without time zone NULL,
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