-- +goose Up
CREATE TABLE piezas (
    nombre VARCHAR(100) NOT NULL
);

-- +goose Down
DROP TABLE piezas;
