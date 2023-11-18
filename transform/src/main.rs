use axum::{
    extract::{Path, State}, 
    routing::get, 
    Router,
    Json
};

mod models;
use models::{AppState, Gene};

const LETHALITY_SQL: &'static str = "SELECT * FROM SYN_LETH WHERE GENE1_ID = :1 OR GENE2_ID = :1";
const ESSENTIALITY_SQL: &'static str ="SELECT * FROM Gene WHERE name = :1 AND essentialScore <> NULL";

async fn get_lethality(
    Path(ident): Path<String>,
    State(state): State<AppState>,
) -> Json<Vec<String>> {
    let connection = match state.get_connection() {
        Ok(conn) => conn,
        Err(_) => {
            dbg!("ERROR: failed to connect to DB");
            return Json(vec!["Internal Server Error".to_string()]);
        }
    };

    // Implement into Response for result and propagate errors with ?
    let gene = Gene::new(&ident.to_uppercase(), &connection);

    if gene.is_err() {
        return Json(vec!["no value".to_string()]);
    };

    println!("{:?}", gene);

    Json(Vec::new())
}

#[tokio::main]
async fn main() {
    let state = AppState::new("system", "lol", "localhost:1521");

    let app: Router = Router::new()
        .route("/:gene", get(get_lethality))
        .with_state(state);

    dbg!("Listening on localhost:3000");
    axum::Server::bind(&"127.0.0.1:3000".parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}
