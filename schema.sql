CREATE TABLE account  (
    id integer PRIMARY KEY
    uuid UUID UNIQUE DEFAULT gen_random_uuid(),
    role role NOT NULL DEFAULT 'guest',
    create_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
)

CREATE TABLE go_job  (
    account integer,
    slug    text NOT NULL,
    nr      integer NULL,
    von     date NOT NULL,
    bis     date NOT NULL,
    PRIMARY KEY (account, slug)
)

CREATE TABLE jahr (
    id integer
)


CREATE TABLE datum (
    id integer,
    datum date,
    feiertag text,
)