CREATE TABLE IF NOT EXISTS public.orders (
    id                UUID        PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    product_name      TEXT        NOT NULL,
    quantity          BIGINT      NOT NULL DEFAULT 0,
    status            TEXT NOT NULL DEFAULT 'pending',
    CONSTRAINT order_status_chk
        CHECK (status IN ('pending', 'completed', 'cancelled'))
);

