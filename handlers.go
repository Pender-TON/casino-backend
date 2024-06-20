package main

import (
	"fmt"
	"os"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

var lastMessageId int64
// start introduces the bot.
func start(b *gotgbot.Bot, ctx *ext.Context) error {
	row1 := []gotgbot.InlineKeyboardButton{
		{
			Text: 	   "Задания",
			CallbackData: "tasks",
		},
	}
	row2 := []gotgbot.InlineKeyboardButton{
        {
            Text:         "Pender Airdrop",
            CallbackData: "casino",
        },
    }
	row3 := []gotgbot.InlineKeyboardButton{
		{
			Text:	   "Пригласить друзей",
			CallbackData: "invite",
		},
	}

    // Create the inline keyboard markup
    inlineKeyboardMarkup := gotgbot.InlineKeyboardMarkup{
        InlineKeyboard: [][]gotgbot.InlineKeyboardButton{row1, row2, row3},
    }

    // Send a message with the inline keyboard
	chatID := ctx.EffectiveChat.Id
	photoURL := "https://imgur.com/a/euIocv3"
    msg, err := b.SendPhoto(chatID, photoURL, &gotgbot.SendPhotoOpts{
        Caption: os.Getenv("WELCOME_MESSAGE"),
		ReplyMarkup: &inlineKeyboardMarkup,
    })
	if err != nil {
    	return err
	}

	lastMessageId = msg.MessageId
	return nil
}

func invite(b *gotgbot.Bot, ctx *ext.Context) error {
	row1 := []gotgbot.InlineKeyboardButton{
		{
			Text:         "Назад",
			CallbackData: "start",
		},
	}

	inlineKeyboardMarkup := gotgbot.InlineKeyboardMarkup{
        InlineKeyboard: [][]gotgbot.InlineKeyboardButton{row1},
    }
	chatID := ctx.EffectiveChat.Id
	refLink := fmt.Sprintf("Вот твоя рефка:\n\nhttps://t.me/pender_referrals_bot?start=%d", ctx.EffectiveUser.Id)
    msg, err := b.SendMessage(chatID, refLink, &gotgbot.SendMessageOpts{
		ReplyMarkup: &inlineKeyboardMarkup,
	})
	if err != nil {
		return fmt.Errorf("failed to send source: %w", err)
	}


	lastMessageId = msg.MessageId
	return nil
}

func casino(b *gotgbot.Bot, ctx *ext.Context) error {
	row1 := []gotgbot.InlineKeyboardButton{
		{
			Text:         "Pender Casino",
			Url: "https://t.me/pender_referrals_bot/casino",
		},
	}
	row2 := []gotgbot.InlineKeyboardButton{
		{
			Text:         "Назад",
			CallbackData: "start",
		},
	}

	inlineKeyboardMarkup := gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{row1, row2},
	}
	chatID := ctx.EffectiveChat.Id
	photoURL := "https://imgur.com/a/dWwL23z"
    msg, err := b.SendPhoto(chatID, photoURL, &gotgbot.SendPhotoOpts{
        Caption: os.Getenv("AIR_MSG"),
		ReplyMarkup: &inlineKeyboardMarkup,
    })
	if err != nil {
		return fmt.Errorf("failed to send source: %w", err)
	}
	lastMessageId = msg.MessageId
	return nil
}

func tasks(b *gotgbot.Bot, ctx *ext.Context) error {
	row1 := []gotgbot.InlineKeyboardButton{
		{
			Text: 	   "Подпишись на канал",
			CallbackData: "channel",
		},
	}
	row2 := []gotgbot.InlineKeyboardButton{
        {
            Text:         "Вступай в наш чат",
            CallbackData: "chat",
        },
    }
	row3 := []gotgbot.InlineKeyboardButton{
		{
			Text:         "Твиттер (Х)",
			CallbackData: "twitter",
		},
	}
	row4 := []gotgbot.InlineKeyboardButton{
		{
			Text:         "[0x Knowledge]",
			CallbackData: "osaka",
		},
	}
	row5 := []gotgbot.InlineKeyboardButton{
		{
			Text:         "Назад",
			CallbackData: "start",
		},
	}

    // Create the inline keyboard markup
    inlineKeyboardMarkup := gotgbot.InlineKeyboardMarkup{
        InlineKeyboard: [][]gotgbot.InlineKeyboardButton{row1, row2, row3, row4, row5},
    }

	chatID := ctx.EffectiveChat.Id
	photoURL := "https://imgur.com/a/nYMQl3c"
    msg, err := b.SendPhoto(chatID, photoURL, &gotgbot.SendPhotoOpts{
        Caption: os.Getenv("TASKS_MESSAGE"),
		ReplyMarkup: &inlineKeyboardMarkup,
	})
	if err != nil {
		return fmt.Errorf("failed to send source: %w", err)
	}

	lastMessageId = msg.MessageId
	return nil
}

func clear(b *gotgbot.Bot, ctx *ext.Context) error {

	_, err := b.DeleteMessage(ctx.EffectiveChat.Id, lastMessageId, nil)
	if err != nil {
		start(b, ctx)
		return fmt.Errorf("failed to delete previous message: %w", err)
	}
	return nil
}

func startHandler(b *gotgbot.Bot, ctx *ext.Context) error {
    
	err := clear(b, ctx)
    if err != nil {
        return err
    }
	
	err = start(b, ctx)
    if err != nil {
        return err
    }

    return nil
}

func tasksHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	
	err := clear(b, ctx)
	if err != nil {
		return err
	}
	
	err = tasks(b, ctx)
	if err != nil {
		return err
	}

	return nil
}

func channelHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	
	err := clear(b, ctx)
	if err != nil {
		return err
	}

	err = channel(b, ctx)
	if err != nil {
		return err
	}

	return nil
}
func inviteHandler(b *gotgbot.Bot, ctx *ext.Context) error {

	err := clear(b, ctx)
	if err != nil {
		return err
	}

	err = invite(b, ctx)
	if err != nil {
		return err
	}

	return nil

}

func chatHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	
	err := clear(b, ctx)
	if err != nil {
		return err
	}

	err = chat(b, ctx)
	if err != nil {
		return err
	}

	return nil
}

func twitterHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	
	err := clear(b, ctx)
	if err != nil {
		return err
	}

	err = twitter(b, ctx)
	if err != nil {
		return err
	}

	return nil
}

func osakaHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	
	err := clear(b, ctx)
	if err != nil {
		return err
	}

	err = osaka(b, ctx)
	if err != nil {
		return err
	}

	return nil
}

func casinoHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	
	err := clear(b, ctx)
	if err != nil {
		return err
	}

	err = casino(b, ctx)
	if err != nil {
		return err
	}

	return nil
}

func checkMembershipHandler(b *gotgbot.Bot, ctx *ext.Context) error {
    userID := ctx.CallbackQuery.From.Id
    isMember, err := isUserInChannel(b, userID, channelID)
    if err != nil {
        return err
    }

    if isMember {
        _, err = ctx.EffectiveMessage.Reply(b, "Задание выполнено!", nil)
    } else {
        _, err = ctx.EffectiveMessage.Reply(b, "Ты не подписан на канал", nil)
    }

    return err
}

func chatMembershipHandler(b *gotgbot.Bot, ctx *ext.Context) error {
    userID := ctx.CallbackQuery.From.Id
    isMember, err := isUserInChannel(b, userID, chatID)
    if err != nil {
        return err
    }

    if isMember {
        _, err = ctx.EffectiveMessage.Reply(b, "Задание выполнено!", nil)
    } else {
        _, err = ctx.EffectiveMessage.Reply(b, "Ты не cостоишь в чате!", nil)
    }

    return err
}

func osakaMembershipHandler(b *gotgbot.Bot, ctx *ext.Context) error {
    userID := ctx.CallbackQuery.From.Id
    isMember, err := isUserInChannel(b, userID, osakaID)
    if err != nil {
        return err
    }

    if isMember {
        _, err = ctx.EffectiveMessage.Reply(b, "Задание выполнено!", nil)
    } else {
        _, err = ctx.EffectiveMessage.Reply(b, "Подпишись на канал!", nil)
    }

    return err
}