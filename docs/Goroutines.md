# Goroutines

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
