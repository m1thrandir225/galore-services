CREATE TABLE cocktail_categories (
     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
     cocktail_id UUID REFERENCES cocktails (id) ON DELETE CASCADE NOT NULL,
     category_id UUID REFERENCES categories (id) ON DELETE CASCADE NOT NULL
);


