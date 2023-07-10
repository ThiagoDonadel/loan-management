CREATE TABLE IF NOT EXISTS loans
(
    id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    method integer,
    loan_value numeric(32,2),
    rate numeric(16,8),
    rate_base_months integer,
    term integer,
    start_date date,
    sign_date timestamp without time zone
);

CREATE TABLE IF NOT EXISTS loan_values (
	id bigserial NOT NULL PRIMARY KEY,
	loan_id uuid NOT NULL REFERENCES loans(id),
	installment_number int,
	payment_date date,
	installment numeric(32,2),
	interest numeric(32,2),
	amortization numeric(32,2),
	balance numeric(32,2)
);