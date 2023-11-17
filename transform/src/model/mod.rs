use serde::Deserialize;

pub mod app_state;
pub use app_state::AppState;

#[derive(Debug, Deserialize)]
pub struct Gene {
    pub name: String,
}
