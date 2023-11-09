--BEING
INSERT INTO BEING (being_id, name) VALUES (1, 'Human');
INSERT INTO BEING (being_id, name) VALUES (2, 'Mouse');
INSERT INTO BEING (being_id, name) VALUES (3, 'Jade');

--GENE
INSERT INTO GENE (gene_id, being_id, name, essential_score) VALUES (1, 1, 'GeneX', 0.6);
INSERT INTO GENE (gene_id, being_id, name, essential_score) VALUES (2, 1, 'GeneD', 0.6);
INSERT INTO GENE (gene_id, being_id, name, essential_score) VALUES (3, 2, 'GeneS', 0.7);
INSERT INTO GENE (gene_id, being_id, name, essential_score) VALUES (4, 2, 'GeneAS', 0.5);
INSERT INTO GENE (gene_id, being_id, name, essential_score) VALUES (5, 3, 'GeneGA', 0.4);
INSERT INTO GENE (gene_id, being_id, name, essential_score) VALUES (6, 3, 'GeneH', 0.4);
INSERT INTO GENE (gene_id, being_id, name, essential_score) VALUES (7, 3, 'GeneJ', 0.2);


INSERT INTO MAPPING (gene1_id, gene2_id) VALUES (1,3);
INSERT INTO MAPPING (gene1_id, gene2_id) VALUES (1,6);
INSERT INTO MAPPING (gene1_id, gene2_id) VALUES (1,7);
INSERT INTO MAPPING (gene1_id, gene2_id) VALUES (2,5);
INSERT INTO MAPPING (gene1_id, gene2_id) VALUES (2,4);


SELECT * FROM GENE;

commit;