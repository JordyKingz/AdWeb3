CREATE TABLE chain_registrations(order_id uuid PRIMARY KEY REFERENCES orders(id),chain_id bigint NOT NULL,contract text NOT NULL,registered_at timestamptz NOT NULL DEFAULT now());
