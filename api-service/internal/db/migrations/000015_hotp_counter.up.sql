CREATE TABLE hotp_counters (
    user_id UUID PRIMARY KEY NOT NULL,
    counter INTEGER NOT NULL,
    last_used TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE  CASCADE
)