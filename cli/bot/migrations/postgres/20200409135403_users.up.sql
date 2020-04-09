CREATE TABLE users (
    username varchar(32) NOT NULL,
    discriminator varchar(4) NOT NULL
)
INHERITS (
    base
);

CREATE TRIGGER tr_on_user_insert
    BEFORE INSERT ON users
    FOR EACH ROW
    EXECUTE PROCEDURE insert_time ();

CREATE TRIGGER tr_on_user_update
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE PROCEDURE update_time ();