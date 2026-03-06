BEGIN;

CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY,
    name TEXT NOT NULL
);

INSERT INTO users (id, name) VALUES
	(1, 'Natalie'),
	(2, 'Domingo'),
	(11, 'Lucy'),
	(99, 'Rasicov')
ON CONFLICT (id) -- When id has conflict
DO UPDATE -- Update the existing data with the new one
SET name = EXCLUDED.name; -- With replacing the existing (EXCLUDED) name

COMMIT;
