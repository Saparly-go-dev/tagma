-- One-time migration for money values.
-- Converts stored manat values such as 12.34 into integer tenne values: 1234.
-- Run only on databases that still store money as fractional numeric/float values.

BEGIN;

DO $$
BEGIN
	IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'public' AND table_name = 'products' AND column_name = 'price') THEN
		ALTER TABLE products ALTER COLUMN price TYPE bigint USING ROUND(price * 100)::bigint;
	END IF;
	IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'public' AND table_name = 'products' AND column_name = 'd_price') THEN
		ALTER TABLE products ALTER COLUMN d_price TYPE bigint USING ROUND(d_price * 100)::bigint;
	END IF;

	IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'public' AND table_name = 'dealer_products' AND column_name = 'price') THEN
		ALTER TABLE dealer_products ALTER COLUMN price TYPE bigint USING ROUND(price * 100)::bigint;
	END IF;

	IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'public' AND table_name = 'orders' AND column_name = 'sum') THEN
		ALTER TABLE orders ALTER COLUMN sum TYPE bigint USING ROUND(sum * 100)::bigint;
	END IF;
	IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'public' AND table_name = 'orders' AND column_name = 'sale_sum') THEN
		ALTER TABLE orders ALTER COLUMN sale_sum TYPE bigint USING ROUND(sale_sum * 100)::bigint;
	END IF;

	IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'public' AND table_name = 'order_lists' AND column_name = 'price') THEN
		ALTER TABLE order_lists ALTER COLUMN price TYPE bigint USING ROUND(price * 100)::bigint;
	END IF;
	IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'public' AND table_name = 'order_lists' AND column_name = 'sale_price') THEN
		ALTER TABLE order_lists ALTER COLUMN sale_price TYPE bigint USING ROUND(sale_price * 100)::bigint;
	END IF;

	IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'public' AND table_name = 'payments' AND column_name = 'currency') THEN
		ALTER TABLE payments ALTER COLUMN currency TYPE bigint USING ROUND(currency * 100)::bigint;
	END IF;

	IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'public' AND table_name = 'debts' AND column_name = 'total') THEN
		ALTER TABLE debts ALTER COLUMN total TYPE bigint USING ROUND(total * 100)::bigint;
	END IF;

	IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'public' AND table_name = 'dealers' AND column_name = 'account') THEN
		ALTER TABLE dealers ALTER COLUMN account TYPE bigint USING ROUND(account * 100)::bigint;
	END IF;
	IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'public' AND table_name = 'dealers' AND column_name = 'debt') THEN
		ALTER TABLE dealers ALTER COLUMN debt TYPE bigint USING ROUND(debt * 100)::bigint;
	END IF;

	IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'public' AND table_name = 'dealer_orders' AND column_name = 'sum') THEN
		ALTER TABLE dealer_orders ALTER COLUMN sum TYPE bigint USING ROUND(sum * 100)::bigint;
	END IF;
	IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'public' AND table_name = 'dealer_orders' AND column_name = 'sale_sum') THEN
		ALTER TABLE dealer_orders ALTER COLUMN sale_sum TYPE bigint USING ROUND(sale_sum * 100)::bigint;
	END IF;

	IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'public' AND table_name = 'dealer_order_lists' AND column_name = 'price') THEN
		ALTER TABLE dealer_order_lists ALTER COLUMN price TYPE bigint USING ROUND(price * 100)::bigint;
	END IF;
	IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'public' AND table_name = 'dealer_order_lists' AND column_name = 'sale_price') THEN
		ALTER TABLE dealer_order_lists ALTER COLUMN sale_price TYPE bigint USING ROUND(sale_price * 100)::bigint;
	END IF;
END $$;

COMMIT;
