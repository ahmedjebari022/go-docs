-- +goose Up 
CREATE TABLE documents (
    id UUID NOT NULL PRIMARY KEY,
    name VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    document_url VARCHAR(255) NOT NULL UNIQUE
);



-- +goose Down
DROP TABLE documents; 
