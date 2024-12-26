CREATE TABLE reset_password_request(
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    user_id UUID  NOT NULL,
    password_reset BOOLEAN NOT NULL DEFAULT false,
    valid_until TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
)