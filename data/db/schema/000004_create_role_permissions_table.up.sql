CREATE TABLE role_permissions(
                                id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                role_id UUID NOT NULL,
                                permission_id UUID NOT NULL,
                                created_at TIMESTAMP NOT NULL DEFAULT now(),
                                CONSTRAINT fk_role_permissions_role
                                    FOREIGN KEY (role_id)
                                        REFERENCES roles(id)
                                        ON DELETE CASCADE,
                                CONSTRAINT fk_role_permissions_permission
                                    FOREIGN KEY (permission_id)
                                        REFERENCES permissions(id)
                                        ON DELETE CASCADE,
                                CONSTRAINT uq_role_permission UNIQUE (role_id, permission_id)
);

CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id);
