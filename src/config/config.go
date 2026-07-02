// SPDX-License-Identifier: AGPL-3.0-or-later
// SPDX-FileCopyrightText: 2026 GOTUNIX Networks <code@gotunix.net>
// SPDX-FileCopyrightText: 2026 Justin Ovens <code@gotunix.net>
// ----------------------------------------------------------------------------------------------- //
//                 $$$$$$\   $$$$$$\ $$$$$$$$\ $$\   $$\ $$\   $$\ $$$$$$\ $$\   $$\               //
//                $$  __$$\ $$  __$$\\__$$  __|$$ |  $$ |$$$\  $$ |\_$$  _|$$ |  $$ |              //
//                $$ /  \__|$$ /  $$ |  $$ |   $$ |  $$ |$$$$\ $$ |  $$ |  \$$\ $$  |              //
//                $$ |$$$$\ $$ |  $$ |  $$ |   $$ |  $$ |$$ $$\$$ |  $$ |   \$$$$  /               //
//                $$ |\_$$ |$$ |  $$ |  $$ |   $$ |  $$ |$$ \$$$$ |  $$ |   $$  $$<                //
//                $$ |  $$ |$$ |  $$ |  $$ |   $$ |  $$ |$$ |\$$$ |  $$ |  $$  /\$$\               //
//                \$$$$$$  | $$$$$$  |  $$ |   \$$$$$$  |$$ | \$$ |$$$$$$\ $$ /  $$ |              //
//                 \______/  \______/   \__|    \______/ \__|  \__|\______|\__|  \__|              //
// ----------------------------------------------------------------------------------------------- //
// Copyright (C) GOTUNIX Networks                                                                  //
// Copyright (C) Justin Ovens                                                                      //
// ----------------------------------------------------------------------------------------------- //
// This program is free software: you can redistribute it and/or modify                            //
// it under the terms of the GNU Affero General Public License as                                  //
// published by the Free Software Foundation, either version 3 of the                              //
// License, or (at your option) any later version.                                                 //
//                                                                                                 //
// This program is distributed in the hope that it will be useful,                                 //
// but WITHOUT ANY WARRANTY; without even the implied warranty of                                  //
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the                                   //
// GNU Affero General Public License for more details.                                             //
//                                                                                                 //
// You should have received a copy of the GNU Affero General Public License                        //
// along with this program.  If not, see <https://www.gnu.org/licenses/>.                          //
// ----------------------------------------------------------------------------------------------- //

package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	DiscordBotToken string
	AllowedUserIDs  []string
	NotifyChannelID string
	NotifyUserIDs   []string
	ServerHost      string
	ServerPort      string
	WebhookSecret   string
)

func Load() {
	_ = godotenv.Load() // Ignore error if .env is missing

	DiscordBotToken = os.Getenv("DISCORD_BOT_TOKEN")
	AllowedUserIDs = parseList(os.Getenv("ALLOWED_USER_IDS"))
	NotifyChannelID = os.Getenv("NOTIFY_CHANNEL_ID")
	NotifyUserIDs = parseList(os.Getenv("NOTIFY_USER_IDS"))

	ServerHost = os.Getenv("SERVER_HOST")
	if ServerHost == "" {
		ServerHost = "0.0.0.0"
	}

	ServerPort = os.Getenv("SERVER_PORT")
	if ServerPort == "" {
		ServerPort = "8765"
	}

	WebhookSecret = os.Getenv("WEBHOOK_SECRET")

	Validate()
}

func parseList(val string) []string {
	var list []string
	if val == "" {
		return list
	}
	parts := strings.Split(val, ",")
	for _, p := range parts {
		if strings.TrimSpace(p) != "" {
			list = append(list, strings.TrimSpace(p))
		}
	}
	return list
}

func Validate() {
	if DiscordBotToken == "" {
		log.Fatal("DISCORD_BOT_TOKEN is not set. Add it to your environment or .env file.")
	}
	if NotifyChannelID == "" && len(NotifyUserIDs) == 0 && len(AllowedUserIDs) == 0 {
		log.Fatal("Set at least one of NOTIFY_CHANNEL_ID, NOTIFY_USER_IDS, or ALLOWED_USER_IDS so the bot knows where to send webhook notifications.")
	}
}
