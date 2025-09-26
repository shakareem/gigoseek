CREATE TABLE IF NOT EXISTS city (
  id SERIAL PRIMARY KEY,
  city_name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS auth_state (
  state TEXT PRIMARY KEY,
  chat_id BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS chat (
  chat_id BIGINT PRIMARY KEY,
  chat_state INTEGER NOT NULL DEFAULT 0,
  city_id INTEGER REFERENCES city(id)
);

CREATE TABLE IF NOT EXISTS token (
  chat_id BIGINT PRIMARY KEY REFERENCES chat(chat_id) ON DELETE CASCADE,
  access_token TEXT NOT NULL,
  token_type TEXT NOT NULL,
  refresh_token TEXT NOT NULL,
  expiry TIMESTAMPTZ NOT NULL
);

INSERT INTO city (city_name) VALUES
  ('Москва'),
  ('Санкт-Петербург'),
  ('Казань')
ON CONFLICT (city_name) DO NOTHING;
