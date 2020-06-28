-- +goose Up
CREATE TABLE public.events
(
    id uuid NOT NULL,
	title character varying(100) COLLATE pg_catalog."default" NOT NULL,
    datetime timestamp without time zone NOT NULL,
	duration integer,
    text text COLLATE pg_catalog."default" NOT NULL,
	userid uuid NOT NULL,
	timenotify integer,
    CONSTRAINT "PK_Events" PRIMARY KEY (id)
)

WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE public.events
    OWNER to igor;

CREATE UNIQUE INDEX "Ix_EventsDate"
    ON public.events USING btree
    (date_part('year'::text, datetime) ASC NULLS LAST, date_part('month'::text, datetime) ASC NULLS LAST, date_part('day'::text, datetime) ASC NULLS LAST, date_part('hour'::text, datetime) ASC NULLS LAST, date_part('minute'::text, datetime) ASC NULLS LAST, date_part('second'::text, datetime) ASC NULLS LAST)
    TABLESPACE pg_default;
-- +goose Down
DROP table public.events;
