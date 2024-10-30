CREATE TABLE mastercat (
    id text not null,
    ipix bigint not null,
    ra double precision not null,
    dec double precision not null,
    cat text not null,
    PRIMARY KEY (id, cat)
);

CREATE INDEX mastercat_ipix_idx ON mastercat (ipix);
