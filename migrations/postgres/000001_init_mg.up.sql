CREATE TABLE projects (
    id          INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name        TEXT NOT NULL,
    created_at  TIMESTAMP DEFAULT NOW()
);
INSERT INTO projects VALUES (default,'Первая запись',default);

CREATE TABLE goods (
    id           INT GENERATED ALWAYS AS IDENTITY,
    project_id   INT REFERENCES projects (id) ON DELETE CASCADE NOT NULL,
    name         TEXT NOT NULL,
    description  TEXT NOT NULL default '',
    priority     INT NOT NULL,
    removed      BOOL DEFAULT FALSE,
    created_at   TIMESTAMP DEFAULT NOW()
);
CREATE INDEX goods_project_id_name ON goods (project_id,name);

CREATE OR REPLACE FUNCTION increment_priority()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.priority := (SELECT COALESCE(MAX(priority),1) FROM goods) + 1;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_priority
    BEFORE INSERT ON goods
    FOR EACH ROW
EXECUTE PROCEDURE increment_priority();