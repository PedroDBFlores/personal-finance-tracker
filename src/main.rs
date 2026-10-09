use std::fmt::Error;
use tokio::main;

mod cmd;

#[main]
async fn main() -> Result<(), Error> {
    println!("Hello from Tracker!");
    Ok(())
}
