-- +goose Up
-- +goose StatementBegin
CREATE TABLE students(
    id BIGSERIAL PRIMARY KEY NOT NULL,
    fisrt_name VARCHAR(30) NOT NULL,
    second_name VARCHAR(30) NOT NULL,
    middle_name VARCHAR(30) DEFAULT NULL,
    gpa FLOAT DEFAULT 0 CHECK(gpa >= 0 and gpa <= 10),
    attendance_rate FLOAT DEFAULT 0 CHECK(attendance_rate >= 0 and attendance_rate <= 1),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE students;
-- +goose StatementEnd
