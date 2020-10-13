#[derive(Queryable)]
pub struct User {
    pub id: Uuid,
    pub osu_id: String,
    pub discord_id: String,
    pub created_at: NaiveDateTime,
    pub updated_at: Option<NaiveDateTime>,
}

#[derive(Insertable)]
#[table_name="users"]
pub struct NewUser<'a> {
    pub discord_id: &'a str,
    pub osu_id: &'a str,
}

pub fn create_user<'a>(conn: &PgConnection, osu_id: &'a str, discord_id: &'a str) -> User {
    use schema::posts;

    let new_user = NewUser {
        discord_id: discord_id,
        osu_id: osu_id,
    };

    diesel::insert_into(users::table)
        .values(&new_user)
        .get_result(conn)
        .expect("Error saving new user")
}
