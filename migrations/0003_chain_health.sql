CREATE TABLE chain_health(chain_id bigint PRIMARY KEY,quarantined boolean NOT NULL DEFAULT false,reason text NOT NULL DEFAULT '',updated_at timestamptz NOT NULL DEFAULT now());
