use oracle::{Connection, Result};
use axum::{
    extract::{Path, State}, 
    routing::get, 
    Router
};

mod models;
use models::{AppState, Gene};

const LETHALITY_SQL: &'static str = "SELECT * FROM SynLeth WHERE geneId1 = :1 OR geneId2 = :1";
const ESSENTIALITY_SQL: &'static str ="SELECT * FROM Gene WHERE name = :1 AND essentialScore <> NULL";

async fn get_lethality(
    Path(gene): Path<String>,
    State(state): State<AppState>,
) -> Json<Vec<String>> {
    let connection = match state.get_connection() {
        Ok(conn) => conn,
        Err(e) => {
            dbg!("ERROR: failed to connect to DB");
            return Json(vec![]);
        }
    };

    // TODO: check errors
    let lethal_rows   = connection.query(LETHALITY_SQL, &[&gene]).unwrap();
    let essential_row = connection.query(ESSENTIALITY_SQL, &[&gene]).unwrap();

    println!("{:?}", lethal_rows);

    for lethal_row in lethal_rows {
        let row = lethal_row.unwrap();

        println!("{:?}", row);
    }

    Json(Vec::new())
}

fn test_query(connection: &Connection) -> Result<()> {
    let sql = "select * from being";
    let mut stmt = connection.statement(sql).build()?;
    let rows = stmt.query(&[])?;

rows    for row_result in rows {
        // print column values
        for (idx, val) in row_result?.sql_values().iter().enumerate() {
            if idx != 0 {
                print!(",");
            }
            print!("{}", val);
        }
        println!();
    }

    Ok(())
}

#[tokio::main]
async fn main() {
    let state = AppState::new("system", "lol", "localhost:1521");
    let connection = state.get_connection();

    let app: Router = Router::new()
        .route("/:gene", get(get_lethality))
        .with_state(state);

    test_query(&connection.unwrap()).unwrap();

    dbg!("Listening on localhost:3000");
    axum::Server::bind(&"127.0.0.1:3000".parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}
