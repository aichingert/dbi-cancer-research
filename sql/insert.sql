--BEING
INSERT INTO BEING (being_id, name) VALUES (1, 'Human');
INSERT INTO BEING (being_id, name) VALUES (2, 'Mouse');
INSERT INTO BEING (being_id, name) VALUES (3, 'Jade');

--GENE
INSERT INTO GENE (gene_id, being_id, name, essential_score) VALUES (2, 1, 'Gene2_Human', 0.6);
INSERT INTO GENE (gene_id, being_id, name, essential_score) VALUES (3, 2, 'Gene1_Mouse', 0.7);
INSERT INTO GENE (gene_id, being_id, name, essential_score) VALUES (4, 2, 'Gene2_Mouse', 0.5);
INSERT INTO GENE (gene_id, being_id, name, essential_score) VALUES (5, 3, 'Gene1_Jada', 0.4);

SELECT * FROM GENE;

commit;