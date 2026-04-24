## Using QuRe

[Return Documentation Index](./README.md)

---

- [User Accounts](#user-accounts)
- [Staff Accounts](#staff-accounts)
- [Admin Accounts](#admin-accounts)
- [Admin Tools](#admin-tools)

**For server commands, see [maintanance](./maintanance.md)**

## Temporary user account

...

## User accounts

...

## Staff accounts

Staff account has access to **Event Reservations** inspect tab, can view staff slots and sign up as staff to an event.


## Admin accounts

If an `admin` account doesn't exist, a new `admin` account is created automatically on server start.

**Check server log output for the `admin` credentials.**

- Username is `admin`
- Password is a random string of characters.

**The password should be changed on first login.**

> It is good policy to not share the `admin` account, but create individual accounts and promote them with the necessary roles. This way any policy violations can be tracked and offending accounts can be demoted or removed without affecting other administrators.

## Admin tools

You can access **account settings,** **reservation** and **admin tools** by clicking your user name at the top right corner of your screen. The last tab is **Admin Tools**.

- [Inspect vations](#reservations)
- [List All Users](#all-users)
- [Delete User](#deleting-a-user)
- [Change User Role](#changeing-a-users-role)
- [Create New Event](#creating-an-event)
- [Editing an Event](#editing-and-event)

### Reservations
The first tool is Reservations. You can use it to inspect the reservations of a given event. This is useful for verifying customers' reservations.

Click **Reservations** then click on the event to inspect. The list will populate with reservations sorted by their time.

You can search for a specific ID by typeing it in the search field.

**This action is GDPR safe. No PII is shown.**

### All Users
The second tab is **All Users**. It is only available to full **admin users**. **Opening the tab counts as PII access**. To comply with GDPR, there must be a good reason to access the list. **This action is logged**.

From the list, you can select a user. A detail card of the user is shown. From the card you can **delete the user** or **change its role**.

### Deleting a User
To delete a user, you must be a full admin user. Open **Settings -> Admin tools -> All Users** and select the user you want to delete. Click **Delete**. You are asked to enter your admin password to confirm the deletion. **This action is logged.** When a user is deleted by the user or an admin, all data related to the user is removed from the system.

### Changeing a User's Role
To change the role of a user, you must be a full admin user. Open **Settings -> Admin tools -> All Users** and select the user you want to modify. Click the ✏️ icon next to the role. Select the desired role from the drop down menu. You are asked to enter your admin password to confirm. **This action is logged.**

### Creating an Event
To create an event, you need to be logged in as an administrator. Click the card with the ➕ symbol on it to open a new event editor. All fields (except *short description*) must be filled to save the event. Event can be either **Published** or **Saved as a Draft**. A draft is visible only to administrators, while Published events are visible to anyone.

To add groups/timeslots to the event, set the group size. You can add more groups by clicking the plus (+) symbol next to the group. You can remove a group by pressing the minus (-) symbol next to it.

### Editing and Event
To create an event, you need to be logged in as an administrator. Select the event you want to edit and click Edit Event. Edit any fields you need to and either **Publish** or **Save as a Draft**. You can hide a published event by saving it as a draft, and you can publish a draft event by publishing it.
