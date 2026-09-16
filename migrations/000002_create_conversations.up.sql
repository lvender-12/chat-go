CREATE TABLE conversations (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    user_one_id BIGINT UNSIGNED NOT NULL,
    user_two_id BIGINT UNSIGNED NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT fk_conversations_user_one FOREIGN KEY (user_one_id)
        REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_conversations_user_two FOREIGN KEY (user_two_id)
        REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT chk_conversations_different_users CHECK (user_one_id <> user_two_id),
    UNIQUE KEY uq_conversations_users (user_one_id, user_two_id),
    KEY idx_conversations_user_one (user_one_id),
    KEY idx_conversations_user_two (user_two_id)
) ENGINE=InnoDB;
