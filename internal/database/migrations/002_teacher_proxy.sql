ALTER TABLE schedules
  ADD COLUMN teacher_id INT REFERENCES users(id) ON DELETE SET NULL,
  ADD COLUMN proxy_pool INT[] NOT NULL DEFAULT '{}';

CREATE TABLE assignments (
  id SERIAL PRIMARY KEY,
  course_id INT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
  teacher_id INT NOT NULL REFERENCES users(id),
  title VARCHAR(200) NOT NULL,
  description TEXT,
  due_date TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE attendance (
  id SERIAL PRIMARY KEY,
  course_id INT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
  student_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  date DATE NOT NULL,
  status VARCHAR(10) NOT NULL CHECK (status IN ('Present','Absent')),
  UNIQUE (course_id, student_id, date)
);

CREATE TABLE leaves (
  id SERIAL PRIMARY KEY,
  teacher_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  start_date DATE NOT NULL,
  end_date DATE NOT NULL,
  reason TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE proxy_assignments (
  id SERIAL PRIMARY KEY,
  schedule_id INT NOT NULL REFERENCES schedules(id) ON DELETE CASCADE,
  date DATE NOT NULL,
  proxy_teacher_id INT NOT NULL REFERENCES users(id),
  UNIQUE (schedule_id, date)
);

CREATE TABLE notifications (
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  message TEXT NOT NULL,
  is_read BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMPTZ DEFAULT NOW()
);