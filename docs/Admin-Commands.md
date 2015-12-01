# Admin Commands

There are some built-in commands that can only be accessed by admin users
(see [User Permissions](./User-Permissions.md)). These include some global
commands that work everywhere, and some listener-specific commands (for example,
an IRC listener may have commands to instruct the bot to join or leave a
channel).

## Global Commands

* `!reload`

  Reload the bot's RiveScript brain without having to completely reboot the bot.
  This is useful for actively developing its brain so you can quickly load your
  new changes.

* `!op <user-id>`

  This adds `user-id` to the admins list. Note that the user ID must match the
  bot's convention for uniquely identifying users (e.g. with the listener
  prefix). An easy way to find the correct user ID is to check the chat logs
  and copy the user's name from there.

* `!deop <user-id>`

  This removes `user-id` from the admins list. The user ID must be formatted
  the same way as for the `!op` command.
