CREATE TABLE categories (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 name TEXT NOT NULL,
 tag TEXT UNIQUE NOT NULL,
 created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);


CREATE TABLE cocktail_categories (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  cocktail_id UUID REFERENCES cocktails (id) ON DELETE CASCADE NOT NULL,
  category_id UUID REFERENCES categories (id) ON DELETE CASCADE NOT NULL
);


