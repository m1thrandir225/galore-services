CREATE TYPE generation_status as ENUM (
    'generating_cocktail',
    'generating_images',
    'error',
    'success'
);

CREATE TABLE generate_cocktail_requests (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    prompt TEXT NOT NULL,
    status generation_status NOT NULL DEFAULT 'generating_cocktail',
    error_message TEXT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_requests_status ON generate_cocktail_requests(status);