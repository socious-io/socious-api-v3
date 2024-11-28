CREATE TYPE project_kind AS ENUM ('JOB', 'SERVICE');

ALTER TABLE projects
    ADD COLUMN kind project_kind DEFAULT 'JOB' NOT NULL,
    ADD COLUMN service_total_hours integer,
    ADD COLUMN service_price integer;
