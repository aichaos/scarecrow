# User Permissions

The Scarecrow chatbot has built-in support for user permissions. Currently the
only user role of any importance is the "Admin" role, which should consist only
of the bot's owner's usernames, and those of any other trusted users.

Users on the Admins list can access the global admin commands (see
[Admin Commands](./Admin-Commands.md)) as well as any other listener-specific
commands.

The admin status of a user can also be queried from within RiveScript by getting
the user variable named `isAdmin` -- this variable is always set to the correct
value before a reply is being fetched, so it will never become out of sync as
admins are added and removed.

## Admin Config

By default, only the local Console user has admin access, so you don't need to
worry about editing the Admin config file (`config/admins.json`) by hand; the
Console user is able to add additional admins using the `!op` command.

See [Admin Commands](./Admin-Commands.md) for a list of global admin commands.

## Reset Admins

If you want to reset the admin list (e.g. because you de-opped `CLI-console`)
and your other admin usernames), just delete the `config/admins.json` and
restart the bot. The `CLI-console` user has admin by default.

## Username Format

Because Scarecrow can chat on multiple platforms simultaneously, there would be
a risk of like-usernames being present on different platforms and being owned by
different users. For example, a Twitter user may have the name "example" and a
user on Slack might also have that as their nickname.

In order to uniquely identify a user across multiple platforms, the names are
formatted in a predictable way. The general format is as follows:

> `<listener prefix>-<username lowercased>`

Additionally, for Slack bots, the team's domain name is suffixed to the end of
their username, so a Slack user might look like `Slack-example@team.slack.com`
(the team name comes from the bots config file; it's technically an arbitrary
string but you should make it match the team's domain as a matter of
consistency).

Thus, to op a user the command might look like one of the following:

* `!op Slack-soandso@team.slack.com`
* `!op XMPP-admin@example.com`
* `!op Twitter-example`

## RiveScript Example

A simple example of querying whether the user is an admin:

```
+ am i an admin
* <get isAdmin> == true => Yes, you are an admin user.
- No, you are not an admin.
```
