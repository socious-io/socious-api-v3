CREATE TYPE project_kind AS ENUM ('JOB', 'SERVICE');

ALTER TABLE projects
    ADD COLUMN kind project_kind DEFAULT 'JOB' NOT NULL,
    ADD COLUMN service_total_hours integer,
    ADD COLUMN service_price integer;

CREATE TABLE service_work_samples (
  id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
  service_id UUID NOT NULL,
  document UUID NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  CONSTRAINT fk_service FOREIGN KEY (service_id) REFERENCES projects(id) ON DELETE CASCADE,
  CONSTRAINT fk_media FOREIGN KEY (document) REFERENCES media(id) ON DELETE CASCADE
);