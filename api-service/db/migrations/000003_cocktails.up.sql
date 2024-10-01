CREATE TABLE cocktails (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  is_alcoholic BOOLEAN DEFAULT TRUE,
  glass TEXT NOT NULL,
  image TEXT NOT NULL,
  instructions TEXT NOT NULL,
  ingredients JSONB NOT NULL,
  embedding vector(768) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);


CREATE TABLE liked_cocktails (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  cocktail_id UUID REFERENCES cocktails (id) ON DELETE CASCADE NOT NULL, 
  user_id UUID REFERENCES users (id) ON DELETE CASCADE NOT NULL
);
