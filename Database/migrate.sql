BEGIN;


CREATE TABLE IF NOT EXISTS public."Cities"
(
    city_name text NOT NULL,
    country text NOT NULL,
    lat numeric NOT NULL,
    lon numeric NOT NULL,
    city_id bigserial NOT NULL,
    PRIMARY KEY (city_name, country),
    UNIQUE (city_id),
    UNIQUE (city_name, country)
);

CREATE TABLE IF NOT EXISTS public."Weather"
(
    date bigint NOT NULL,
    temperature numeric,
    date_txt text,
    city_id bigint NOT NULL,
    PRIMARY KEY (city_id, date)
);

ALTER TABLE IF EXISTS public."Weather"
    ADD FOREIGN KEY (city_id)
    REFERENCES public."Cities" (city_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;

END;