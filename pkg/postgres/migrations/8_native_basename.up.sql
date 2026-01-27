CREATE OR REPLACE FUNCTION basename(path text, slash text DEFAULT '/')
RETURNS text AS $$
DECLARE
    re text;
    result text;
BEGIN
    IF slash IS NULL OR length(slash) != 1 THEN
        RAISE EXCEPTION 'Slash must be a single character';
    END IF;

    re := '[^' || regexp_replace(slash, '(.)', '\\\\\1', 'g') || ']+$';

    result := substring(path FROM re);

    RETURN result;
END;
$$ LANGUAGE plpgsql IMMUTABLE STRICT;
