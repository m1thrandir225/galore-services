CREATE TABLE created_cocktails (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL, 
  image TEXT NOT NULL,
  ingredients JSONB NOT NULL, 
  instructions JSONB NOT NULL, 
  description TEXT NOT NULL, 
  user_id UUID REFERENCES users (id) ON DELETE CASCADE NOT NULL,
  embedding vector(768) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE fcm_tokens (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  token TEXT NOT NULL, 
  device_id TEXT NOT NULL,
  user_id UUID REFERENCES users (id) ON DELETE CASCADE NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE liked_flavours (
  flavour_id UUID REFERENCES flavours (id) ON DELETE CASCADE NOT NULL, 
  user_id UUID REFERENCES users (id) ON DELETE CASCADE NOT NULL, 

  PRIMARY KEY(flavour_id, user_id)
);

CREATE TABLE notifications (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users (id) ON DELETE CASCADE NOT NULL, 
  notification_type_id UUID REFERENCES notification_types (id) ON DELETE RESTRICT NOT NULL,
  opened BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
