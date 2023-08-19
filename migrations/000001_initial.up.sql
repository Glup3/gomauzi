BEGIN;
CREATE TABLE IF NOT EXISTS public.user (
  id SERIAL PRIMARY KEY,
  email TEXT NOT NULL,
  username TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS public.currency (
  id SERIAL PRIMARY KEY,
  code VARCHAR(3) NOT NULL
);
CREATE TABLE IF NOT EXISTS public.account (
  id SERIAL PRIMARY KEY,
  account_name TEXT NOT NULL,
  user_id INTEGER NOT NULL,
  CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES public.user(id)
);
CREATE TABLE IF NOT EXISTS public.category (
  id SERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  created_by INTEGER NOT NULL,
  CONSTRAINT fk_user FOREIGN KEY(created_by) REFERENCES public.user(id)
);
CREATE TABLE IF NOT EXISTS public.trecord (
  id SERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  amount_smallest_unit BIGINT NOT NULL,
  record_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  account_id INTEGER NOT NULL,
  category_id INTEGER NOT NULL,
  currency_id INTEGER NOT NULL,
  CONSTRAINT fk_account FOREIGN KEY(account_id) REFERENCES public.account(id),
  CONSTRAINT fk_category FOREIGN KEY(category_id) REFERENCES public.category(id),
  CONSTRAINT fk_currency FOREIGN KEY(currency_id) REFERENCES public.currency(id)
);
COMMIT;