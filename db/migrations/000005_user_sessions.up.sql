CREATE TABLE sessions (
  id UUID PRIMARY KEY,
  email TEXT REFERENCES users (email) ON DELETE CASCADE NOT NULL,
  refresh_token TEXT NOT NULL,
  user_agent TEXT NOT NULL,
  client_ip TEXT NOT NULL, 
  is_blocked BOOLEAN NOT NULL DEFAULT FALSE, 
  expires_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
