CREATE TABLE tasks_abandoned (
    id INT REFERENCES tasks(id) ON DELETE SET NULL,
    name_task VARCHAR(50) REFERENCES tasks(name) ON DELETE SET NULL,
    content TEXT, 
    abandoned_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);