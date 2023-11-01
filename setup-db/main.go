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
	WORM  = 3
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

func main() {
	//connStr := fmt.Sprintf("oracle://SYSTEM:lol@localhost:%d/ORACLE?PREFETCH_ROWS=%d", 1521, 500)

	/*db, _ := sql.Open("oracle", connStr)
	err := db.Ping()
	createTables(db)*/

	writeBeings()

	human, _ := os.ReadFile("../data/syn-leth/Human_SL.csv")
	humanSplit := strings.Split(string(human), "\n")

	var humanFormatted [][]string
	for _, humanInformation := range humanSplit {
		humanFormatted = append(humanFormatted, strings.Split(humanInformation, ","))
	}

	mouse, _ := os.ReadFile("../data/syn-leth/Mouse_SL.csv")
	mouseSplit := strings.Split(string(mouse), "\n")

	var mouseFormatted [][]string
	for _, mouseInformation := range mouseSplit {
		mouseFormatted = append(mouseFormatted, strings.Split(mouseInformation, ","))
	}

	worm, _ := os.ReadFile("../data/syn-leth/Worm_SL.csv")
	wormSplit := strings.Split(string(worm), "\n")

	var wormFormatted [][]string
	for _, wormInformation := range wormSplit {
		wormFormatted = append(wormFormatted, strings.Split(wormInformation, ","))
	}

	genes, dependencies := fillGenesAndDependencies(humanFormatted, mouseFormatted, wormFormatted)
	writeGenes(genes)
	writeDependencies(dependencies)
}

func fillGenesAndDependencies(human [][]string, mouse [][]string, worm [][]string) ([]gene, []dependency) {
	humanGenes, humanDependencies := fillBeing(human, HUMAN)
	mouseGenes, mouseDependencies := fillBeing(mouse, MOUSE)
	wormGenes, wormDependencies := fillBeing(worm, WORM)

	humanGenes = append(humanGenes, mouseGenes...)
	humanDependencies = append(humanDependencies, mouseDependencies...)

	return append(humanGenes, wormGenes...), append(humanDependencies, wormDependencies...)
}

func fillBeing(being [][]string, beingId int) ([]gene, []dependency) {
	var genes []gene
	var dependencies []dependency

	for _, result := range being {
		id1, _ := strconv.Atoi(result[1])
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

		id2, _ := strconv.Atoi(result[3])
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

		score, _ := strconv.ParseFloat(result[7], 32)

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
	beingNames := []string{"Human", "Mouse", "Worm"}
	var beings = "BeingId;Name\n"

	for i := 0; i < len(beingNames); i++ {
		beings += fmt.Sprintf("%d;%s\n", i+1, beingNames[i])
	}

	_ = os.WriteFile("../output/Beings.csv", []byte(beings), 0644)
}

func writeGenes(genes []gene) {
	var beings = "GeneId;Name;EssentialScore;Being\n"

	for _, gene := range genes {
		if gene.essentialScore == nil {
			beings += fmt.Sprintf("%d;%s;;%d\n", gene.id, gene.name, gene.being)
		} else {
			beings += fmt.Sprintf("%d;%s;%f;%d\n", gene.id, gene.name, *gene.essentialScore, gene.being)
		}
	}

	_ = os.WriteFile("../output/Genes.csv", []byte(beings), 0644)
}

func writeDependencies(dependencies []dependency) {
	var beings = "Gene1Id;Gene2Id;Score\n"

	for _, dependency := range dependencies {
		beings += fmt.Sprintf("%d;%d;%f\n", dependency.gene1.id, dependency.gene2.id, dependency.score)
	}

	_ = os.WriteFile("../output/Dependencies.csv", []byte(beings), 0644)
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
