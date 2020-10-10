use std::env;

use serenity::{
    async_trait,
    model::{channel::Message, gateway::Ready},
    prelude::*,
};

use dotenv::dotenv;

extern crate reqwest;
extern crate diesel;
extern crate dotenv;

mod db;

struct Handler;

static PREFIX: &str = "%";

#[derive(Deserialize, Debug)]
struct OsuResponseUser {
    user_id: String,
    username: String,
    join_date: String,
    count300: String,
    count100: String,
    count50: string,
    playcount: String,
    ranked_score: String,
    total_score: String,
    pp_rank: String,
    level: String,
    pp_raw: String,
    accuracy: String,
    count_rank_ss: String,
    count_rank_ssh: String,
    count_rank_s: String,
    count_rank_sh: String,
    count_rank_a: String,
    country: String,
    total_seconds_played: String,
    pp_country_rank: String,
}

#[async_trait]
impl EventHandler for Handler {
    async fn message(&self, ctx: Context, msg: Message) {
        // tokenize the input:
        let args = msg.content.split(" ");
        let args = args.collect::<Vec<&str>>();

        if args[0] == format!("{}osu", PREFIX) {
            let formatted_name = args[1..args.len()].join("_");

            // load from the api
            let key: &str = "";
            let res = reqwest::get(format!("https://osu.ppy.sh/api/get_user_best?u={}&k={}", formatted_name, key)).await;

            match res {
                Ok(v) => {
                    let body: Result<OsuResponseUser, E> = v.json().await;
                    match body {
                        Ok(v) => {
                            if let Err(why) = msg.channel_id.say(&ctx.http, "Found user!").await {
                                println!("Error sending message: {:?}", why);
                            }
                        }
                        Err(err) => {
                            if let Err(why) = msg.channel_id.say(&ctx.http, "Problem getting user info!").await {
                                println!("Error sending message: {:?}", why);
                            }
                        }
                    },
                },
                Err(err) => {
                    if let Err(why) = msg.channel_id.say(&ctx.http, "Problem getting user info!").await {
                        println!("Error sending message: {:?}", why);
                    }
                },
            }

        }
    }

    async fn ready(&self, _: Context, ready: Ready) {
        println!("{} is connected!", ready.user.name);
    }
}

#[tokio::main]
async fn main() {
    dotenv().ok();
    let token = env::var("DISCORD_TOKEN")
        .expect("Expected a token in the environment");

    let connection = db::establish_connection();

    let mut client = Client::new(&token)
        .event_handler(Handler)
        .await
        .expect("Err creating client");

    if let Err(why) = client.start().await {
        println!("Client error: {:?}", why);
    }
}

