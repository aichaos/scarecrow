# Configuration

Scarecrow is configured through some JSON files in the `config/` folder. The
config format may change in the future, as JSON isn't the most user-friendly
format and lacks support for comments, which would allow a config file to be
self-documenting.

# Bot Configuration (bots.json)

This is the config file for the bots and their connections to various messaging
services.

Look at `bots-sample.json` for an example, and create a file named `bots.json`
and fill in the information for your bots. You can simply copy/paste the
sample file and then edit it.

The file looks something like this (with comments):

```javascript
{
  // Personality: configure global details about your bot.
  "personality": {
    // Give your bot a name. Currently isn't used anywhere...
    "name": "Scarecrow",

    // Brain configuration for your bot. Currently `backend` isn't used but this
    // bot may support more brains than just RiveScript in the future.
    "brain": {
      "backend": "RiveScript",

      // This is the path on disk to your RiveScript *.rive files.
      "replies": "./replies/standard"
    }
  },

  // Listeners: array of interfaces for how people can communicate with your
  // chatbot. You can have many listeners here. The default example config only
  // has one entry for each type of listener but you can have many entries
  // (although having multiple Consoles may get messy and confusing...)
  "listeners": [
    {
      "type": "Slack",   // Each Listener has a type
      "enabled": false,  // Set to `true` to enable this listener on start-up
      "settings": {      // Listener-specific configuration
          "api_token": "XXXX-NOT-A-REAL-TOKEN-XXXX", // Slack Bot API token
          "username": "scarecrow",    // The Slack bot's real username
          "team": "example.slack.com" // Your Slack team's domain
      }
    },
    {
      "type": "Console",
      "enabled": true,
      "settings": {
          "username": "Scarecrow" // This username is shown in the console when
                                  // the bot sends you a reply. It can be
                                  // anything you want.
      }
    }
  ]
}
```

## Listener Settings

Here are the specific `settings` for each type of listener.

**NOTE: All settings have string values.** So for the boolean settings, you
write the word `"true"` or `"false"` in quotes. Numbers are also quoted, etc.

### Console

* `username`: This username is shown in the terminal when the bot sends you a
  response.

### Slack

* `api_token`: This is your Slack Bot API token.
* `username`: The bot's username in Slack. You should make sure this accurately
  matches the bot's real username, as Scarecrow uses this name to look up its
  own user ID to identify when somebody at-mentions it in a Slack channel.
* `team`: The Slack team's domain. This is an arbitrary string but it should
  look like `teamname.slack.com`. The team is appended onto the end of Slack
  usernames for the purpose of keeping like-usernames in different Slack teams
  separate under the hood, for user variable, logging and admin purposes.

### XMPP

Required fields:

* `username`: The bot's XMPP username, like `scarecrow@jabber.com`
* `password`: The bot's XMPP password.
* `server`: The name of the XMPP server, like `example.com`
* `port`: The port number to use on the XMPP server, like `5222`

Optional fields:

* `debug`: If `"true"`, debugging is enabled in the XMPP module.
* `notls`: If `"true"`, TLS will not be used with the initial connection to the
  XMPP server (it can be combined with `starttls: true` if the server supports
  StartTLS)
* `starttls`: If `"true"`, the connection will be upgraded with StartTLS.
* `tls-no-verify`: If `"true"`, the TLS server certificate is not verified
  (this setting will be needed if you use a self-signed certificate on the
  XMPP server).

## Listener Examples

Full example config snippets for certain listeners.

### Google Hangouts via XMPP

```json
{
    "type": "XMPP",
    "enabled": true,
    "settings": {
        "username": "scarecrow.bot@gmail.com",
        "password": "XXX-PASSWORD-XXX",
        "server": "talk.google.com",
        "port": "443"
    }
}
```

### XMPP with StartTLS

Scenario: your XMPP server doesn't speak TLS directly (traditionally was done on
port `5223`), but listens on the clear-text port `5222` and supports (or
requires) StartTLS to upgrade the connection to have encryption.

```json
{
    "type": "XMPP",
    "enabled": true,
    "settings": {
        "username": "scarecrow@jabber.com",
        "password": "XXX-PASSWORD-XXX",
        "server": "jabber.com",
        "port": "5222",
        "notls": "true",
        "starttls": "true"
    }
}
```

Note: we use the `notls` ("No TLS") option because the server connection doesn't
begin as TLS. The `starttls` option then tells it to use StartTLS over the
non-secure connection.
