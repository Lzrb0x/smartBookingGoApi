BEGIN;

CREATE UNIQUE INDEX IF NOT EXISTS idx_owners_user_id_unique ON owners (user_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_employees_user_barbershop_unique ON employees (user_id, barbershop_id);

COMMIT;
