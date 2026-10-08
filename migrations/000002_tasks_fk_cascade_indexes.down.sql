DROP INDEX IF EXISTS todoapp.tasks_created_at_idx;
DROP INDEX IF EXISTS todoapp.tasks_author_user_id_idx;

ALTER TABLE todoapp.tasks
	ALTER COLUMN created_at DROP DEFAULT;

ALTER TABLE todoapp.tasks
	DROP CONSTRAINT IF EXISTS tasks_author_user_id_fkey,
	ADD CONSTRAINT tasks_author_user_id_fkey
		FOREIGN KEY (author_user_id) REFERENCES todoapp.users(id);
