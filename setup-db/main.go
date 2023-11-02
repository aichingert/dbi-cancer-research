package main

import (
	"fmt"
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
	//connStr := fmt.Sprintf("oracle://SYSTEM:lol@localhost:1521/ORACLE?PREFETCH_ROWS=%d", 500)

	/*db, _ := sql.Open("oracle", connStr)
	err := db.Ping()
	createTables(db)*/

	writeBeings()

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

	genes, dependencies := fillGenesAndDependencies(
		formatCSV(human, true),
		formatCSV(mouse, true),
		formatCSV(yeast, true))

	mouseMappings, err := os.ReadFile("../data/oma/mouse-mapping.csv")

	if err != nil {
		fmt.Println("Error Mouse Mapping", err)
	}

	yeastMappings, err := os.ReadFile("../data/oma/yeast-mapping.csv")

	if err != nil {
		fmt.Println("Error Yeast Mapping", err)
	}

	mappings := fillMappings(genes, formatCSV(mouseMappings, false), formatCSV(yeastMappings, false))

	writeGenes(genes)
	writeDependencies(dependencies)
	writeMappings(mappings)
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

func writeBeings() {
	beingNames := []string{"Human", "Mouse", "Yeast"}
	var beings = "BeingId,Name\n"

	for i := 0; i < len(beingNames); i++ {
		beings += fmt.Sprintf("%d,%s\n", i+1, beingNames[i])
	}

	_ = os.WriteFile("../output/Beings.csv", []byte(beings), 0644)
}

func writeGenes(genes []gene) {
	var fileGenes = "GeneId,Name,EssentialScore,Being\n"

	for _, gene := range genes {
		if gene.essentialScore == nil {
			fileGenes += fmt.Sprintf("%d,%s,,%d\n", gene.id, gene.name, gene.being)
		} else {
			fileGenes += fmt.Sprintf("%d,%s,%f,%d\n", gene.id, gene.name, *gene.essentialScore, gene.being)
		}
	}

	_ = os.WriteFile("../output/Genes.csv", []byte(fileGenes), 0644)
}

func writeDependencies(dependencies []dependency) {
	var fileDependencies = "Gene1Id,Gene2Id,Score\n"

	for _, dependency := range dependencies {
		fileDependencies += fmt.Sprintf("%d,%d,%f\n", dependency.gene1.id, dependency.gene2.id, dependency.score)
	}

	_ = os.WriteFile("../output/Dependencies.csv", []byte(fileDependencies), 0644)
}

func writeMappings(mappings []mapping) {
	var fileMappings = "Gene1Id,Gene2Id\n"

	for _, mapping := range mappings {
		fileMappings += fmt.Sprintf("%d,%d\n", mapping.gene1.id, mapping.gene2.id)
	}

	_ = os.WriteFile("../output/Mappings.csv", []byte(fileMappings), 0644)
}

/*func createTables(db *sql.DB) {
	_, _ = db.Exec("CREATE IF NOT EXISTS TABLE BEING(" +
		"being_id NUMBER(1) PRIMARY KEY NOT NULL," +
		"name VARCHAR2(5) NOT NULL);")

	_, _ = db.Exec("CREATE IF NOT EXISTS TABLE GENE(" +
		"gene_id NUMBER(10) PRIMARY KEY NOT NULL," +
		"being_id NUMBER(1) NOT NULL," +
		"name VARCHAR2(15) NOT NULL," +
		"essential_score FLOAT(5) NULL," +
		"CONSTRAINT BEING_GENE_FK FOREIGN KEY (being_id) REFERENCES BEING(being_id));")

	_, _ = db.Exec("CREATE IF NOT EXISTS TABLE SYN_LETH(" +
		"gene1_id NUMBER(10) NOT NULL," +
		"gene2_id NUMBER(10) NOT NULL," +
		"score FLOAT(5) NOT NULL," +
		"PRIMARY KEY (gene1_id, gene2_id)," +
		"CONSTRAINT GENE1_DEP_FK FOREIGN KEY (gene1_id) REFERENCES GENE(gene_id)," +
		"CONSTRAINT GENE2_DEP_FK FOREIGN KEY (gene2_id) REFERENCES GENE(gene_id));")

	_, _ = db.Exec("CREATE IF NOT EXISTS TABLE Mapping(" +
		"gene1_id NUMBER(10) NOT NULL," +
		"gene2_id NUMBER(10) NOT NULL," +
		"PRIMARY KEY (gene1_id, gene2_id)," +
		"CONSTRAINT GENE1_MAP_FK FOREIGN KEY (gene1_id) REFERENCES GENE(gene_id)," +
		"CONSTRAINT GENE2_MAP_FK FOREIGN KEY (gene2_id) REFERENCES GENE(gene_id));")
}*/
