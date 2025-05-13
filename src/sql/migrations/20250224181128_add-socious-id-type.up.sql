ALTER TYPE oauth_connected_providers ADD VALUE 'SOCIOUS_ID';
CREATE UNIQUE INDEX unique_provider_mui ON oauth_connects (matrix_unique_id, provider);