# Configuration Files

# Bot Configuration (bots.json)

This is the config file for the bots and their connections to various messaging
services.

Look at `bots-sample.json` for an example, and create a file named `bots.json`
and fill in the information for your bots. You can simply copy/paste the
sample file and then edit it.

The file looks something like this, with comments:

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
  // has one for Slack and one for Console but you can have multiple for each
  // (although having multiple Consoles may get messy and confusing...)
  "listeners": [
    {
      "type": "Slack",   // Each Listener has a type
      "enabled": false,  // Set to `true` to enable this listener on start-up
      "settings": {      // Listener-specific configuration
          "api_token": "XXXX-NOT-A-REAL-TOKEN-XXXX", // Slack Bot API token
          "username": "scarecrow" // Enter the Slack bot's real username here
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

### XMPP

Required fields:

* `username`: The bot's XMPP username, like `scarecrow@jabber.com`
* `password`: The bot's XMPP password.
* `server`: The name of the XMPP server, like `example.com`
* `port`: The port number to use on the XMPP server, like `5222`

Optional fields:

* `debug`: If `"true"`, debugging is enabled in the XMPP module.
* `tls-disable`: If `"true"`, TLS will not be used with the XMPP server.
* `tls-no-verify`: If `"true"`, the TLS server certificate is not verified
  (this setting will be needed if you use a self-signed certificate on the
  XMPP server).
