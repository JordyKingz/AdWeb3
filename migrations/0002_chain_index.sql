ALTER TABLE orders ADD COLUMN IF NOT EXISTS chain_order_id text UNIQUE;
ALTER TABLE refund_intents ADD COLUMN IF NOT EXISTS chain_refund_id text UNIQUE;
CREATE TABLE IF NOT EXISTS chain_checkpoints(chain_id bigint NOT NULL,contract text NOT NULL,block_number bigint NOT NULL,block_hash text NOT NULL,PRIMARY KEY(chain_id,contract));
