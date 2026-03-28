CREATE TABLE generate_cocktail_drafts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    request_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    instructions JSONB NOT NULL,
    ingredients JSONB NOT NULL,
    main_image_prompt TEXT not null,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (request_id) REFERENCES generate_cocktail_requests(id) ON DELETE CASCADE,
    UNIQUE(request_id)
);