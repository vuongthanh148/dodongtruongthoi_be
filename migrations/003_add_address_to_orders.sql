-- Add customer address to orders
ALTER TABLE orders
ADD COLUMN IF NOT EXISTS address TEXT;
