CREATE TABLE account  (
    id integer PRIMARY KEY,
    -- uuid UUID UNIQUE DEFAULT gen_random_uuid(),
    -- role role NOT NULL DEFAULT 'guest',
    create_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE c11_job  (
    code     varchar(8) PRIMARY KEY,
    account  integer NOT NULL,
    nr       integer NOT NULL,
    eintritt date NOT NULL,
    austritt date NOT NULL,
    UNIQUE (code, account),
    UNIQUE (account, nr)
);


CREATE TABLE c11_urlaub  (
    account integer NOT NULL,
    job     varchar(8) NOT NULL, -- REFERENCES go_job()
    jahr    integer NOT NULL,

    nr   integer NOT NULL,
    von  date    NOT NULL,
    bis  date    NOT NULL,

    num_urlaub float NOT NULL,
    num_ausgl  float NOT NULL,
    num_sonder float NOT NULL,

    grund      text NOT NULL,
    beantragt  date NOT NULL,
    genehmigt  date NOT NULL,
    kommentar  text NOT NULL,

    PRIMARY KEY (account, job, jahr, nr),
    -- FOREIGN KEY (account) REFERENCES account (account),
    FOREIGN KEY (account,job) REFERENCES c11_job (account, code),
    FOREIGN KEY (account,job,jahr) REFERENCES c11_arbeitsjahr (account,job,jahr)
);
 


--CREATE TABLE go_kalendertag  (

CREATE TABLE c11_datum (
    id       integer PRIMARY KEY,
    datum    date NOT NULL UNIQUE,
    jahr     integer NOT NULL,
    monat    integer NOT NULL,
    tag      integer NOT NULL,
    jahrtag  integer NOT NULL,
    kw_jahr  integer NOT NULL,
    kw       integer NOT NULL,
    kw_tag   integer NOT NULL,
    feiertag text NOT NULL DEFAULT '',
    UNIQUE (jahr, monat, tag),
    UNIQUE (kw_jahr, kw, kw_tag),
    UNIQUE (jahr, jahrtag)
);


CREATE TABLE c11_arbeitsjahr  (
    account integer NOT NULL,
    job     varchar(8) NOT NULL, -- REFERENCES go_job()
    jahr    integer NOT NULL,

    urlaub_vorjahr        NUMERIC(5, 2) NOT NULL DEFAULT 0,
    urlaub_anspruch       NUMERIC(5, 2) NOT NULL DEFAULT 0,
    urlaub_tage           NUMERIC(5, 2) NOT NULL DEFAULT 0,
    urlaub_geplant        NUMERIC(5, 2) NOT NULL DEFAULT 0,
    urlaub_rest           NUMERIC(5, 2) NOT NULL DEFAULT 0,
    
    soll             float NOT NULL DEFAULT 0,
    ist              float NOT NULL DEFAULT 0,
    diff             float NOT NULL DEFAULT 0,
    PRIMARY KEY (account,job,jahr)
);

CREATE TABLE c11_arbeitsmonat  (
    account integer NOT NULL,
    job     varchar(8) NOT NULL,
    jahr    integer NOT NULL,
    monat   integer NOT NULL,
    soll    float NOT NULL DEFAULT 0,
    ist     float NOT NULL DEFAULT 0,
    diff    float NOT NULL DEFAULT 0,
    PRIMARY KEY (account,job,jahr,monat),
    FOREIGN KEY (account,job,jahr) REFERENCES c11_arbeitsjahr (account,job,jahr)
);

CREATE TABLE c11_arbeitstag  (
    account integer,
    job     varchar(8) NOT NULL,
    datum   date NOT NULL,
    
    jahr    integer NOT NULL,
    monat   integer NOT NULL,
    --woche integer NOT NULL,

    status       char(1) NOT NULL,
    kategorie    char(1) NOT NULL,
    -- krankmeldung | boolean                  |           | not null | false
    -- urlaubstage       NUMERIC(5, 2) NOT NULL,
    -- freizeitausgleich boolean,  -- | double precision
    -- krank             | boolean                  |           | not null | 
    -- krankheit         | text  
    soll         float NOT NULL,
    start        timestamp with time zone NULL,
    ende         timestamp with time zone NULL,
    brutto       float NOT NULL,
    pausen       float NOT NULL,
    extra        float NOT NULL,
    netto        float NOT NULL,
    diff         float NOT NULL,
    
    -- saldo        float NOT NULL,

    kommentar    text  NOT NULL,

    PRIMARY KEY (account, datum),
    FOREIGN KEY (account,job) REFERENCES c11_job (account,code),
    FOREIGN KEY (account,job,jahr) REFERENCES c11_arbeitsjahr (account,job,jahr)
);


CREATE TABLE c11_zeitspanne (
    account integer,
 -- job     varchar(8) NOT NULL,
    datum   date                     NOT NULL,
    nr      integer                  NOT NULL,
    status  varchar(1)               NOT NULL,
    start   timestamp with time zone,
    ende    timestamp with time zone,
    dauer   float                    NOT NULL DEFAULT 0,
 -- titel         character varying(100) NOT NULL,
 -- story         character varying(100) NOT NULL,
 -- beschreibung  text  NOT NULL,
 -- grund         character varying(255)  NOT NULL,
 -- arbeitszeit   boolean  NOT NULL,
    PRIMARY KEY (account, datum, nr),
    FOREIGN KEY (account, datum) REFERENCES c11_arbeitstag (account, datum)
);
              