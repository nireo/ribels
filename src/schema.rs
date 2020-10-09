table! {
    user (id) {
        id -> Uuid,
        osu_id -> Text,
        discord_id -> Text,
        created_at -> Timestamp,
        updated_at -> Nullable<Timestamp>,
    }
}
