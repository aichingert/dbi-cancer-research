use serde::{Serialize, Deserialize};
use sqlx::postgres::PgConnection;

#[derive(Serialize, Deserialize)]
pub struct LethalGenes {
    pub request_gene: Gene,
    pub human_genes: Vec<Lethal>,
    pub mouse_genes: Vec<Lethal>,
    pub yeast_genes: Vec<Lethal>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct Lethal {
    pub gene: Gene,
    pub lethality_score: f32,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct Gene {
    pub gene_id: i32,
    being_id: i16,
    name: String,
    essential_score: Option<f32>,
}

pub struct Mapping {
    gene1_id: i32,
    gene2_id: i32,
}

pub struct SynLeth {
    gene1_id: i32,
    gene2_id: i32,
    score: f32,
}

impl Gene {
    pub async fn new(ident: &str, conn: &mut PgConnection) -> Result<Gene, sqlx::Error> {
        sqlx::query_as!(Gene, "SELECT * FROM gene WHERE name = UPPER($1)", ident)
            .fetch_one(conn)
            .await
    }

    pub async fn from_id(id: i32, conn: &mut PgConnection) -> Result<Gene, sqlx::Error> {
        sqlx::query_as!(Gene, "SELECT * FROM gene WHERE gene_id = $1", id)
            .fetch_one(conn)
            .await
    }

    pub async fn lethal_pairs(id: i32, conn: &mut PgConnection) -> Result<Vec<(i32, f32)>, sqlx::Error> {
        let pairs = sqlx::query_as!(SynLeth, "SELECT * FROM Syn_Leth WHERE gene1_id = $1 OR gene2_id = $1", id)
            .fetch_all(conn)
            .await?;

        Ok(pairs.iter().map(|pair| (if pair.gene1_id == id
            { pair.gene2_id } else
            { pair.gene1_id }, pair.score)
        ).collect::<Vec<_>>())
    }

    pub async fn get_mappings(id: i32, conn: &mut PgConnection) -> Result<Vec<i32>, sqlx::Error> {
        let mappings = sqlx::query_as!(Mapping, "SELECT * FROM mapping WHERE gene1_id = $1 OR gene2_id = $1", id)
            .fetch_all(conn)
            .await?;

        Ok(mappings.iter().map(|mapping| if mapping.gene1_id == id
            { mapping.gene2_id } else 
            { mapping.gene1_id }
            ).collect::<Vec<_>>()
        )
    }
}

impl LethalGenes {
    pub async fn new(gene: Gene, conn: &mut PgConnection) -> Result<LethalGenes, sqlx::Error> {
        let mut lethal_genes = Self {
            request_gene: gene,
            human_genes: Vec::new(),
            mouse_genes: Vec::new(),
            yeast_genes: Vec::new(),
        };

        for (id, score) in Gene::lethal_pairs(lethal_genes.request_gene.gene_id, &mut *conn).await? {
            lethal_genes.human_genes.push(Lethal {
                gene: Gene::from_id(id, &mut *conn).await?,
                lethality_score: score,
            });
        }

        for mapping in Gene::get_mappings(lethal_genes.request_gene.gene_id, &mut *conn).await? {
            let mapped_gene = Gene::from_id(mapping, conn).await?;
            let mut is_lethal: bool = false;

            for pair in Gene::lethal_pairs(mapped_gene.gene_id, conn).await? {
                let mappings = Gene::get_mappings(pair.0, conn).await?;

                for mapping in mappings {
                    if Gene::from_id(mapping, conn).await?.being_id == 1 {
                        is_lethal = true;
                        break;
                    }
                }

                let lethal = Lethal { 
                    gene: Gene::from_id(pair.0, conn).await?, 
                    lethality_score: pair.1 
                };

                match (is_lethal, lethal.gene.being_id) {
                    (true, 2) => lethal_genes.mouse_genes.push(lethal),
                    (true, 3) => lethal_genes.mouse_genes.push(lethal),
                    _ => (),
                }
            } 
        }

        Ok(lethal_genes)
    }
}
