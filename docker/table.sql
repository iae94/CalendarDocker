-- Table: public."Events"

-- DROP TABLE public."Events";

CREATE TABLE public."Events"
(
    id serial NOT NULL ,
    "UUID" uuid,
    "Summary" text COLLATE pg_catalog."default",
    "Description" text COLLATE pg_catalog."default",
    "User" text COLLATE pg_catalog."default",
    "StartDate" timestamp without time zone,
    "EndDate" timestamp without time zone,
    "NotifyTime" timestamp without time zone,
    CONSTRAINT "Events_pkey" PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE public."Events"
    OWNER to api_user;

GRANT ALL ON TABLE public."Events" TO api_user;
GRANT ALL ON TABLE public."Events" TO scheduler_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO api_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO scheduler_user;