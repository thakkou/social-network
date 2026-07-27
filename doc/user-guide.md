# User Guide

## Getting Started

### Registration

1. Navigate to `/register`
2. Fill in required fields: **First Name**, **Last Name**, **Email**, **Password**, **Date of Birth**
3. Optionally add: **Nickname**, **About Me**, **Avatar** (JPEG/PNG/GIF, max 1MB)
4. Click **"create account →"**

### Login

1. Navigate to `/login`
2. Enter your **email** or **nickname** and **password**
3. Click **"sign in →"**

> **Quick Login:** In development mode, one-click login buttons for test users (Alice, Bob, Chloe, Farid, Isabella, Jack) are available on the login page.

## Navigation

The app has a persistent sidebar on the left with these sections:

### Navigate
- **Feed** (`/`) — Main feed with all public/almost-private posts
- **Profile** (`/profile`) — Your own profile page
- **Groups** (`/groups`) — All your groups and group discovery
- **Messages** (`/messages`) — Direct and group messaging
- **Notifications** (`/notifications`) — All your notifications
- **Settings** (`/settings`) — Account settings

### My Groups
- Shows groups you've joined with their logo/avatar and name
- Click a group to open its page
- **"new group"** button to create a new group

### Following / Followers
- Lists users you follow and who follow you
- Click a user to visit their profile
- Shows avatar images when available

## Profile

### Your Profile
- View your posts, followers, following, and groups
- Toggle **public/private** profile (with confirmation pop-up)
- Edit your nickname, about me, and avatar in Settings

### Other Users' Profiles
- View their posts, followers, and following
- **Follow/Unfollow** — click to follow; if private, a request is sent
- Private profiles are only visible to accepted followers

## Posts

### Creating a Post

1. On the feed page, use the post composer at the top
2. Add a **title**, **content**, and optionally an **image/GIF**
3. Select **categories** (optional)
4. Choose **privacy level**:
   - **Public** — everyone can see it
   - **Almost Private** — only your followers can see it
   - **Private** — only specific users you choose can see it
5. Click **"Post"**

### Interacting with Posts
- **Like / Dislike** — click the thumbs up/down buttons
- **Comment** — click the comment icon to expand comments
- **Delete** — delete your own posts (with confirmation)

## Groups

### Creating a Group

1. Click **"new group"** in the sidebar or go to `/groups`
2. Enter a **title** and **description**
3. Optionally upload a **logo** and **background** image
4. Optionally invite followers to join immediately

### Joining a Group

- **Public groups** — click "Join" to join immediately
- **Private groups** — click "Request to Join"; the creator will approve or reject

### Group Features

Once you're a member, you can:

- **Create posts** — share text and images visible only to group members
- **Create events** — set a title, description, date/time, and RSVP options (going/not going)
- **Chat** — use the group chat room for real-time messaging (link in the messages sidebar)
- **Invite members** — invite your followers to join
- **Leave group** — if you no longer want to be a member

### Group Admin (Creator Only)

- Approve/reject join requests
- Kick members
- Update group details
- Delete posts and events

## Messaging

### Direct Messages

1. Go to **Messages** (`/messages`) or use the right sidebar
2. Select a conversation or click a user from the "discover" section
3. Type your message and press Enter
4. Send **emojis** using the emoji picker

**Note:** You can only message users that follow you or that you follow.

### Group Chat

1. Join a group
2. Go to Messages — shared groups appear in the right sidebar
3. Click the group to open the group chat room
4. Messages are delivered in real-time to all online members

### Online Status

A green dot indicates a user is online (connected via WebSocket). A gray dot means offline.

## Notifications

Notifications appear in real-time (bell icon in the header). Click the bell or go to `/notifications` to view all.

### Notification Types

| Type | What Happens |
|------|-------------|
| Follow Request | Someone with a private profile wants to follow you — accept or reject |
| New Follower | Someone started following you (public profile or request accepted) |
| Follow Accepted | Your follow request was accepted |
| Group Invitation | You're invited to join a group — accept or reject |
| Group Join Request | Someone wants to join your group — approve or reject |
| Group Event | A new event was created in one of your groups |
| Post Reaction | Someone liked or disliked your post |
| Comment | Someone commented on your post |

## Settings

Available at `/settings`:

- **Update Profile** — change nickname, about me, avatar
- **Privacy** — toggle profile between public and private (with confirmation)
- **Your Groups** — view groups you've created and manage them

## Search

Use the search bar in the header to find:
- **Users** — search by name or nickname
- **Groups** — search by group title

Results appear as you type with avatars and quick-links to profiles or groups.

## Docker Deployment

```bash
# Build and start everything
./script.sh

# Or manually:
docker compose up --build

# Access:
# Frontend: http://localhost:3000
# Backend:  http://localhost:8080
```
