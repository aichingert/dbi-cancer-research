CREATE OR REPLACE PACKAGE CANCER_Research AS
    PROCEDURE get_connected_genes(geneID NUMBER);
END CANCER_Research;
/

CREATE OR REPLACE PACKAGE BODY CANCER_Research AS

    PROCEDURE get_connected_genes(geneID NUMBER) IS
    BEGIN
        DBMS_OUTPUT.PUT_LINE('-------------------------------------------------');
        FOR gen IN (SELECT gene2_id
                         FROM MAPPING
                         WHERE gene1_id = geneID)
            LOOP
                DBMS_OUTPUT.PUT_LINE('Connected Gene ID: ' || gen.GENE2_ID);
            END LOOP;
        DBMS_OUTPUT.PUT_LINE('-------------------------------------------------');

    END get_connected_genes;
END CANCER_Research;
/
BEGIN
CANCER_Research.get_connected_genes(1);
END;
/