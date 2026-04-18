CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS warehouses (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name         VARCHAR(255) NOT NULL,
    office_from_id VARCHAR(100) NOT NULL UNIQUE,
    address      TEXT        NOT NULL DEFAULT '',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS routes (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    warehouse_id UUID        NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
    route_id     VARCHAR(100) NOT NULL,
    name         VARCHAR(255) NOT NULL DEFAULT '',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(warehouse_id, route_id)
);

CREATE TABLE IF NOT EXISTS users (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    username      VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role          VARCHAR(20)  NOT NULL CHECK (role IN ('admin', 'driver')),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS drivers (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    warehouse_id UUID REFERENCES warehouses(id) ON DELETE SET NULL,
    route_id     UUID REFERENCES routes(id) ON DELETE SET NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS thresholds (
    id           UUID    PRIMARY KEY DEFAULT gen_random_uuid(),
    warehouse_id UUID    NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
    route_id     UUID    NOT NULL REFERENCES routes(id) ON DELETE CASCADE,
    value        FLOAT8  NOT NULL,
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(warehouse_id, route_id)
);

CREATE TABLE IF NOT EXISTS forecasts (
    id              UUID    PRIMARY KEY DEFAULT gen_random_uuid(),
    warehouse_id    UUID    NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
    route_id        UUID    NOT NULL REFERENCES routes(id) ON DELETE CASCADE,
    forecast_time   TIMESTAMPTZ NOT NULL,
    horizon_hours   INT     NOT NULL DEFAULT 2,
    predicted_count FLOAT8  NOT NULL,
    actual_count    FLOAT8,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_forecasts_warehouse_route ON forecasts(warehouse_id, route_id);
CREATE INDEX IF NOT EXISTS idx_forecasts_time ON forecasts(forecast_time);

CREATE TABLE IF NOT EXISTS truck_calls (
    id               UUID    PRIMARY KEY DEFAULT gen_random_uuid(),
    warehouse_id     UUID    NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
    route_id         UUID    NOT NULL REFERENCES routes(id) ON DELETE CASCADE,
    driver_id        UUID    REFERENCES drivers(id) ON DELETE SET NULL,
    forecast_value   FLOAT8  NOT NULL,
    threshold_value  FLOAT8  NOT NULL,
    called_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    status           VARCHAR(20) NOT NULL DEFAULT 'pending'
                     CHECK (status IN ('pending', 'accepted', 'completed', 'missed')),
    timeliness       VARCHAR(20) CHECK (timeliness IN ('on_time', 'late', 'early')),
    timeliness_at    TIMESTAMPTZ,
    actual_containers INT
);

CREATE INDEX IF NOT EXISTS idx_truck_calls_warehouse_route ON truck_calls(warehouse_id, route_id);
CREATE INDEX IF NOT EXISTS idx_truck_calls_driver ON truck_calls(driver_id);
CREATE INDEX IF NOT EXISTS idx_truck_calls_status ON truck_calls(status);

CREATE TABLE IF NOT EXISTS raw_data (
    id            UUID    PRIMARY KEY DEFAULT gen_random_uuid(),
    route_id      VARCHAR(100) NOT NULL,
    office_from_id VARCHAR(100) NOT NULL,
    timestamp     TIMESTAMPTZ NOT NULL,
    status_1      FLOAT8,
    status_2      FLOAT8,
    status_3      FLOAT8,
    status_4      FLOAT8,
    status_5      FLOAT8,
    status_6      FLOAT8,
    status_7      FLOAT8,
    status_8      FLOAT8,
    target_2h     FLOAT8,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_raw_data_route ON raw_data(route_id, office_from_id);
CREATE INDEX IF NOT EXISTS idx_raw_data_timestamp ON raw_data(timestamp);