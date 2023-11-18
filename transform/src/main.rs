use axum::{
    extract::{Path, State}, 
    routing::get, 
    response::IntoResponse,
    http::StatusCode,
    Router,
    Json,
};

mod models;
use models::{AppState, Gene};

async fn get_lethality(
    Path(ident): Path<String>,
    State(state): State<AppState>,
) -> impl IntoResponse {
    let connection = match state.get_connection() {
        Ok(conn) => conn,
        Err(_) => {
            dbg!("ERROR: failed to connect to DB");
            return (StatusCode::INTERNAL_SERVER_ERROR, Json(None));
        }
    };

    let gene = match Gene::new(&ident.to_uppercase(), &connection) {
        Ok(gene) => gene,
        Err(_) => return (StatusCode::BAD_REQUEST, Json(None))
    };

    let lethal_for_human_tests = gene.is_lethal_for(&connection).unwrap();

    let rats = Gene::filter(&gene.map_to_being(2, &connection).unwrap(), &connection);
    let yeast = Gene::filter(&gene.map_to_being(3, &connection).unwrap(), &connection);

    println!("{:?}", rats);
    println!("{:?}", yeast);

    (StatusCode::OK, Json(Some(lethal_for_human_tests)))
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
