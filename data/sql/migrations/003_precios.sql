-- +goose Up
CREATE TABLE precios (
    precio  FLOAT NOT NULL,
    pieza_id INT NOT NULL,
    proveedor_id INT NOT NULL
);

-- +goose Down
DROP TABLE precios;
