ALTER TABLE todoapp.tasks
	DROP CONSTRAINT tasks_author_user_id_fkey,
	ADD CONSTRAINT tasks_author_user_id_fkey
		FOREIGN KEY (author_user_id) REFERENCES todoapp.users(id) ON DELETE CASCADE;

ALTER TABLE todoapp.tasks
	ALTER COLUMN created_at SET DEFAULT now();

CREATE INDEX tasks_author_user_id_idx ON todoapp.tasks (author_user_id);
CREATE INDEX tasks_created_at_idx ON todoapp.tasks (created_at);
