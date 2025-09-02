INSERT INTO wallets (id, user_id, address, network, testnet) VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (user_id, network, testnet) DO UPDATE SET
    address = EXCLUDED.address,
    updated_at = NOW();
