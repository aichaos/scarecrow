# Scarecrow

Scarecrow is a chatbot written in Go. It connects to Slack and can be chatted
with in your terminal window, and it will probably be updated to connect to more
things in the future too.

It uses [RiveScript](http://www.rivescript.com/) as its brain back-end, it
remembers information about the people it chats with, keeps log files, etc.

# Features

* [Slack](https://www.slack.com/) integration
  * Users can chat with it over direct message and carry on a conversation
  * It can join public channels where it will sit in silence until a user talks
    directly to it, either by at-mentioning its username or starting a message
    with its name.
* Goroutines are spawned for each individual bot connection, so you can run
  multiple instances of Slack bots from one program.
* Chat with the bot on the console, too

# Install and Build

I don't have Godep set up, so install the dependencies manually for now.

```bash
go get github.com/nlopes/slack               # commit 9153359e4c6e3dd6
go get github.com/jprichardson/readline-go   # commit 69d39989b17ebcfb
```

Build it with `make build`

Run it with `./scarecrow`

Command line options are pretty basic: `--debug` for debug mode and
`--version` to get the version number.

# Configuration

The bot is configured through JSON files in the `config` folder. You'll have to
edit the JSON files by hand, sorry. :frowning:

There is an example config file in `config/bots-sample.json` -- simply copy this
file and name it `bots.json` and edit it to configure your bot.

Here is an abridged version (with comments) of the config file:

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
      "api_token": "XXXX-NOT-A-REAL-TOKEN-XXXX", // Slack Bot API token
      "username": "scarecrow" // Enter the Slack bot's real username here
    },
    {
      "type": "Console",
      "enabled": true,
      "username": "Scarecrow" // This username is shown in the console when the
                              // bot sends you a reply. It can be anything you
                              // want.
    }
  ]
}
```

# To Do

This is my first real chatbot written in Go and it isn't very perfect...

Here's my list of things I want to improve on it:

* Make the Listener system implement Golang interfaces.
  * There's `src/listeners/listeners.go` that defines an interface, but it's
    not being used by the actual listeners. Instead it's just a reference for
    developers (me) to know what functions the listener should have.
  * I tried making it an interface originally but it was proving to be pretty
    difficult (errors like "Can't use a `*SlackListener` in place of a
    `*Listener` and stuff...")
  * As a consequence, the main code (`src/scarecrow.go`) needs special logic and
    extra lines of code to handle the Slack vs. Console listeners, instead of
    (ideally) being able to just have an array of listeners it can loop over
    and deal with.
* Revisit how I'm using Goroutines and see if it's the optimal way of doing it.
  * See below for the Goroutine layout of this app.

## Goroutines

The main thread does all the work up until the `Start()` function is called.

Then, for each active Listener:

* Two channels are created:
  * `ReplyRequest`: Messages from the listener that need a RiveScript reply.
  * `ReplyAnswer`: Responses from RiveScript going back to the listener.
* A Goroutine (`ManageRequestChannel`) is spawned which watches the
  `ReplyRequest` channel, to get a reply and send back the response over the
  other channel.
* The listener itself also spawns number of Goroutines:
  * The Slack listener spawns one from the Slack RTM API module (unavoidable)
    and also one for maintaining its own `MainLoop()`
    * The `MainLoop` checks for messages from the Slack RTM channel *and* the
      `ReplyAnswer` channel (for sending replies to the users).
  * The Console listener spawns two Goroutines: one to read from the terminal
    (`os.Stdin`) and write to a message channel, and another to run its
    `MainLoop()` -- which, like the Slack listener, reads from the readline
    channel and the `ReplyAnswer` channel.

# License

```
Admiral Chatbot
Copyright (C) 2015  Noah Petherbridge

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License along
with this program; if not, write to the Free Software Foundation, Inc.,
51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
```

# See Also

RiveScript's official homepage, <http://www.rivescript.com/>

The RiveScript Go module, <https://github.com/aichaos/rivescript-go>
