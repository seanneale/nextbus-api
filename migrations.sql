CREATE TABLE IF NOT EXISTS nextbus.stops
(
    name text COLLATE pg_catalog."default" NOT NULL,
    "kmbStopId" text COLLATE pg_catalog."default",
    "kmbNameEn" text COLLATE pg_catalog."default",
    "kmbNameTc" text COLLATE pg_catalog."default",
    "kmbNameSc" text COLLATE pg_catalog."default",
    latitude numeric,
    longitude numeric,
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    CONSTRAINT stops_pkey PRIMARY KEY (id)
)
