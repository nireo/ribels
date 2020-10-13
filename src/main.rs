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
mod models;

struct Handler;

static PREFIX: &str = "%";

#[async_trait]
impl EventHandler for Handler {
    async fn message(&self, ctx: Context, msg: Message) {
        // tokenize the input:
        let args = msg.content.split(" ");
        let args = args.collect::<Vec<&str>>();

        if args[0] == format!("{}set", PREFIX) {
            let formatted_name = args[1..args.len()].join("-");
            let conn = db::establish_connection();  
            let user = models::create_user(&conn, &formatted_name, &msg.author.id.to_string());
            
            if let Err(why) = msg.channel_id.say(&ctx.http, "Added to the database!").await {
                println!("Error adding  to database")
            }
        }

        /*
        if args[0] == format!("{}osu", PREFIX) {
            let formatted_name = args[1..args.len()].join("_");

            // load from the api
            let key: &str = "";
            let res = reqwest::get(format!("https://osu.ppy.sh/api/get_user_best?u={}&k={}", formatted_name, key)).await;

            /*
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
                    }
                },
                Err(err) => {
                    if let Err(why) = msg.channel_id.say(&ctx.http, "Problem getting user info!").await {
                        println!("Error sending message: {:?}", why);
                    }
                },
            } */
        } */
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

    let mut client = Client::new(&token)
        .event_handler(Handler)
        .await
        .expect("Err creating client");

    if let Err(why) = client.start().await {
        println!("Client error: {:?}", why);
    }
}

