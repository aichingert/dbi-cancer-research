use axum::{
    async_trait,
    extract::{Path, FromRef, FromRequestParts, State}, 
    routing::get, 
    response::IntoResponse,
    http::{request::Parts, StatusCode},
    Router,
    Json,
};
use sqlx::postgres::{PgPool, PgPoolOptions};

mod models;
use models::{Gene, Lethal, LethalGenes};

struct DatabaseConnection(sqlx::pool::PoolConnection<sqlx::Postgres>);

#[async_trait]
impl<S> FromRequestParts<S> for DatabaseConnection
where
    PgPool: FromRef<S>,
    S: Send + Sync
{
    type Rejection = (StatusCode, String);

    async fn from_request_parts(_parts: &mut Parts, state: &S) -> Result<Self, Self::Rejection> {
        let pool = PgPool::from_ref(state);
        let conn = pool.acquire().await.map_err(internal_error)?;
        Ok(Self(conn))
    }
}

fn internal_error<E>(err: E) -> (StatusCode, String) 
where
    E: std::error::Error
{
    (StatusCode::INTERNAL_SERVER_ERROR, err.to_string())
}

async fn get_lethality(
    Path(ident): Path<String>,
    DatabaseConnection(mut conn): DatabaseConnection,
) -> impl IntoResponse {
    let gene = match Gene::new(&ident, &mut *conn).await {
        Ok(g) => g,
        Err(_)   => return (StatusCode::BAD_REQUEST, Json(None)),
    };
    
    (StatusCode::OK, Json(Some(LethalGenes::new(gene, &mut *conn).await.unwrap())))
}

#[tokio::main]
async fn main() {
    let connection_str = "postgres://postgres:lol@localhost";

    let pool = PgPoolOptions::new()
        .max_connections(5)
        .acquire_timeout(std::time::Duration::from_secs(3))
        .connect(connection_str).await.expect("can't connect to database");

    let app: Router = Router::new()
        .route("/:gene", get(get_lethality))
        .with_state(pool);

    dbg!("Listening on localhost:3000");
    axum::Server::bind(&"127.0.0.1:3000".parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}
