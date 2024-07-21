package main

import (
	"database/sql"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"piezas-proveedores/.gen/proveedores/public/model"
	"piezas-proveedores/.gen/proveedores/public/table"
)

func main() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")

	}

	dbUrl := os.Getenv("DB_URL")

	if dbUrl == "" {
		log.Fatal("DB_URL not found")
	}

	conn, err := sql.Open("postgres", dbUrl)

	if err != nil {
		log.Fatal("Error connecting to database")

	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world!"))
	})

	r.Post("/pieza/add", func(w http.ResponseWriter, r *http.Request) {

		var pieza model.Piezas

		err := json.NewDecoder(r.Body).Decode(&pieza)

		if err != nil {
			return
		}

		smt := table.Piezas.INSERT(table.Piezas.Nombre, table.Piezas.PiezaID).MODEL(pieza)

		err = smt.Query(conn, &pieza)

		if err != nil {
			return
		}

		_, err = smt.Exec(conn)

		if err != nil {
			return

		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Pieza añadida"))

		json.NewEncoder(w).Encode(pieza)

	})

	r.Post("/proveedor/add", func(w http.ResponseWriter, r *http.Request) {

		var proveedor model.Proveedores

		err := json.NewDecoder(r.Body).Decode(&proveedor)

		if err != nil {
			return
		}

		smt := table.Proveedores.INSERT(table.Proveedores.Nombre, table.Proveedores.Codigo, table.Proveedores.ProveedorID).MODEL(proveedor)

		err = smt.Query(conn, &proveedor)

		if err != nil {
			return
		}

		_, err = smt.Exec(conn)

		if err != nil {
			return

		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Proveedor añadido"))

		json.NewEncoder(w).Encode(proveedor)
	})

	r.Post("/precios/add", func(w http.ResponseWriter, r *http.Request) {
		var precio model.Precios

		err := json.NewDecoder(r.Body).Decode(&precio)

		if err != nil {
			return
		}

		smt := table.Precios.INSERT(table.Precios.Precio, table.Precios.PiezaID, table.Precios.ProveedorID).MODEL(precio)

		err = smt.Query(conn, &precio)

		if err != nil {
			return
		}

		_, err = smt.Exec(conn)

		if err != nil {
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Precio añadido"))
	})

	r.Post("/piezas-proveedores/add", func(w http.ResponseWriter, r *http.Request) {
		var piezaProveedor model.PiezasProveedores

		err := json.NewDecoder(r.Body).Decode(&piezaProveedor)

		if err != nil {
			return
		}

		smt := table.PiezasProveedores.INSERT(table.PiezasProveedores.PiezaID, table.PiezasProveedores.ProveedorID).MODEL(piezaProveedor)

		err = smt.Query(conn, &piezaProveedor)

		if err != nil {
			return
		}

		_, err = smt.Exec(conn)

		if err != nil {
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Pieza-Proveedor añadido"))

	})

	r.Get("/consulta/1", func(w http.ResponseWriter, r *http.Request) {
		var piezas []model.Piezas

		smt := table.Piezas.SELECT(table.Piezas.AllColumns).FROM(table.Piezas)

		err := smt.Query(conn, &piezas)

		if err != nil {
			return
		}

		_, err = smt.Exec(conn)

		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(piezas)

	})

	r.Get("/consulta/2", func(w http.ResponseWriter, r *http.Request) {
		var proveedores []model.Proveedores

		smt := table.Proveedores.SELECT(table.Proveedores.AllColumns).FROM(table.Proveedores)

		err := smt.Query(conn, &proveedores)

		if err != nil {
			return
		}

		_, err = smt.Exec(conn)

		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(proveedores)

	})

	r.Get("/consulta/3", func(w http.ResponseWriter, r *http.Request) {
		var precios struct {
			Avg float64 `db:"avg" json:"avg"`
		}
		// SELECT AVG(precio) FROM precios;
		smt := table.Precios.SELECT(postgres.AVG(table.Precios.Precio).AS("avg")).FROM(table.Precios)

		err := smt.Query(conn, &precios)

		if err != nil {
			return
		}

		_, err = smt.Exec(conn)

		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(precios)
	})

	r.Get("/consulta/4", func(w http.ResponseWriter, r *http.Request) {

		var piezasProveedores []struct {
			Nombre string `db:"nombre" json:"nombre"`
		}

		smt := table.Proveedores.
			SELECT(table.Proveedores.Nombre.AS("nombre")).
			FROM(
				table.Proveedores.
					INNER_JOIN(
						table.Piezas,
						table.Proveedores.ProveedorID.EQ(table.Piezas.PiezaID),
					),
			).
			WHERE(table.Piezas.PiezaID.EQ(postgres.Int(1)))

		err := smt.Query(conn, &piezasProveedores)

		if err != nil {
			return
		}

		_, err = smt.Exec(conn)

		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(piezasProveedores)
	})

	r.Get("/consulta/5", func(w http.ResponseWriter, r *http.Request) {

		// Parametro: Codigo del proveedor

		var piezas []struct {
			Nombre string `db:"nombre" json:"nombre"`
		}

		var smt = table.Piezas.
			SELECT(table.Piezas.Nombre.AS("nombre")).
			FROM(
				table.Piezas.
					INNER_JOIN(
						table.Proveedores,
						table.Piezas.PiezaID.EQ(table.Proveedores.ProveedorID),
					),
			).
			WHERE(table.Proveedores.Codigo.EQ(postgres.String("gfdsghfgdh")))

		err = smt.Query(conn, &piezas)

		if err != nil {
			return
		}

		_, err = smt.Exec(conn)

		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(piezas)

	})

	// consulta 6

	r.Get("/consulta/6", func(w http.ResponseWriter, r *http.Request) {

		var Resultado struct {
			NombrePieza     string `db:"nombre_pieza"`
			NombreProveedor string `db:"nombre_proveedor"`
			Precio          int    `db:"precio"`
		}

		var smt = table.Piezas.
			SELECT(table.Piezas.Nombre.AS("nombre_pieza"), table.Precios.Precio.AS("precio"), table.Proveedores.Nombre.AS("nombre_proveedor")).
			FROM(
				table.Piezas.
					INNER_JOIN(
						table.PiezasProveedores,
						table.Piezas.PiezaID.EQ(table.PiezasProveedores.PiezaID),
					).
					INNER_JOIN(
						table.Precios,
						table.PiezasProveedores.PiezaID.EQ(table.Precios.PiezaID).AND(
							table.PiezasProveedores.ProveedorID.EQ(table.Precios.ProveedorID),
						),
					).
					INNER_JOIN(
						table.Proveedores,
						table.PiezasProveedores.ProveedorID.EQ(table.Proveedores.ProveedorID),
					),
			).
			ORDER_BY(table.Precios.Precio.DESC())

		err := smt.Query(conn, &Resultado)

		if err != nil {
			return
		}

		_, err = smt.Exec(conn)

		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(Resultado)

	})

	http.ListenAndServe(":3000", r)
}
