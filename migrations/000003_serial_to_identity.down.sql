DO $$
DECLARE
	next_id BIGINT;
BEGIN
	EXECUTE format(
		'SELECT GREATEST(CASE WHEN is_called THEN last_value + 1 ELSE last_value END, %s) FROM %s',
		(SELECT COALESCE(MAX(id), 0) + 1 FROM todoapp.tasks),
		pg_get_serial_sequence('todoapp.tasks', 'id')
	) INTO next_id;

	ALTER TABLE todoapp.tasks ALTER COLUMN id DROP IDENTITY IF EXISTS;
	CREATE SEQUENCE todoapp.tasks_id_seq AS INTEGER OWNED BY todoapp.tasks.id;
	PERFORM setval('todoapp.tasks_id_seq', next_id, false);
	ALTER TABLE todoapp.tasks ALTER COLUMN id SET DEFAULT nextval('todoapp.tasks_id_seq');

	EXECUTE format(
		'SELECT GREATEST(CASE WHEN is_called THEN last_value + 1 ELSE last_value END, %s) FROM %s',
		(SELECT COALESCE(MAX(id), 0) + 1 FROM todoapp.users),
		pg_get_serial_sequence('todoapp.users', 'id')
	) INTO next_id;

	ALTER TABLE todoapp.users ALTER COLUMN id DROP IDENTITY IF EXISTS;
	CREATE SEQUENCE todoapp.users_id_seq AS INTEGER OWNED BY todoapp.users.id;
	PERFORM setval('todoapp.users_id_seq', next_id, false);
	ALTER TABLE todoapp.users ALTER COLUMN id SET DEFAULT nextval('todoapp.users_id_seq');
END $$;
