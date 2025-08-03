CREATE TABLE IF NOT EXISTS public.wallets (
    id SERIAL PRIMARY KEY,
    addr character varying(64) UNIQUE NOT NULL,
    balance NUMERIC(12,2)
);

CREATE TABLE IF NOT EXISTS public.transfers (
    id SERIAL PRIMARY KEY,
    id_from INTEGER NOT NULL REFERENCES public.wallets(id),
    id_to INTEGER NOT NULL REFERENCES public.wallets(id),
    amount NUMERIC(12,2) NOT NULL CHECK (amount > 0),
    created_at TIMESTAMPTZ DEFAULT now()
);