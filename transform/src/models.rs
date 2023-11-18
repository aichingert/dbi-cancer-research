use std::sync::Arc;
use oracle::Connection;
use serde::Deserialize;

const GENE_SQL: &'static str = "SELECT gene_id, being_id, name, essential_score FROM gene WHERE name = :1";
const GENE_LETHALITY: &'static str = "SELECT GENE1_ID, GENE2_ID, SCORE FROM SYN_LETH WHERE GENE1_ID = :1 OR GENE2_ID = :1";

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

#[derive(Debug, Deserialize)]
pub struct Gene {
    pub id: i64,
    pub name: String,
    pub being: String,
    pub score: Option<f32>,
}

impl Gene {
    pub fn new(ident: &str, connection: &Connection) -> Result<Self, oracle::Error> {
        let row = connection.query_row(GENE_SQL, &[&ident])?;

        Ok(Self {
            id: row.get("gene_id").unwrap(),
            name: row.get("name").unwrap(),
            being: row.get("being_id").unwrap(),
            score: row.get("essential_score").unwrap(),
        })
    }

    pub fn is_lethal(&self, connection: &Connection) -> Result<String, oracle::Error> {
        let results = connection.query(GENE_LETHALITY, &[&self.id])?;



        Ok(String::new())
    }
}
