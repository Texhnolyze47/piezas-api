-- +goose Up
ALTER TABLE piezas ADD COLUMN pieza_id INT NOT NULL;
ALTER TABLE proveedores ADD COLUMN proveedor_id INT NOT NULL;

-- +goose Down
ALTER TABLE piezas DROP COLUMN pieza_id;
ALTER TABLE proveedores DROP COLUMN proveedor_id;
