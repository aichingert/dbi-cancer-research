use oracle::{Connection, Result};
use axum::{
    extract::{Path, State}, 
    routing::get, 
    Router
};

mod model;
use model::{AppState, Gene};

// 1.
// GENE -> ID
//
// being_type
// gene_name
//
// 2. 
// Map -> Animals
// lethality and essentiality checking
// |
// v
// Resulting in genes that need to be mapped
// to human genes again 

async fn get_lethality() {

}

fn test_query(connection: &Connection) -> Result<()> {
    let sql = "select * from being";
    let mut stmt = connection.statement(sql).build()?;
    let rows = stmt.query(&[])?;

    for row_result in rows {
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
