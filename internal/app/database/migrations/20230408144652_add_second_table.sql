-- +goose Up
-- +goose StatementBegin

CREATE TABLE groups(
    id BIGSERIAL Primary key not null,
    amount_of_students int CHECK(amount_of_students > 0) DEFAULT 30 NOT NULL,
    group_name VARCHAR(30) DEFAULT NULL,
    st_year int CHECK(st_year >= 1 and st_year <= 4) DEFAULT 1 NOT NULL --курс обучения (1-4)
);

ALTER TABLE students
ADD group_id BIGINT references groups(id),;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE groups;

ALTER TABLE students
DROP COLUMN group_id;
-- +goose StatementEnd
