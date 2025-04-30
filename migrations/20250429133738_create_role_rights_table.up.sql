CREATE TABLE role_rights (
    id SERIAL PRIMARY KEY,
    role_id INT NOT NULL,
    section VARCHAR(255) NOT NULL,
    route VARCHAR(255) NOT NULL,
    r_create TINYINT DEFAULT 0,
    r_read TINYINT DEFAULT 0,
    r_update TINYINT DEFAULT 0,
    r_delete TINYINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_role_rights_role FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE
);

CREATE INDEX idx_role_rights_role_id ON role_rights (role_id);
