CREATE SCHEMA movie_service;

CREATE TABLE IF NOT EXISTS movie_service.movies ( 
    id bigserial PRIMARY KEY,
    name varchar,
    description varchar,
    path varchar
);

CREATE INDEX partner_index ON movie_service.movies USING btree (name);

INSERT INTO movie_service.movies (name, description) VALUES ('Матрица', 'Фантастическое кино');
INSERT INTO movie_service.movies (name, description) VALUES ('Терминатор', 'Фантастическое кино');
INSERT INTO movie_service.movies (name, description) VALUES ('Пираты карибского моря', 'Приключения');
INSERT INTO movie_service.movies (name, description) VALUES ('Побег из шаушенко', 'Драма');