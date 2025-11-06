CREATE OR REPLACE FUNCTION phash_distance(lhash bigint, rhash bigint)
RETURNS bigint
LANGUAGE sql
IMMUTABLE
AS $$
  SELECT length(replace(((lhash::bit(64) # rhash::bit(64))::text), '0', ''));
$$;
