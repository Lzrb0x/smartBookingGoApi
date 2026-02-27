BEGIN;

-- users: identifies any user of the system
CREATE TABLE IF NOT EXISTS users (
    id              BIGSERIAL PRIMARY KEY,
    active          BOOLEAN       NOT NULL DEFAULT TRUE,
    created_on      TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    user_identifier VARCHAR(100)  NOT NULL UNIQUE,
    name            VARCHAR(150)  NOT NULL,
    email           VARCHAR(254)  NOT NULL UNIQUE,
    password        TEXT          NOT NULL,
    phone           VARCHAR(30),
    is_complete     BOOLEAN       NOT NULL DEFAULT FALSE
);

-- owners: dono da barbearia
CREATE TABLE IF NOT EXISTS owners (
    id      BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id)
);

-- barbershops: estabelecimento
CREATE TABLE IF NOT EXISTS barbershops (
    id               BIGSERIAL PRIMARY KEY,
    barbershop_name  VARCHAR(200) NOT NULL,
    address          TEXT,
    phone            VARCHAR(30),
    owner_id         BIGINT NOT NULL REFERENCES owners(id)
);

-- employees: funcionário/barbeiro
CREATE TABLE IF NOT EXISTS employees (
    id             BIGSERIAL PRIMARY KEY,
    user_id        BIGINT NOT NULL REFERENCES users(id),
    barbershop_id  BIGINT NOT NULL REFERENCES barbershops(id)
);

-- services: catálogo global de serviços
CREATE TABLE IF NOT EXISTS services (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(150) NOT NULL,
    description TEXT
);

-- barbershop_services: customização do serviço para a barbearia
CREATE TABLE IF NOT EXISTS barbershop_services (
    id                   BIGSERIAL PRIMARY KEY,
    price                NUMERIC(10, 2) NOT NULL,
    duration             INT            NOT NULL,
    description_override TEXT,
    barbershop_id        BIGINT NOT NULL REFERENCES barbershops(id),
    service_id           BIGINT NOT NULL REFERENCES services(id)
);

-- service_employees: associação N-N entre funcionário e serviço da barbearia
CREATE TABLE IF NOT EXISTS service_employees (
    id                    BIGSERIAL PRIMARY KEY,
    employee_id           BIGINT NOT NULL REFERENCES employees(id),
    barbershop_service_id BIGINT NOT NULL REFERENCES barbershop_services(id),
    UNIQUE (employee_id, barbershop_service_id)
);

-- employee_working_hours: horário padrão semanal
CREATE TABLE IF NOT EXISTS employee_working_hours (
    id          BIGSERIAL PRIMARY KEY,
    employee_id BIGINT NOT NULL REFERENCES employees(id),
    day_of_week SMALLINT NOT NULL CHECK (day_of_week BETWEEN 0 AND 6),
    start_time  TIME,
    end_time    TIME,
    is_day_off  BOOLEAN NOT NULL DEFAULT FALSE
);

-- employee_working_hour_overrides: exceções pontuais de horário
CREATE TABLE IF NOT EXISTS employee_working_hour_overrides (
    id          BIGSERIAL PRIMARY KEY,
    employee_id BIGINT NOT NULL REFERENCES employees(id),
    date        DATE   NOT NULL,
    start_time  TIME,
    end_time    TIME,
    is_day_off  BOOLEAN NOT NULL DEFAULT FALSE
);

-- bookings: agendamentos confirmados
CREATE TABLE IF NOT EXISTS bookings (
    id                    BIGSERIAL PRIMARY KEY,
    customer_id           BIGINT NOT NULL REFERENCES users(id),
    employee_id           BIGINT NOT NULL REFERENCES employees(id),
    barbershop_id         BIGINT NOT NULL REFERENCES barbershops(id),
    barbershop_service_id BIGINT NOT NULL REFERENCES barbershop_services(id),
    date                  DATE   NOT NULL,
    start_time            TIME   NOT NULL,
    end_time              TIME   NOT NULL
);

COMMIT;
