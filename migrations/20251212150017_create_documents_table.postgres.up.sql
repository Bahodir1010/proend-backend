-- migrations/..._create_documents_table.up.sql
CREATE TABLE documents (
    id UUID PRIMARY KEY,
    template_id UUID NOT NULL, -- Which template was used
    filename VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'draft',
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);