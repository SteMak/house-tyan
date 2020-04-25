CREATE OR REPLACE FUNCTION insert_time ()
    RETURNS TRIGGER
    AS $$
BEGIN
    NEW.inserted_at = now() at time zone 'utc';
    RETURN NEW;
END;
$$
LANGUAGE 'plpgsql';

CREATE OR REPLACE FUNCTION update_time ()
    RETURNS TRIGGER
    AS $$
BEGIN
    NEW.updated_at = now() at time zone 'utc';
    RETURN NEW;
END;
$$
LANGUAGE 'plpgsql';

