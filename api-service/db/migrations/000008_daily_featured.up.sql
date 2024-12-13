    CREATE TABLE daily_featured_cocktails (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        cocktail_id UUID references cocktails (id) ON DELETE CASCADE NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT (CURRENT_DATE AT TIME ZONE 'UTC')::timestamp
    );
