-- DROP TRIGGER IF EXISTS users_updated_at ON tasks;
DROP TRIGGER IF EXISTS exhibits_updated_at ON exhibits;
DROP TRIGGER IF EXISTS users_updated_at ON users;
DROP FUNCTION IF EXISTS set_updated_at();

DROP TABLE IF EXISTS exhibits_exhibitions;
DROP TABLE IF EXISTS exhibitions;
DROP TABLE IF EXISTS category;
-- DROP TABLE IF EXISTS material;
DROP TABLE IF EXISTS exhibits;
DROP TABLE IF EXISTS users;