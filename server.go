package main

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type JsonQuery struct {
	Nature_culture          string
	Nature_culture_speciale string
	Type_local              string
	Code_commune            []string
	Code_postal             []string
	Code_departement        string
	Nom_commune             string
	Date_debut              string
	Date_fin                string
	Surface_terrain_min     string
	Surface_terrain_max     string
	Surface_reelle_bati_min string
	Surface_reelle_bati_max string
}

type Mutation struct {
	Date_mutation             string
	Valeur_fonciere           string
	Adresse_numero            string
	Adresse_nom_voie          string
	Code_postal               string
	Nom_commune               string
	Code_commune              string
	Id_parcelle               string
	Type_local                string
	Nature_culture            string
	Nature_culture_speciale   string
	Surface_terrain           string
	Surface_relle_bati        string
	Nombre_pieces_principales string
	Longitude                 string
	Latitude                  string
}

type QueryParam struct {
	clause string
	value  interface{}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logInfo("HTTP %s", r.RequestURI)
		next.ServeHTTP(w, r) // Call the next handler, which can be another middleware in the chain, or the final handler.
	})
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

//
// Send uniques types (eg. Maison, Appartement)
//
func uniq_types_handler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	json_data, err := json.Marshal(db.Types)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(json_data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logError("%s", err.Error())
	}
}

//
// Send mutations to HTTP client
//
func query_handler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	//
	// Fill the JsonQuery from HTTP form request
	//
	var jq JsonQuery
	err := json.NewDecoder(r.Body).Decode(&jq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	qmap := make(map[string]QueryParam)
	qmap["Nature_culture"] = QueryParam{"nature_culture = ?", jq.Nature_culture}
	qmap["Nature_culture_speciale"] = QueryParam{"nature_culture = ?", jq.Nature_culture_speciale}
	qmap["Type_local"] = QueryParam{"type_local = ?", jq.Type_local}
	qmap["Code_commune"] = QueryParam{"code_commune IN (?)", strings.Join(jq.Code_commune[:], ",")}
	qmap["Code_postal"] = QueryParam{"code_postal IN (?)", strings.Join(jq.Code_postal[:], ",")}
	qmap["Nom_commune"] = QueryParam{"nom_commune LIKE ?", "%" + jq.Nom_commune + "%"}
	qmap["Code_departement"] = QueryParam{"code_departement = ?", jq.Code_departement}
	qmap["Date_debut"] = QueryParam{"date_mutation >= ?", jq.Date_debut}
	qmap["Date_fin"] = QueryParam{"date_mutation <= ?", jq.Date_fin}
	qmap["Surface_terrain_min"] = QueryParam{"surface_terrain != \"\" AND surface_terrain >= ?", jq.Surface_terrain_min}
	qmap["Surface_terrain_max"] = QueryParam{"surface_terrain != \"\" AND surface_terrain <= ?", jq.Surface_terrain_max}
	qmap["Surface_reelle_bati_min"] = QueryParam{"surface_reelle_bati != \"\" AND surface_reelle_bati >= ?", jq.Surface_reelle_bati_min}
	qmap["Surface_reelle_bati_max"] = QueryParam{"surface_reelle_bati != \"\" AND surface_reelle_bati <= ?", jq.Surface_reelle_bati_max}

	//
	// Dynamically fill clauses and values slice from json query
	//
	var clauses []string
	var values []interface{}

	st := reflect.ValueOf(&jq).Elem()
	rt := st.Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		rv := reflect.ValueOf(&jq)
		value := reflect.Indirect(rv).FieldByName(field.Name)
		if field.Type.Kind() == reflect.String || field.Type.Kind() == reflect.Slice {
			if value.Len() > 0 {
				clauses = append(clauses, qmap[field.Name].clause)
				values = append(values, qmap[field.Name].value)
			}
		}
	}

	//
	// Build the SQL query
	//
	query_string := `SELECT date_mutation, valeur_fonciere,
	 adresse_numero, adresse_nom_voie, code_postal, nom_commune, code_commune,
	 id_parcelle,
	 type_local, nature_culture, nature_culture_speciale,
	 surface_terrain, surface_reelle_bati, nombre_pieces_principales,
	 longitude, latitude
	 FROM dvf WHERE `
	query_string += strings.Join(clauses[:], " AND ")
	query_string += " ORDER BY id_parcelle ASC, date_mutation DESC"
	logInfo("%s", query_string)

	//
	// Exec query
	//
	var mutations []*Mutation
	rows, err := db.cnx.Query(query_string, values...)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logError("SQL request %s", err.Error())
		return
	}

	defer rows.Close()
	for rows.Next() {
		m := new(Mutation)
		_ = rows.Scan(
			&m.Date_mutation, &m.Valeur_fonciere,
			&m.Adresse_numero, &m.Adresse_nom_voie, &m.Code_postal, &m.Nom_commune, &m.Code_commune,
			&m.Id_parcelle,
			&m.Type_local, &m.Nature_culture, &m.Nature_culture_speciale,
			&m.Surface_terrain, &m.Surface_relle_bati, &m.Nombre_pieces_principales,
			&m.Longitude, &m.Latitude,
		)
		mutations = append(mutations, m)
	}

	//
	// Encode JSON response and send it to the client
	//
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(mutations)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logError("JSON %s", err.Error())
	}
}

//
// Init and start HTTP server
//
func server_start() {
	logInfo("%s", "Starting server")

	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.HandleFunc("/types", uniq_types_handler).Methods("GET")
	r.HandleFunc("/query", query_handler).Methods("POST")
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("www/dist"))))

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
