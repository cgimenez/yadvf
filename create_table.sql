-- DROP TABLE dvf;
CREATE TABLE dvf
(
	id_mutation TEXT,
	date_mutation TEXT,
	numero_disposition INTEGER,
	nature_mutation TEXT,
	valeur_fonciere REAL,
	adresse_numero INTEGER,
	adresse_suffixe TEXT,
	adresse_nom_voie TEXT,
	adresse_code_voie TEXT,
	code_postal TEXT,
	code_commune TEXT,
	nom_commune TEXT,
	code_departement TEXT,
	ancien_code_commune TEXT,
	ancien_nom_commune TEXT,
	id_parcelle TEXT,
	ancien_id_parcelle TEXT,
	numero_volume TEXT,
	lot1_numero TEXT,
	lot1_surface_carrez REAL,
	lot2_numero TEXT,
	lot2_surface_carrez REAL,
	lot3_numero TEXT,
	lot3_surface_carrez REAL,
	lot4_numero TEXT,
	lot4_surface_carrez REAL,
	lot5_numero TEXT,
	lot5_surface_carrez REAL,
	nombre_lots INTEGER,
	code_type_local TEXT,
	type_local TEXT,
	surface_reelle_bati REAL,
	nombre_pieces_principales INTEGER,
	code_nature_culture TEXT,
	nature_culture TEXT,
	code_nature_culture_speciale TEXT,
	nature_culture_speciale TEXT,
	surface_terrain REAL,
	longitude REAL,
	latitude REAL
);

CREATE INDEX idx_date_mutation ON dvf(date_mutation);
CREATE INDEX idx_nom_commune ON dvf(nom_commune);
CREATE INDEX idx_code_postal ON dvf(code_postal);
CREATE INDEX idx_code_departement ON dvf(code_departement);
CREATE INDEX idx_nature_culture ON dvf(nature_culture);
CREATE INDEX idx_type_local ON dvf(type_local);
CREATE INDEX idx_nature_culture_speciale ON dvf(nature_culture_speciale);