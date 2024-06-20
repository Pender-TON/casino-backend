package main

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

const channelID int64 = -1002030251950

const chatID int64 = -1002192560490

func isUserInChannel(b *gotgbot.Bot, userID int64, channelID int64) (bool, error) {
    chatMember, err := b.GetChatMember(channelID, userID, nil)

    if err != nil {
        return false, fmt.Errorf("failed to get chat member: %w", err)
    }

    return chatMember.GetStatus() != "left" && chatMember.GetStatus() != "kicked" && chatMember.GetStatus() != "restricted", nil
}

func isUserInChat(b *gotgbot.Bot, userID int64, chatID int64) (bool, error) {
    chatMember, err := b.GetChatMember(chatID, userID, nil)

    if err != nil {
        return false, fmt.Errorf("failed to get chat member: %w", err)
    }

    return chatMember.GetStatus() != "left" && chatMember.GetStatus() != "kicked" && chatMember.GetStatus() != "restricted", nil
}