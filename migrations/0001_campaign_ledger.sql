CREATE TABLE IF NOT EXISTS schema_migrations(version text PRIMARY KEY);
CREATE TABLE IF NOT EXISTS orders(
 id uuid PRIMARY KEY, advertiser text NOT NULL, idempotency_key text NOT NULL,
 fingerprint text NOT NULL, payer text NOT NULL, amount bigint NOT NULL CHECK(amount>0),
 funded bigint NOT NULL DEFAULT 0 CHECK(funded>=0), spent bigint NOT NULL DEFAULT 0 CHECK(spent>=0),
 chain_id bigint NOT NULL CHECK(chain_id IN (31337,43113)), token text NOT NULL, terms text NOT NULL,
 start_at bigint NOT NULL, end_at bigint NOT NULL CHECK(end_at>start_at),
 allocation_version integer NOT NULL DEFAULT 1 CHECK(allocation_version=1),
 CHECK(spent<=funded AND funded<=amount), UNIQUE(advertiser,idempotency_key));
CREATE TABLE IF NOT EXISTS chain_events(id text PRIMARY KEY,fingerprint text NOT NULL,order_id uuid NOT NULL REFERENCES orders(id),kind text NOT NULL,block_number bigint NOT NULL,tx_hash text NOT NULL);
CREATE TABLE IF NOT EXISTS journals(id text PRIMARY KEY,description text NOT NULL,created_at timestamptz NOT NULL DEFAULT now());
CREATE TABLE IF NOT EXISTS ledger_entries(journal_id text NOT NULL REFERENCES journals(id),account text NOT NULL,amount bigint NOT NULL CHECK(amount<>0),PRIMARY KEY(journal_id,account));
CREATE OR REPLACE FUNCTION protect_history() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RAISE EXCEPTION 'immutable financial history'; END $$;
DROP TRIGGER IF EXISTS immutable_entries ON ledger_entries;
CREATE TRIGGER immutable_entries BEFORE UPDATE OR DELETE ON ledger_entries FOR EACH ROW EXECUTE FUNCTION protect_history();
DROP TRIGGER IF EXISTS immutable_journals ON journals;
CREATE TRIGGER immutable_journals BEFORE UPDATE OR DELETE ON journals FOR EACH ROW EXECUTE FUNCTION protect_history();
CREATE OR REPLACE FUNCTION check_journal_balance() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN IF (SELECT COALESCE(sum(amount),0) FROM ledger_entries WHERE journal_id=NEW.journal_id)<>0 THEN RAISE EXCEPTION 'unbalanced journal'; END IF; RETURN NULL; END $$;
DROP TRIGGER IF EXISTS balanced_journal ON ledger_entries;
CREATE CONSTRAINT TRIGGER balanced_journal AFTER INSERT ON ledger_entries DEFERRABLE INITIALLY DEFERRED FOR EACH ROW EXECUTE FUNCTION check_journal_balance();
CREATE TABLE IF NOT EXISTS refund_intents(id uuid PRIMARY KEY,order_id uuid NOT NULL REFERENCES orders(id),actor text NOT NULL,idempotency_key text NOT NULL,amount bigint NOT NULL CHECK(amount>0),status text NOT NULL DEFAULT 'pending' CHECK(status IN ('pending','submitted','confirmed','failed')),UNIQUE(actor,idempotency_key));
CREATE TABLE IF NOT EXISTS outbox(id bigserial PRIMARY KEY,kind text NOT NULL,aggregate_id text NOT NULL,payload jsonb NOT NULL,created_at timestamptz NOT NULL DEFAULT now(),completed_at timestamptz,UNIQUE(kind,aggregate_id));
CREATE TABLE IF NOT EXISTS audit_log(id bigserial PRIMARY KEY,actor text NOT NULL,action text NOT NULL,resource text NOT NULL,created_at timestamptz NOT NULL DEFAULT now());
