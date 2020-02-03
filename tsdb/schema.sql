CREATE DATABASE factorio;

\c factorio;

CREATE USER grafana WITH PASSWORD 'grafana';

CREATE TABLE circuit_network_signals (
  tick        INT         NOT NULL,
  time        TIMESTAMPTZ NOT NULL,
  network_id  INT         NOT NULL,
  signal_type TEXT        NOT NULL,
  signal_name TEXT        NOT NULL,
  count       INT         NOT NULL
);

-- 60 ticks per bucket
-- 1 bucket per second
-- 2 KB per bucket (1 circuit network, 1 signal)
-- 86400 seconds per day
-- ~= 10.368 GB/day
SELECT create_hypertable('circuit_network_signals', 'tick', chunk_time_interval => 648000);

GRANT SELECT ON circuit_network_signals TO grafana;

CREATE USER frpc WITH PASSWORD 'frpc';

GRANT INSERT ON circuit_network_signals TO frpc;
