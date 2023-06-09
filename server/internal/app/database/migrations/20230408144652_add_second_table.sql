-- +goose Up
-- +goose StatementBegin

CREATE TABLE groups(
    id BIGSERIAL Primary key not null,
    group_name VARCHAR(30) DEFAULT NULL,
    st_year int CHECK(st_year >= 1 and st_year <= 4) DEFAULT 1 NOT NULL --курс обучения (1-4)
);

ALTER TABLE students
ADD group_id BIGINT references groups(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE students
DROP COLUMN group_id;
DROP TABLE groups;
-- +goose StatementEnd
