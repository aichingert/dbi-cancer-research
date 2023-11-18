use std::sync::Arc;
use oracle::Connection;
use serde::{Serialize, Deserialize};

type SQL = &'static str;

const GENE_FROM_IDENT_SQL: SQL = "SELECT gene_id, being_id, name, essential_score FROM gene WHERE name = :1";
const GENE_FROM_ID_SQL: SQL = "SELECT gene_id, being_id, name, essential_score FROM gene WHERE gene_id = :1";

const GENE_LETHALITY_SQL: SQL = "SELECT GENE1_ID, GENE2_ID, SCORE FROM SYN_LETH WHERE GENE1_ID = :1 OR GENE2_ID = :1";
const GENES_MAPPING_SQL: SQL = "SELECT GENE1_ID, GENE2_ID FROM MAPPING WHERE GENE1_ID = :1 OR GENE2_ID = :1";

#[derive(Clone)]
pub struct AppState {
    username: Arc<str>,
    password: Arc<str>,
    database: Arc<str>,
}

impl AppState {
    pub fn new<'a>(username: &'a str, password: &'a str, database: &'a str) -> Self {
        Self {
            username: username.into(),
            password: password.into(),
            database: database.into(),
        }
    }

    pub fn get_connection(&self) -> Result<Connection, oracle::Error> {
        Connection::connect::<&str, &str, &str>(self.username.as_ref(), self.password.as_ref(), self.database.as_ref())
    }
}

#[derive(Serialize, Deserialize)]
pub struct LethalGenes {
    pub request_gene: Gene,
    pub human_genes: Vec<Lethal>,
    pub mouse_genes: Vec<Lethal>,
    pub yeast_genes: Vec<Lethal>,
}

#[derive(Serialize, Deserialize)]
pub struct Lethal {
    gene: Gene,
    lethality_score: f32,
}

#[derive(Serialize, Deserialize)]
pub struct Gene {
    id: i64,
    being: i64,
    name: String,
    essentiality_score: Option<f32>,
}

impl Gene {
    pub fn new(ident: &str, sql: Option<(SQL, i64)>, conn: &Connection) -> Result<Self, oracle::Error> {
        let row = if let Some((sql, param)) = sql {
            conn.query_row(sql, &[&param])?
        } else {
            conn.query_row(GENE_FROM_IDENT_SQL, &[&ident])?
        };

        Ok(Self {
            id: row.get("gene_id").unwrap(),
            name: row.get("name").unwrap(),
            being: row.get("being_id").unwrap(),
            essentiality_score: row.get("essential_score").unwrap(),
        })
    }

    pub fn is_lethal_for(&self, connection: &Connection) -> Result<Vec<Lethal>, oracle::Error> {
        let results = connection.query(GENE_LETHALITY_SQL, &[&self.id])?;
        let mut lethal_genes = Vec::new();

        for row_res in results {
            let row = row_res?;
            let lethality: f32 = row.get("score")?;
 
            if lethality < 0.65 {
                continue;
            }

            let id = self.get_other_id(row.get("gene1_id")?, row.get("gene2_id")?);
            let gene = Gene::new("", Some((GENE_FROM_ID_SQL, id)), connection)?;

            lethal_genes.push(Lethal {
                gene,
                lethality_score: lethality,
            });
        }

        Ok(lethal_genes)
    }

    pub fn map_to_being(&self, being: i64, conn: &Connection) -> Result<Vec<Gene>, oracle::Error> {
        let results = conn.query(GENES_MAPPING_SQL, &[&self.id])?;
        let mut mapped_genes = Vec::new();

        for row_result in results {
            let row = row_result?;

            let gene = Gene::new("",
                Some((GENE_FROM_ID_SQL, self.get_other_id(row.get("gene1_id")?, row.get("gene2_id")?))),
                conn
            )?;

            if gene.being != being {
                continue;
            }

            mapped_genes.push(gene);
        }

        Ok(mapped_genes)
    }

    pub fn filter(genes: &Vec<Gene>, connection: &Connection) -> Vec<Lethal> {
        let mut lethal = Vec::new();

        for gene in genes {
            let lethal_genes = gene.is_lethal_for(connection).unwrap();

            for lethal_gene in lethal_genes {
                if lethal_gene.gene.map_to_being(1, connection).unwrap().len() > 0 {
                    lethal.push(lethal_gene);
                }
            }
        }

        lethal
    }

    fn get_other_id(&self, gene1_id: i64, gene2_id: i64) -> i64 {
        if self.id == gene1_id {
            gene2_id
        } else {
            gene1_id
        }
    }
}
