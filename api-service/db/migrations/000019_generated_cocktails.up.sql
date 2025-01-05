CREATE TABLE generated_cocktails (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    request_id UUID NOT NULL,
    draft_id UUID NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    main_image_url TEXT NOT NULL,
    instructions JSONB NOT NULL,
    ingredients JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (request_id) REFERENCES generate_cocktail_requests(id) ON DELETE RESTRICT,
    FOREIGN KEY (draft_id) REFERENCES generate_cocktail_drafts(id) ON DELETE RESTRICT,
    UNIQUE(request_id),
    UNIQUE(draft_id)
);