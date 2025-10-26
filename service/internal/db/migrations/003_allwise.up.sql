CREATE TABLE allwise (
    id text not null,
    cntr bigint not null,
    w1mpro double precision,
    w1sigmpro double precision,
    w2mpro double precision,
    w2sigmpro double precision,
    w3mpro double precision,
    w3sigmpro double precision,
    w4mpro double precision,
    w4sigmpro double precision,
    J_m_2mass double precision,
    J_msig_2mass double precision,
    H_m_2mass double precision,
    H_msig_2mass double precision,
    K_m_2mass double precision,
    K_msig_2mass double precision,
    PRIMARY KEY (id)
);

CREATE INDEX allwise_cntr_idx ON allwise (cntr);
