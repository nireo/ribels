pub struct User {
    pub id: Uuid,
    pub osu_id: String,
    pub discord_id: String,
    pub created_at: NaiveDateTime,
    pub updated_at: Option<NaiveDateTime>,
}

impl User {
}