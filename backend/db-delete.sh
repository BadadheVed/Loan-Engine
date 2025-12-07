DB_URL=${1:-"postgres://ved:root@localhost:5432/clickpe"}
psql $DB_URL -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;" -c "\dt"