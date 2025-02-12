CREATE TABLE IF NOT exists followers (
    user_id INT NOT NULL,
    follower_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (user_id, follower_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT check_user_id_follower_id CHECK (user_id != follower_id)
);
