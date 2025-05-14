CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    message_content VARCHAR(140) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    status INT8 NOT NULL,
    external_message_id VARCHAR(100)
);

CREATE INDEX idx_messages_external_message_id ON messages (external_message_id);
CREATE INDEX idx_messages_created_at ON messages (created_at);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE PLPGSQL;

CREATE TRIGGER update_messages_updated_at
BEFORE UPDATE ON messages
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();