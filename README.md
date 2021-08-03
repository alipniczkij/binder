# Binder

---
There is a problem that the built-in gitlab notifier spams events not only about opening an issue, but also about closing, reopening it and so on. 
There was a desire to notify only about the events of the opening of the issue. 
Therefore, this service was written.

---
The web server that validates events from gitlab and performs the necessary notifications in Slack channels.

In the current version, notifications are sent only when an issue is opened.

---
To use commands in slack you need to create slack bot with commands below and permissions, and add it to Slack channels

## Start service
You need to fill the config file and create bbolt storage file before start.
```
go mod download
go run cmd/binder/main.go
```

## Commands

- /subscribe {group or project name} (ex. alipniczkij/binder)
  
  To subscribe a specific channel to notifications from the gitlab group or project, you must:
    1. Add your bot to channel/chat
    2. Add a webhook in gitlab's group or project. The URL is where the binder works. Trigger - Issue Events
    3. Send the above command at the Slack channel/chat

**IMPORTANT**

If you want to subscribe your channel to notifications for a specific project (for example, your_group/binder), you need to check whether there is a webhook on the group to which the project belongs.
If there is a webhook for group, then you do NOT need to add a webhook for the project. If you add a webhook to the project as well as group, notifications will be duplicated in all channels

- /unbsubscribe {group or project name}

  Unsubscribes your channel/chat from notifications for the specified group or project

- /label {labels names separated by escape}

- /unlabel {labels names separated by escape}

  The commands regulate notifications by labels, linking to specific channels/chats (not gitlab groups/projects)
  
  Implemented logic:
  1. If you used unlabel, then all issues will be sent to the channel, except those that will have unwanted labels
  2. If you used a label, only those issues that have the desired labels will be sent to the channel. Anything else won't be sent
  3. If you haven't used anything, then all the issues will come to your channel

The /label and /unlabel commands have the -delete subcommand, which allows you to delete labels previously added by the corresponding commands.
Example:
`/label -delete ux`  
- /list

  The command shows all subscriptions and labels/unlabels of group
  ```
  Subscribed: ...
  Labels: ...
  Unlabels: ...
  ```

---

## Bbolt storage

There is structure for storage:

- Buckets - command name
- Keys:
  1. For /subscribe and /unsubscribe commands – group or project name
  2. For /label and /unlabel commands – Slack channel id
- Values:
  1. For /subscribe and /unsubscribe commands – slice of Slack channels IDs
  2. For /label and /unlabel commands – slice of labels names

## TODO

- Add commentaries
- Tests
- Gitlab webhook token support
- Docker
