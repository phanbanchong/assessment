CREATE TABLE IF NOT EXISTS expenses ( id SERIAL PRIMARY KEY, title TEXT, amount FLOAT, note TEXT, tags TEXT[]);
INSERT INTO expenses (id, title, amount, note, tags) VALUES(1,'test-title1', 13, 'test-note1', ARRAY['tag1', 'tag2']);
INSERT INTO expenses (id, title, amount, note, tags) VALUES(2,'test-title2', 14, 'test-note2', ARRAY['tag3', 'tag4']);
ALTER TABLE expenses ALTER COLUMN id SET DEFAULT 3;