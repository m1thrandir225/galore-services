CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
  email TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  password TEXT NOT NULL,
  avatar_url TEXT NOT NULL, 
  enabled_push_notifications BOOLEAN DEFAULT FALSE NOT NULL,
  enabled_email_notifications BOOLEAN DEFAULT FALSE NOT NULL, 
  created_at TIMESTAMPTZ DEFAULT 'now()' NOT NULL,
)
