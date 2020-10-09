use diesel::pg::PgConnection;
use diesel::r2d2::ConnectionManager;
use lazy_static::lazy_static;
use std::env;
use diesel::r2d2;

type Pool = r2d2::Pool<ConnectionManager<PgConnection>>;
pub type DbConnection = r2d2::PooledConnection<ConnectionManager<PgConnection>>;

embed_migrations!();

lazy_static! {
    static ref POOL: Pool = {
        let db_url = env::var("DATABASE_URL").expect("database url not added");
        let manager = ConnectionManager::<PgConnection>::new(db_url);
        Pool::new(manager).expect("Failed to create db pool")
    }
}

pub fn init() {
    info!("Initializing the database");
    lazy_static::initialize(&POOL);
    let conn = connection().expect("Failed to get database connection");
    embed_migrations::run(&conn).unwrap();
}


pub fn connection() -> Result<DbConnection, &'static str> {
    POOL.get().map_err(|e| Err("could not get database connection"))
}

