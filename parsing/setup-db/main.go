package main

import (
	"database/sql"
	"fmt"
	_ "github.com/godror/godror"
	"os"
	"strconv"
	"strings"
)

const (
	HUMAN = 1
	MOUSE = 2
	YEAST = 3
)

type gene struct {
	id             int
	name           string
	essentialScore *float32
	being          int
}

type dependency struct {
	gene1 gene
	gene2 gene
	score float32
}

type mapping struct {
	gene1 gene
	gene2 gene
}

func main() {
	db, err := sql.Open("godror", "SYSTEM/lol@localhost")

	if err != nil {
		fmt.Println("conn str wrong")
	}

	err = db.Ping()

	if err != nil {
		fmt.Println("db not working ;(", err)
	}

	dropTables(db)
	createTables(db)

	fmt.Println("Reading Files")

	human, err := os.ReadFile("../data/syn-leth/Human_SL.csv")

	if err != nil {
		fmt.Println("Error Mouse", err)
	}

	mouse, err := os.ReadFile("../data/syn-leth/Mouse_SL.csv")

	if err != nil {
		fmt.Println("Error Mouse", err)
	}

	yeast, err := os.ReadFile("../data/syn-leth/Yeast_SL.csv")

	if err != nil {
		fmt.Println("Error Yeast", err)
	}

	scores, err := os.ReadFile("../data/OGEE/CSEGs_CEGs.txt")

	if err != nil {
		fmt.Println("Error essential scores", err)
	}

	mouseMappings, err := os.ReadFile("../data/oma/mouse-mapping.csv")

	if err != nil {
		fmt.Println("Error Mouse Mapping", err)
	}

	yeastMappings, err := os.ReadFile("../data/oma/yeast-mapping.csv")

	if err != nil {
		fmt.Println("Error Yeast Mapping", err)
	}

	fmt.Println("Processing genes and dependencies")
	genes, dependencies := fillGenesAndDependencies(
		formatCSV(human, true),
		formatCSV(mouse, true),
		formatCSV(yeast, true))

	fmt.Println("Adding essential score")
	addEssentialScore(genes, formatEssentialScores(scores))

	fmt.Println("Processing Mappings")
	mappings := fillMappings(genes, formatCSV(mouseMappings, false), formatCSV(yeastMappings, false))

	fmt.Println("Writing data to db (0/4)")

	fmt.Println("Writing beings to db (1/4)")
	writeBeingsToDb(db)

	fmt.Println("Writing genes to db (2/4)")
	writeGenesToDb(db, genes)

	fmt.Println("Writing dependencies to db (3/4)")
	writeDependenciesToDb(db, dependencies)

	fmt.Println("Writing mappings to db (4/4)")
	writeMappingsToDb(db, mappings)
}

func fillMappings(genes []gene, mouseMappings [][]string, yeastMappings [][]string) []mapping {
	mouse := fillBeingMapping(genes, mouseMappings, MOUSE)
	yeast := fillBeingMapping(genes, yeastMappings, YEAST)

	return append(mouse, yeast...)
}

func fillBeingMapping(genes []gene, beingMapping [][]string, beingId int) []mapping {
	var mappings []mapping

	for _, mappingRow := range beingMapping {
		id1, err1 := strconv.Atoi(mappingRow[0])
		id2, err2 := strconv.Atoi(mappingRow[1])

		if err1 == nil && err2 == nil {
			gene1 := containsOrGet(genes, id1)
			gene2 := containsOrGet(genes, id2)

			if gene1 != nil && gene2 != nil && gene1.being == HUMAN && gene2.being == beingId {
				mappings = append(mappings, mapping{
					gene1: *gene1,
					gene2: *gene2,
				})
			}
		}
	}

	return mappings
}

func formatEssentialScores(scores []byte) [][]string {
	test := string(scores)
	beingSplit := strings.Split(test, "\n")[12:]

	var beingFormatted [][]string

	for _, beingInformation := range beingSplit {
		beingFormatted = append(beingFormatted, strings.Split(beingInformation, "\t"))
	}

	return beingFormatted
}

func addEssentialScore(genes []gene, essentialScores [][]string) {
	for _, essentialScore := range essentialScores {
		geneId := getGeneIdByName(essentialScore[0], genes)

		if geneId != -1 {
			var score float32

			switch essentialScore[1] {
			case "CSEGs":
				score = 0.8
				break
			case "CEGs":
				score = 0.1
				break
			}

			genes[geneId].essentialScore = &score
		}
	}
}

func formatCSV(being []byte, isSynLeth bool) [][]string {
	beingSplit := strings.Split(string(being), "\n")[1:]

	var beingFormatted [][]string
	for _, beingInformation := range beingSplit {
		var tempInformation = strings.Split(beingInformation, ",")

		if len(tempInformation) != 8 && isSynLeth {
			for i := 0; i < len(tempInformation); i++ {
				if tempInformation[i][0] == '"' {
					startSlice := tempInformation[:i]
					endSlice := tempInformation[(i + len(tempInformation) - 7):]

					concatInformation := tempInformation[i]

					for j := 0; j < len(tempInformation)-8; j++ {
						concatInformation += "," + tempInformation[i+j+1]
					}

					tempInformation = append(append(startSlice, concatInformation), endSlice...)
					break
				}
			}
		}

		beingFormatted = append(beingFormatted, tempInformation)
	}

	return beingFormatted
}

func fillGenesAndDependencies(human [][]string, mouse [][]string, yeast [][]string) ([]gene, []dependency) {
	humanGenes, humanDependencies := fillBeing(human, HUMAN)
	mouseGenes, mouseDependencies := fillBeing(mouse, MOUSE)
	yeastGenes, yeastDependencies := fillBeing(yeast, YEAST)

	humanGenes = append(humanGenes, mouseGenes...)
	humanDependencies = append(humanDependencies, mouseDependencies...)

	return append(humanGenes, yeastGenes...), append(humanDependencies, yeastDependencies...)
}

func fillBeing(being [][]string, beingId int) ([]gene, []dependency) {
	var genes []gene
	var dependencies []dependency

	for _, result := range being {
		id1, err := strconv.Atoi(result[1])

		if err != nil {
			fmt.Println(result)
			fmt.Println("Error Id1", err, beingId)
		}

		gene1 := containsOrGet(genes, id1)

		if gene1 == nil {
			gene1 = &gene{
				id:             id1,
				name:           result[0],
				essentialScore: nil,
				being:          beingId,
			}

			genes = append(genes, *gene1)
		}

		id2, err := strconv.Atoi(result[3])

		if err != nil {
			fmt.Println(result)
			fmt.Println("Error Id2", err, beingId)
		}

		gene2 := containsOrGet(genes, id2)
		if gene2 == nil {
			gene2 = &gene{
				id:             id2,
				name:           result[2],
				essentialScore: nil,
				being:          beingId,
			}
			genes = append(genes, *gene2)
		}

		score, err := strconv.ParseFloat(result[7], 32)

		if err != nil {
			fmt.Println(result)
			fmt.Println("Error Score", err, beingId)
		}

		dependencies = append(dependencies, dependency{
			gene1: *gene1,
			gene2: *gene2,
			score: float32(score),
		})
	}

	return genes, dependencies
}

func containsOrGet(genes []gene, id int) *gene {
	for _, gene := range genes {
		if gene.id == id {
			return &gene
		}
	}

	return nil
}

func getGeneIdByName(name string, genes []gene) int {
	for i, gene := range genes {
		if gene.name == name {
			return i
		}
	}

	return -1
}

func writeBeingsToDb(db *sql.DB) {
	beingNames := []string{"Human", "Mouse", "Yeast"}
	insert, err := db.Prepare(`INSERT INTO BEING (being_id,name) VALUES (:1,:2)`)

	if err != nil {
		fmt.Println("ERROR preparing insert Being", err)
		return
	}

	defer func(insert *sql.Stmt) {
		err = insert.Close()
		if err != nil {
			fmt.Println("huh")
		}
	}(insert)

	for i := 0; i < len(beingNames); i++ {
		_, err = insert.Exec(i+1, beingNames[i])

		if err != nil {
			fmt.Println("ERROR inserting into Being", err)
		}
	}
}

func writeGenesToDb(db *sql.DB, genes []gene) {
	insert, err := db.Prepare(`INSERT INTO GENE (gene_id, name, essential_score, being_id) VALUES (:1,:2,:3,:4)`)

	if err != nil {
		fmt.Println("ERROR preparing insert Gene", err)
	}

	defer func(insert *sql.Stmt) {
		err = insert.Close()
		if err != nil {
			fmt.Println("huh")
		}
	}(insert)

	for _, gene := range genes {
		if gene.essentialScore == nil {
			_, err = insert.Exec(gene.id, gene.name, nil, gene.being)
		} else {
			_, err = insert.Exec(gene.id, gene.name, *gene.essentialScore, gene.being)
		}

		if err != nil {
			fmt.Println("ERROR inserting into Gene", err)
		}
	}
}

func writeDependenciesToDb(db *sql.DB, dependencies []dependency) {
	insert, err := db.Prepare(`INSERT INTO SYN_LETH (gene1_id, gene2_id, score) VALUES (:1,:2,:3)`)

	if err != nil {
		fmt.Println("ERROR preparing insert SYN_LETH", err)
	}

	defer func(insert *sql.Stmt) {
		err = insert.Close()
		if err != nil {
			fmt.Println("huh")
		}
	}(insert)

	for _, dependency := range dependencies {
		_, err = insert.Exec(dependency.gene1.id, dependency.gene2.id, dependency.score)

		if err != nil {
			fmt.Println("ERROR inserting into SYN_LETH", err)
		}
	}
}

func writeMappingsToDb(db *sql.DB, mappings []mapping) {
	insert, err := db.Prepare(`INSERT INTO MAPPING (gene1_id, gene2_id) VALUES (:1,:2)`)

	if err != nil {
		fmt.Println("ERROR preparing insert Mapping", err)
	}

	defer func(insert *sql.Stmt) {
		err = insert.Close()
		if err != nil {
			fmt.Println("huh")
		}
	}(insert)

	for _, mapping := range mappings {
		_, _ = insert.Exec(mapping.gene1.id, mapping.gene2.id)
	}
}

func createTables(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE BEING(
    being_id NUMBER(1) PRIMARY KEY NOT NULL,
    name VARCHAR2(5) NOT NULL)`)

	if err != nil {
		fmt.Println("Being", err)
	}

	_, err = db.Exec(`CREATE TABLE GENE(
		gene_id NUMBER(10) PRIMARY KEY NOT NULL,
		being_id NUMBER(1) NOT NULL,
		name VARCHAR2(15) NOT NULL,
		essential_score FLOAT(5) NULL, 
		CONSTRAINT BEING_GENE_FK FOREIGN KEY (being_id) REFERENCES BEING(being_id))`)

	if err != nil {
		fmt.Println("GENE", err)
	}

	_, err = db.Exec(`CREATE TABLE SYN_LETH(
		gene1_id NUMBER(10) NOT NULL,
		gene2_id NUMBER(10) NOT NULL,
		score FLOAT(5) NOT NULL,
		PRIMARY KEY (gene1_id, gene2_id),
		CONSTRAINT GENE1_DEP_FK FOREIGN KEY (gene1_id) REFERENCES GENE(gene_id),
		CONSTRAINT GENE2_DEP_FK FOREIGN KEY (gene2_id) REFERENCES GENE(gene_id))`)

	if err != nil {
		fmt.Println("SYN_LETH", err)
	}

	_, err = db.Exec(`CREATE TABLE Mapping(
		gene1_id NUMBER(10) NOT NULL,
		gene2_id NUMBER(10) NOT NULL,
		PRIMARY KEY (gene1_id, gene2_id),
		CONSTRAINT GENE1_MAP_FK FOREIGN KEY (gene1_id) REFERENCES GENE(gene_id),
		CONSTRAINT GENE2_MAP_FK FOREIGN KEY (gene2_id) REFERENCES GENE(gene_id))`)

	if err != nil {
		fmt.Println("Mapping", err)
	}
}

func dropTables(db *sql.DB) {
	_, _ = db.Exec(`DROP TABLE MAPPING`)
	_, _ = db.Exec(`DROP TABLE SYN_LETH`)
	_, _ = db.Exec(`DROP TABLE GENE`)
	_, _ = db.Exec(`DROP TABLE BEING`)
}
