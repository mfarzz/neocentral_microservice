-- auth-service/migrations/000001_create_users.down.sql
-- Rollback: drop all auth tables in reverse dependency order.

DROP TABLE IF EXISTS lecturers;
DROP TABLE IF EXISTS students;
DROP TABLE IF EXISTS user_has_roles;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS users;
