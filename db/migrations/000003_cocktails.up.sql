CREATE TABLE cocktails (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  is_alcoholic BOOLEAN DEFAULT TRUE,
  glass TEXT NOT NULL,
  image TEXT NOT NULL,
  instructions JSONB NOT NULL,
  ingredients JSONB NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
