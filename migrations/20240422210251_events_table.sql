-- +goose Up
-- +goose StatementBegin
CREATE TABLE test.events (
    eventID Int64,
    eventType String,
    userID Int64,
    eventTime DateTime,
    payload String
) ENGINE = MergeTree
ORDER BY (eventID, eventTime);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE test.events;
-- +goose StatementEnd
