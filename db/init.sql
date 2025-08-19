\c ingresso_go;

-- ksuid extension
CREATE OR REPLACE FUNCTION ksuid() RETURNS TEXT AS $$

DECLARE
	V_TIME TIMESTAMP WITH TIME ZONE := NULL;
	V_SECONDS NUMERIC(50) := NULL;
	V_NUMERIC NUMERIC(50) := NULL;
	V_EPOCH NUMERIC(50) := 1400000000;
	V_BASE62 TEXT := '';
	V_ALPHABET CHAR ARRAY[62] := ARRAY[
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J',
		'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 
		'U', 'V', 'W', 'X', 'Y', 'Z', 
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 
		'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't',
		'u', 'v', 'w', 'x', 'y', 'z'];
	i INTEGER := 0;
BEGIN
	V_TIME := clock_timestamp();
	V_SECONDS := floor(EXTRACT(EPOCH FROM V_TIME)) - V_EPOCH;
	V_NUMERIC := V_SECONDS * pow(2::NUMERIC(50), 128) 
		+ ((random()::NUMERIC(70,20) * pow(2::NUMERIC(70,20), 48))::NUMERIC(50) * pow(2::NUMERIC(50), 80)::NUMERIC(50))
		+ ((random()::NUMERIC(70,20) * pow(2::NUMERIC(70,20), 40))::NUMERIC(50) * pow(2::NUMERIC(50), 40)::NUMERIC(50))
		+  (random()::NUMERIC(70,20) * pow(2::NUMERIC(70,20), 40))::NUMERIC(50);

	while V_NUMERIC <> 0 loop
		V_BASE62 := V_BASE62 || V_ALPHABET[mod(V_NUMERIC, 62) + 1];
		V_NUMERIC := div(V_NUMERIC, 62);
	end loop;
	V_BASE62 := reverse(V_BASE62);
	V_BASE62 := lpad(V_BASE62, 27, '0');
	return V_BASE62;
	
end $$ language plpgsql;

CREATE TABLE IF NOT EXISTS session (
    id VARCHAR(27) PRIMARY KEY DEFAULT ksuid(),
    movie_id VARCHAR NOT NULL,
	date DATE NOT NULL,
	start_time VARCHAR NOT NULL,
    room VARCHAR NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(movie_id, start_time, room)
);

CREATE TABLE IF NOT EXISTS ticket (
    id VARCHAR(27) PRIMARY KEY DEFAULT ksuid(),
    session_id VARCHAR(27) REFERENCES session(id) ON DELETE CASCADE,
    user_id VARCHAR NOT NULL,
    seats VARCHAR NOT NULL,
    email VARCHAR NOT NULL,
    amount INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(session_id, seats)
);