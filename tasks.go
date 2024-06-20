package main

import (
	"fmt"
	"os"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)
func channel(b *gotgbot.Bot, ctx *ext.Context) error {
	row1 := []gotgbot.InlineKeyboardButton{
		{
			Text: 	   "Подпишись на канал",
			Url: "https://t.me/pender_official",
		},
	}
	row2 := []gotgbot.InlineKeyboardButton{
        {
            Text:         "Проверить",
            CallbackData: "checkChannel",
        },
    }
	row3 := []gotgbot.InlineKeyboardButton{
		{
			Text:         "Назад",
			CallbackData: "tasks",
		},
	}

    // Create the inline keyboard markup
    inlineKeyboardMarkup := gotgbot.InlineKeyboardMarkup{
        InlineKeyboard: [][]gotgbot.InlineKeyboardButton{row1, row2, row3},
    }

	chatID := ctx.EffectiveChat.Id
    msg, err := b.SendMessage(chatID, os.Getenv("CHANNEL_MSG"), &gotgbot.SendMessageOpts{
		ReplyMarkup: &inlineKeyboardMarkup,
	})
	if err != nil {
		return fmt.Errorf("failed to send source: %w", err)
	}

	lastMessageId = msg.MessageId
	return nil
}

func chat(b *gotgbot.Bot, ctx *ext.Context) error {
	row1 := []gotgbot.InlineKeyboardButton{
		{
			Text: 	   "Вступай в наш чат",
			Url: "https://t.me/+WCc9YAzTT8pmNmQy",
		},
	}
	row2 := []gotgbot.InlineKeyboardButton{
        {
            Text:         "Проверить",
            CallbackData: "checkChat",
        },
    }
	row3 := []gotgbot.InlineKeyboardButton{
		{
			Text:         "Назад",
			CallbackData: "tasks",
		},
	}

    // Create the inline keyboard markup
    inlineKeyboardMarkup := gotgbot.InlineKeyboardMarkup{
        InlineKeyboard: [][]gotgbot.InlineKeyboardButton{row1, row2, row3},
    }

	chatID := ctx.EffectiveChat.Id
    msg, err := b.SendMessage(chatID, os.Getenv("CHAT_MSG"), &gotgbot.SendMessageOpts{
		ReplyMarkup: &inlineKeyboardMarkup,
	})
	if err != nil {
		return fmt.Errorf("failed to send source: %w", err)
	}


	lastMessageId = msg.MessageId
	return nil
}

func twitter(b *gotgbot.Bot, ctx *ext.Context) error {
	row1 := []gotgbot.InlineKeyboardButton{
		{
			Text: 	   "X.com",
			Url: "https://x.com/pender_ton",
		},
	}
	row2 := []gotgbot.InlineKeyboardButton{
		{
			Text:         "Назад",
			CallbackData: "tasks",
		},
	}

    // Create the inline keyboard markup
    inlineKeyboardMarkup := gotgbot.InlineKeyboardMarkup{
        InlineKeyboard: [][]gotgbot.InlineKeyboardButton{row1, row2},
    }

	chatID := ctx.EffectiveChat.Id
    msg, err := b.SendMessage(chatID, os.Getenv("TWITTER_MSG"), &gotgbot.SendMessageOpts{
		ReplyMarkup: &inlineKeyboardMarkup,
	})
	if err != nil {
		return fmt.Errorf("failed to send source: %w", err)
	}


	lastMessageId = msg.MessageId
	return nil
}

func osaka(b *gotgbot.Bot, ctx *ext.Context) error {
	row1 := []gotgbot.InlineKeyboardButton{
		{
			Text: 	   "[0x_Knowledge]",
			Url: "https://t.me/knowledge_0x",
		},
	}
	row2 := []gotgbot.InlineKeyboardButton{
        {
            Text:         "Проверить",
            CallbackData: "checkOsaka",
        },
    }
	row3 := []gotgbot.InlineKeyboardButton{
		{
			Text:         "Назад",
			CallbackData: "tasks",
		},
	}

    // Create the inline keyboard markup
    inlineKeyboardMarkup := gotgbot.InlineKeyboardMarkup{
        InlineKeyboard: [][]gotgbot.InlineKeyboardButton{row1, row2, row3},
    }

	chatID := ctx.EffectiveChat.Id
    msg, err := b.SendMessage(chatID, os.Getenv("OSAKA_MSG"), &gotgbot.SendMessageOpts{
		ReplyMarkup: &inlineKeyboardMarkup,
	})
	if err != nil {
		return fmt.Errorf("failed to send source: %w", err)
	}


	lastMessageId = msg.MessageId
	return nil
}