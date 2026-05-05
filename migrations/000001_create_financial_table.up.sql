CREATE TABLE IF NOT EXISTS rols (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    userName VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    rol_id INTEGER NOT NULL REFERENCES rols(id)
);

CREATE TABLE IF NOT EXISTS stores (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    balance BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS storeAccess (
    store_id BIGINT NOT NULL REFERENCES stores(id),
    user_id BIGINT NOT NULL REFERENCES users(id),

    PRIMARY KEY (store_id, user_id)
);

CREATE TABLE IF NOT EXISTS vaults (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    store_id BIGINT NOT NULL REFERENCES stores(id),
    name VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    balance BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS amountCategory (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    description VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS movements (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    detail TEXT NOT NULL,
    amount BIGINT NOT NULL,
    amount_category_id INTEGER NOT NULL REFERENCES amountCategory(id),
    vault_id BIGINT NOT NULL REFERENCES vaults(id),
    user_id BIGINT NOT NULL REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS transfers (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    amount BIGINT NOT NULL CHECK (amount > 0),
    description TEXT,

    status TEXT NOT NULL DEFAULT 'PENDING'
        CHECK (status IN ('PENDING', 'APPROVED', 'REJECTED')),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    approved_at TIMESTAMP,
    rejected_at TIMESTAMP,

    vault_orig_id BIGINT NOT NULL REFERENCES vaults(id),
    vault_dest_id BIGINT NOT NULL REFERENCES vaults(id),

    sender_user_id BIGINT NOT NULL REFERENCES users(id),
    receiver_user_id BIGINT REFERENCES users(id),

    CONSTRAINT different_vaults
        CHECK (vault_orig_id <> vault_dest_id),

    CONSTRAINT valid_transfer_status
        CHECK (
            (status = 'APPROVED' AND approved_at IS NOT NULL AND rejected_at IS NULL)
         OR (status = 'REJECTED' AND rejected_at IS NOT NULL AND approved_at IS NULL)
         OR (status = 'PENDING' AND approved_at IS NULL AND rejected_at IS NULL)
        )
);

