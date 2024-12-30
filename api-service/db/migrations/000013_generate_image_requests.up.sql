CREATE TYPE image_generation_status as ENUM (
    'generating',
    'success',
    'error'
);

CREATE TABLE generate_image_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    draft_id UUID NOT NULL REFERENCES generate_cocktail_drafts(id),
    prompt TEXT NOT NULL,
    status image_generation_status NOT NULL DEFAULT  'generating',
    image_url TEXT,
    error_message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_image_requests_status ON generate_image_requests(status);
