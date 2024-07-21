-- +goose Up
CREATE TABLE proveedores  (
    nombre VARCHAR(100) NOT NULL,
    codigo VARCHAR(20) NOT NULL
);

-- +goose Down
DROP TABLE proveedores;
