# ******************************************************************************
# Settings
# ******************************************************************************
SET foreign_key_checks = 1;
SET time_zone = '+00:00';

# ******************************************************************************
# Create tables
# ******************************************************************************
CREATE TABLE note (
    id INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,

    name TEXT NOT NULL,

    user_id INT(10) UNSIGNED NOT NULL,

    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,

    CONSTRAINT `f_note_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,

    PRIMARY KEY (id)
);

