-- +goose Up
CREATE TABLE piezas_proveedores (
    pieza_id INT NOT NULL,
    proveedor_id INT NOT NULL
);

-- +goose Down
DROP TABLE piezas_proveedores;
