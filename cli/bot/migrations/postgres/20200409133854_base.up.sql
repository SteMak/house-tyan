CREATE TABLE base (
    id varchar(18) NOT NULL PRIMARY KEY,
    inserted_at timestamp without time zone NULL,
    updated_at timestamp without time zone NULL
);

CREATE FUNCTION insert_time ()
    RETURNS TRIGGER
    AS $$
BEGIN
    NEW.inserted_at = now() at time zone 'utc';
    RETURN NEW;
END;
$$
LANGUAGE 'plpgsql';

CREATE FUNCTION update_time ()
    RETURNS TRIGGER
    AS $$
BEGIN
    NEW.updated_at = now() at time zone 'utc';
    RETURN NEW;
END;
$$
LANGUAGE 'plpgsql';

