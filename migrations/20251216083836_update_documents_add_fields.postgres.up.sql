-- migrations/..._update_documents_add_fields.up.sql
ALTER TABLE documents
ADD COLUMN fio VARCHAR(255),
ADD COLUMN lavozim VARCHAR(255),
ADD COLUMN oylik VARCHAR(255),
ADD COLUMN stavka VARCHAR(50),
ADD COLUMN username VARCHAR(255),
ADD COLUMN order_type VARCHAR(255);