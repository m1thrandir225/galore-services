CREATE TABLE flavours (
      id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
      name TEXT NOT NULL UNIQUE,
      created_at TIMESTAMPTZ NOT NULL DEFAULT now ()
);