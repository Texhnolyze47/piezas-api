-- +goose Up


-- +goose Down
ALTER TABLE piezas DROP COLUMN pieza_id;
ALTER TABLE proveedores DROP COLUMN proveedor_id;
