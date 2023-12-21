CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION trigger_set_timestamp()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS users (
    id              serial primary key,
    chat_id         varchar (20) not null,
    name            text not null,
    external_info   text default '',

    created_at      timestamp default CURRENT_TIMESTAMP,
    updated_at      timestamp default CURRENT_TIMESTAMP
);

CREATE TRIGGER set_timestamp
    BEFORE UPDATE ON users
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TYPE session_status AS ENUM (
    'WAITING', 'PROCESSING', 'DONE'
);

CREATE TYPE session_type AS ENUM (
    'REVIEW', 'INFO'
    );

CREATE TABLE IF NOT EXISTS sessions (
    id          serial primary key,
    external_id UUID NOT NULL DEFAULT uuid_generate_v1(),
    user_id     integer NOT NULL,
    status      session_status DEFAULT 'WAITING',
    type      session_type DEFAULT 'REVIEW',
    product_name text DEFAULT '',
    created_at  timestamp default CURRENT_TIMESTAMP,
    updated_at  timestamp default CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TRIGGER set_timestamp
    BEFORE UPDATE ON sessions
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TABLE IF NOT EXISTS session_logs (
    id          serial primary key,
    session_id  integer NOT NULL,
    text        text NOT NULL,
    created_at  timestamp default CURRENT_TIMESTAMP,

    FOREIGN KEY (session_id) REFERENCES sessions (id)
);


