use std::env;

use serenity::{
    async_trait,
    model::{channel::Message, gateway::Ready},
    prelude::*,
};

use dotenv::dotenv;

mod db;

struct Handler;

static PREFIX: char = '%';

#[async_trait]
impl EventHandler for Handler {
    async fn message(&self, ctx: Context, msg: Message) {
        // tokenize the input:
        let args = msg.content.split(" ");
        let args = args.collect::<Vec<&str>>();

        match args[0] {
            format!("{}recent", PREFIX) => {
                if let Err(why) = msg.channel_id.say(&ctx.http, "Your recent plays!").await {
                    println!("Error sending message: {:?}", why);
                }
            }

            format!("{}osutop", PREFIX) => {
                if let Err(why) = msg.channel_id.say(&ctx.http, "Your recent plays!").await {
                    println!("Error sending message: {:?}", why);
                }
            }

            format!("{}set", PREFIX) => {
                if let Err(why) = msg.channel_id.say(&ctx.http, "Your recent plays!").await {
                    println!("Error sending message: {:?}", why);
                }
            }

            format!("{}osu", PREFIX) => {
                if let Err(why) = msg.channel_id.say(&ctx.http, "Your recent plays!").await {
                    println!("Error sending message: {:?}", why);
                }
            }

            format!("{}recent", PREFIX) => {
                if let Err(why) = msg.channel_id.say(&ctx.http, "Your recent plays!").await {
                    println!("Error sending message: {:?}", why);
                }
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

    db::init();

    let mut client = Client::new(&token)
        .event_handler(Handler)
        .await
        .expect("Err creating client");

    if let Err(why) = client.start().await {
        println!("Client error: {:?}", why);
    }
}

