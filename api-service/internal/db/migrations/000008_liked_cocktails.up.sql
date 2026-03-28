CREATE TABLE liked_cocktails (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  cocktail_id UUID REFERENCES cocktails (id) ON DELETE CASCADE NOT NULL, 
  user_id UUID REFERENCES users (id) ON DELETE CASCADE NOT NULL
);
