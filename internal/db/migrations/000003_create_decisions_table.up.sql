CREATE TABLE IF NOT EXISTS decisions (
    actor_user_id UUID NOT NULL,
    recipient_user_id UUID NOT NULL,
    liked_recipient BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(actor_user_id, recipient_user_id),
    FOREIGN KEY (actor_user_id) REFERENCES users(user_id),
    FOREIGN KEY (recipient_user_id) REFERENCES users(user_id)
);
