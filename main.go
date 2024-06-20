package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type VerifyRequest struct {
    Hash string `json:"hash"`
    Data string `json:"data"`
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
	var err error
	err = godotenv.Load()
	if err != nil {
    	log.Fatal("Error loading .env file")
  	}	
    clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_CONNECTION_STRING"))
    client, err = mongo.Connect(ctx, clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Connected to MongoDB Atlas!")

	err = godotenv.Load()
	if err != nil {
    	log.Fatal("Error loading .env file")
  	}

    token := os.Getenv("BOT_TOKEN")
    if token == "" {
        panic("TOKEN environment variable is empty")
    }
  
    // Create bot from environment value.
    b, err := gotgbot.NewBot(token, nil)
    if err != nil {
        panic("failed to create new bot: " + err.Error())
    }
  
    // Create updater and dispatcher.
    dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
        // If an error is returned by a handler, log it and continue going.
        Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
            log.Println("an error occurred while handling update:", err.Error())
            return ext.DispatcherActionNoop
        },
        MaxRoutines: ext.DefaultMaxRoutines,
    })
    updater := ext.NewUpdater(dispatcher, nil)
  
    // /start command to introduce the bot
    dispatcher.AddHandler(handlers.NewCommand("start", start))
    dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("start"), startHandler))
  
    dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("tasks"), tasksHandler))
    dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("channel"), channelHandler))
    dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("chat"), chatHandler))
    dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("twitter"), twitterHandler))
    dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("checkChannel"), checkMembershipHandler))
    dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("checkChat"), chatMembershipHandler))
    dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("invite"), inviteHandler))
    dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("casino"), casinoHandler))
    dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("osaka"), osakaHandler))
    dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("checkOsaka"), osakaMembershipHandler))
  
    // Start receiving updates.
    err = updater.StartPolling(b, &ext.PollingOpts{
        DropPendingUpdates: true,
        GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
            Timeout: 9,
            RequestOpts: &gotgbot.RequestOpts{
                Timeout: time.Second * 10,
            },
        },
    })
    if err != nil {
        panic("failed to start polling: " + err.Error())
    }
    log.Printf("%s has been started...\n", b.User.Username)
  
    // Idle, to keep updates coming in, and avoid bot stopping.
    updater.Idle()
}
