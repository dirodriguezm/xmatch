CREATE TABLE IF NOT EXISTS mastercat (
    id text not null,
    ipix bigint not null,
    ra double precision not null,
    dec double precision not null,
    cat text not null,
    PRIMARY KEY (id, cat)
);

CREATE INDEX IF NOT EXISTS mastercat_ipix_idx ON mastercat (ipix);
