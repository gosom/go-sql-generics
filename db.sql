CREATE TABLE notes (
  id SERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  content TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

INSERT INTO notes
    (title, content)
    VALUES
    ('title1', 'some data...'),
    ('title2', 'some data2...'),
    ('title3', 'some data3...'),
    ('title4', 'some data4...'),
    ('title5', 'some data5...'),
    ('title6', 'some data6...'),
    ('title7', 'some data7...'),
    ('title8', 'some data8...');
