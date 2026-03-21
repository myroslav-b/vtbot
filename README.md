# VTBot (VirusTotal Telegram Bot)

[Українська версія (Ukrainian version) 🇺🇦](README_UA.md)

A Telegram bot that checks incoming files using the [VirusTotal API](https://www.virustotal.com/). The bot searches for files by their SHA-256 hash or uploads them for a new analysis if they are not found in the database.

## Features

- **File hash check (SHA-256):** Instantly returns a report if the file is already in the VirusTotal database.
- **File Uploading:** Automatically uploads unknown files for scanning.
- **Queue System:** Safely handles multiple file validation requests concurrently.
- **Rate Limit Handling:** Accommodates the free VirusTotal API tier limits (4 requests/minute) with built-in request intervals.
- **Size Limit Check:** Rejects files over 20MB to comply with free-tier capabilities.
- **Docker Support:** Ready to be deployed quickly using Docker and Docker Compose.

## Prerequisites

To run this bot, you will need:
- A Telegram Bot API Token (from [BotFather](https://t.me/BotFather)).
- A VirusTotal API Key.
- Docker & Docker Compose (optional, for containerized deployment).

## Installation & Running

### Using Docker Compose (Recommended)

1. Clone the repository to your machine.
2. Ensure you have Docker and Docker Compose installed.
3. Configure your API keys. You can do this by setting environment variables in your terminal, or by creating an `.env` file in the root directory:
   ```env
   TELEGRAM_TOKEN=your_telegram_bot_token
   VT_API_KEY=your_virustotal_api_key
   ```
4. Build and start the bot:
   ```bash
   docker-compose up -d --build
   ```

### Manual Build

If you prefer not to use Docker, you need to have Go installed (`1.25` or newer):

1. Export required environment variables:
   ```bash
   export TELEGRAM_TOKEN="your_telegram_bot_token"
   export VT_API_KEY="your_virustotal_api_key"
   ```
2. Build the project:
   ```bash
   go build -o bot ./cmd/bot
   ```
3. Run the application:
   ```bash
   ./bot
   ```

## Usage

Simply start a chat with the bot and send any document or file. The bot will acknowledge the submission, analyze the file, and send back a report summarizing the `malicious`, `suspicious`, and `harmless` detection scores.

### Using in Groups and Channels

To allow the bot to read and check files sent in a group or channel, you must do one of the following:
- **Disable Privacy Mode** for your bot via [@BotFather](https://t.me/BotFather) (send `/setprivacy`, choose your bot, and set to `Disable`).
- OR **Promote the bot to Administrator** in the group/channel. 

*Note: The standard Telegram Bot API has a 20 MB file size limit for downloading files. Files larger than 20 MB cannot be downloaded or checked by the bot.*
