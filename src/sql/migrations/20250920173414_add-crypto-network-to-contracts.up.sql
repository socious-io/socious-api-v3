ALTER TABLE contracts ADD COLUMN crypto_network network_type;

-- Update existing crypto contracts to have BSC as the default network
UPDATE contracts
SET crypto_network = 'bsc'
WHERE crypto_currency IS NOT NULL
  AND crypto_network IS NULL;