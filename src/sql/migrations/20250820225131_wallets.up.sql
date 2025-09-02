CREATE TYPE network_type AS ENUM ('bsc', 'sepolia', 'cardano', 'midnight');

CREATE TABLE wallets (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    address TEXT NOT NULL,
    network network_type NOT NULL,
    testnet boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);

CREATE UNIQUE INDEX idx_wallets_user_id_network ON wallets (user_id, network, testnet);
