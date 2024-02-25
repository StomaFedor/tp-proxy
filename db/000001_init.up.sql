CREATE TABLE "request"
(
    id            SERIAL PRIMARY KEY,
    data          JSON   NOT NULL
);

CREATE TABLE "responce"
(
    id            SERIAL PRIMARY KEY,
    data          JSON   NOT NULL
);
