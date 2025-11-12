# gigoseek
Concert recomendation Telegram bot service based on favorite artists. Service retrieves artists from your spotify account and searches their concerts on [Timepad](https://timepad.ru).

Link to bot: https://t.me/gigoseek_bot?text=%2Fstart

## Usage

1. **Start the bot**\
   Send `/start` to begin. The bot ask you to authenticate via Spotify and enter your city.

2. **Authenticate**\
   Use `/auth` to log in with your Spotify account.
   After successful authentication, you can use `/favorites` and `/concerts`.

3. **Available commands**

   ```
   /start         — start working with the bot
   /help          — show available commands
   /auth          — authenticate via Spotify
   /favorites     — show your favorite artists
   /concerts      — show upcoming concerts
   /change_city   — change your city
   ```

4. **Change city**
   Send `/change_city` and then enter your city name — the bot will find concerts near you.
