-- +goose Up
-- +goose StatementBegin
ALTER TABLE hotel
ADD COLUMN manager_id INT REFERENCES my_user (id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE hotel
DROP COLUMN manager_id;
-- +goose StatementEnd
