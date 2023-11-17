use std::sync::Arc;
use oracle::Connection;
use serde::Deserialize;

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
    pub name: String,
}
