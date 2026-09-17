CREATE TABLE friend_requests (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    sender_id BIGINT UNSIGNED NOT NULL,
    receiver_id BIGINT UNSIGNED NOT NULL,

    status ENUM(
        'pending',
        'accepted',
        'rejected'
    ) NOT NULL DEFAULT 'pending',

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (id),

    KEY idx_friend_requests_sender (sender_id),
    KEY idx_friend_requests_receiver (receiver_id),
    KEY idx_friend_requests_status (status),

    CONSTRAINT fk_friend_requests_sender
        FOREIGN KEY (sender_id)
        REFERENCES users (id)
        ON DELETE CASCADE,

    CONSTRAINT fk_friend_requests_receiver
        FOREIGN KEY (receiver_id)
        REFERENCES users (id)
        ON DELETE CASCADE,

    CONSTRAINT chk_friend_requests_different_users
        CHECK (sender_id <> receiver_id)
) ENGINE=InnoDB
DEFAULT CHARSET=utf8mb4
COLLATE=utf8mb4_0900_ai_ci;
