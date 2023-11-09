DROP TABLE MAPPING;
DROP TABLE SYN_LETH;
DROP TABLE GENE;
DROP TABLE BEING;

CREATE TABLE BEING(
		being_id NUMBER(1) PRIMARY KEY NOT NULL,
		name VARCHAR2(5) NOT NULL
);
CREATE TABLE GENE(
        gene_id NUMBER(10) PRIMARY KEY NOT NULL,
		being_id NUMBER(1) NOT NULL,
        name VARCHAR2(15) NOT NULL,
		essential_score FLOAT(5) NULL,
		CONSTRAINT BEING_GENE_FK FOREIGN KEY (being_id) REFERENCES BEING(being_id)
);

CREATE TABLE SYN_LETH(
		gene1_id NUMBER(10) NOT NULL,
		gene2_id NUMBER(10) NOT NULL,
		score FLOAT(5) NOT NULL,
		PRIMARY KEY (gene1_id, gene2_id),
        CONSTRAINT GENE1_DEP_FK FOREIGN KEY (gene1_id) REFERENCES GENE(gene_id),
        CONSTRAINT GENE2_DEP_FK FOREIGN KEY (gene2_id) REFERENCES GENE(gene_id)
);

Create TABLE MAPPING(
        gene1_id NUMBER(10) NOT NULL,
        gene2_id NUMBER(10) NOT NULL,
        PRIMARY KEY (gene1_id, gene2_id),
        CONSTRAINT GENE1_MAP_FK FOREIGN KEY (gene1_id) REFERENCES GENE(gene_id),
        CONSTRAINT GENE2_MAP_FK FOREIGN KEY (gene2_id) REFERENCES GENE(gene_id)
);
commit;