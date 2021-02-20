CREATE TABLE books(
  isbn        TEXT PRIMARY KEY,
  title       TEXT,
  author      TEXT,
  publisher   TEXT,
  pubdate     TEXT,
  cover       TEXT,
  keywords    TEXT,
  ccode       TEXT,
  genre       TEXT,
  description TEXT,
  contents    TEXT,
  created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_books_updated_at
AFTER UPDATE ON books
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE INDEX ON books (title);
