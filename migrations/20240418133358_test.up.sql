CREATE OR REPLACE FUNCTION make_uid() RETURNS text AS $$
DECLARE
    new_uid text;
    done bool;
BEGIN
    done := false;
    WHILE NOT done LOOP
        new_uid := md5(''||now()::text||random()::text);
        done := NOT exists(SELECT 1 FROM wallets WHERE id=new_uid);
    END LOOP;
    RETURN new_uid;
END;
$$ LANGUAGE PLPGSQL VOLATILE;

CREATE TABLE IF NOT EXISTS wallets
(
	id TEXT DEFAULT make_uid()::text NOT NULL UNIQUE,
	balance INTEGER DEFAULT 0 CHECK (balance >= 0) NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions
(
	time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    from_wallet_id TEXT NOT NULL REFERENCES wallets(id),
    to_wallet_id TEXT NOT NULL REFERENCES wallets(id),
    amount INTEGER NOT NULL CHECK (amount > 0)
);

CREATE OR REPLACE FUNCTION generate_data(num_wallets INTEGER, num_transactions INTEGER)
RETURNS VOID AS $$
DECLARE
    i INTEGER;
    from_wallet TEXT;
    to_wallet TEXT;
    trans_amount INTEGER;
BEGIN
    -- Генерация кошельков
    FOR i IN 1..num_wallets LOOP
        INSERT INTO wallets (id, balance) VALUES (make_uid()::text, floor(random() * 10000 + 1)::INTEGER);
    END LOOP;

    -- Генерация транзакций
    FOR i IN 1..num_transactions LOOP
        -- Выбор случайных кошельков для транзакции
        SELECT id INTO from_wallet FROM wallets ORDER BY random() LIMIT 1;
        SELECT id INTO to_wallet FROM wallets WHERE id != from_wallet ORDER BY random() LIMIT 1;

        -- Вычисление случайной суммы транзакции, не превышающей баланс отправителя
        SELECT balance INTO trans_amount FROM wallets WHERE id = from_wallet;
        trans_amount := floor(random() * trans_amount + 1)::INTEGER;

        -- Выполнение транзакции
        BEGIN
            UPDATE wallets SET balance = balance - trans_amount WHERE id = from_wallet;
            UPDATE wallets SET balance = balance + trans_amount WHERE id = to_wallet;
            INSERT INTO transactions (from_wallet_id, to_wallet_id, amount) VALUES (from_wallet, to_wallet, trans_amount);
        EXCEPTION WHEN check_violation THEN
            -- Если транзакция не удалась из-за недостатка средств, пропустить ее
            CONTINUE;
        END;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

SELECT generate_data(100, 1000);

DROP FUNCTION IF EXISTS generate_data;
