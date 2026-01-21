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

CREATE TABLE IF NOT EXISTS nextbus.routes
(
    id uuid NOT NULL,
    "route_no" text COLLATE pg_catalog."default" NOT NULL,
    company text COLLATE pg_catalog."default" NOT NULL,
    bound text COLLATE pg_catalog."default" NOT NULL,
    "service_type" smallint NOT NULL,
    "orig_en" text COLLATE pg_catalog."default" NOT NULL,
    "orig_sc" text COLLATE pg_catalog."default",
    "orig_tc" text COLLATE pg_catalog."default",
    "dest_en" text COLLATE pg_catalog."default" NOT NULL,
    "dest_sc" text COLLATE pg_catalog."default",
    "dest_tc" text COLLATE pg_catalog."default",
    CONSTRAINT routes_pkey PRIMARY KEY (id)
)
