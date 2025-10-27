CREATE TABLE roles(
                      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                      title_ru VARCHAR(255) NOT NULL,
                      title_en VARCHAR(255),
                      title_kk VARCHAR(255),
                      description_ru TEXT NOT NULL,
                      description_kk TEXT,
                      description_en TEXT,
                      value VARCHAR(280) UNIQUE NOT NULL,
                      created_at TIMESTAMP NOT NULL DEFAULT now(),
                      updated_at TIMESTAMP NOT NULL DEFAULT now(),
                      deleted_at TIMESTAMP
);

CREATE INDEX idx_roles_value ON roles(value);