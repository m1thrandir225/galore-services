CREATE TABLE category_flavour (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id UUID references categories(id) ON DELETE CASCADE NOT NULL,
    flavour_id UUID references  flavours(id) ON DELETE CASCADE  NOT NULL
);
