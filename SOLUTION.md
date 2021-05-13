# Github GO Go2
## Pre requisitos
### Base de datos
El proyecto usa una base de datos sqlite3, con una unica table donde guarda toda la informacion de cada issue

El nombre de la base de datos es: goRepoDB.db
El query para generar la tabla es:

CREATE TABLE "issues" (
	"numero"	INTEGER NOT NULL,
	"url"	TEXT,
	"nombre"	TEXT,
	"autor"	TEXT,
	"mile_nombre"	TEXT,
	"mile_desc"	TEXT,
	"tags"	TEXT,
	PRIMARY KEY("numero")
);

### go.mod
module github.com/CarlosTrejo2308/GoChallenge

go 1.16

require github.com/mattn/go-sqlite3 v1.14.7 // indirect


## Solucion
El proyecto consiste en un archivo main.go el cual se empieza ejecutando la funcion main

El flujo del programa es el siguiente:
1. Se hace una peticion http a la api de github para obtener los issues con el tag de "Go2" del repositorio de GO
2. Itereamos por cada issue y manipulamos la respuesta que nos regresa la api para guardar solo los datos que nos interesa
3. Guardamos los datos de cada issue en la base de datos
4. Leemos la base de datos para imprimir cada issue en un formato deseado

## Mejoras
Ademas de las pruebas unitarias y de sistema, estas son algunas mejoras que noto que se podrian hacer al codigo, pero por falta de tiempo todavia no las implemento
1. Separar las funciones del proyecto en diferentes archivos
- Main: Donde hara el llamado a los demas programas
- Api: La que hara la conexion a la api de github
- Datos: El que va a filtar y guardar los datos de la respuesta que regresa la api
- Base de datos: El encargado de escribir y guardar en la base de datos
2. En la linea numero 245:
url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues?labels=%s&page=1&per_page=100", user, repo, label)
solo obtenemos 100 issues, agregaria un ciclo para iterar y obtener mas issues
3. Mejorar el algoritmo para obtener los datos de tags y de milestone, ya que por ahora esta iterando por cada variable
4. A la hora de guardar datos, debio a que el numero de issue es PK, esta lo rechaza si es que se intenta guardar otraa vez: agregaria un check para asegurarme que no exista un issiue que queramos guardar
5. Mejoraria y normalizare la base de datos. Ya que por ahora contamos con datos reduentantes, y los tags los pondria en una tabla aparte

### Query planeados para la BD:
CREATE TABLE "milestone" (
	"id"	INTEGER NOT NULL,
	"name"	TEXT,
	"desc"	TEXT,
	PRIMARY KEY("id")
);


CREATE TABLE "issue" (
	"id"	INTEGER NOT NULL,
	"url"	TEXT,
	"name"	TEXT,
	"author"	TEXT,
	"fk_milestone"	INTEGER,
	PRIMARY KEY("id"),
	FOREIGN KEY("fk_milestone") REFERENCES "milestone"("id") ON DELETE SET NULL
);

CREATE TABLE "tags" (
	"id"	INTEGER NOT NULL,
	"name"	TEXT,
	PRIMARY KEY("id")
);

CREATE TABLE "hasTags" (
	"fk_issue"	INTEGER,
	"fk_tag"	INTEGER,
	FOREIGN KEY("fk_issue") REFERENCES "issue"("id"),
	FOREIGN KEY("fk_tag") REFERENCES "tags"("id")
);