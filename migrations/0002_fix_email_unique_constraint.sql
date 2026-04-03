BEGIN;

-- Remove a constraint UNIQUE da coluna email
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_key;

-- Garante idempotência quando a migration roda mais de uma vez
DROP INDEX IF EXISTS users_email_unique_idx;

-- Altera a coluna email para permitir NULL
ALTER TABLE users ALTER COLUMN email DROP NOT NULL;

-- Cria um índice parcial UNIQUE que só se aplica quando o email não é NULL e não é vazio
-- Isso permite múltiplos usuários sem email, mas garante unicidade quando o email é fornecido
CREATE UNIQUE INDEX users_email_unique_idx ON users (email) WHERE email IS NOT NULL AND email != '';

COMMIT;
