CREATE TABLE IF NOT EXISTS likes (
    id SERIAL PRIMARY KEY,
    actor_user_id INT NOT NULL,
    recipient_user_id INT NOT NULL,
    liked_recipient BOOLEAN NOT NULL
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(actor_user_id, recipient_user_id),
    FOREIGN KEY (actor_user_id) REFERENCES users(user_id),
    FOREIGN KEY (recipient_user_id) REFERENCES users(user_id)
);

CREATE INDEX IF NOT EXISTS idx_likes_recipient_user_id ON likes(recipient_user_id);
CREATE INDEX IF NOT EXISTS idx_likes_actor_recipient ON likes(actor_user_id, recipient_user_id);
