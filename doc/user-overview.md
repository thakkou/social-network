# User Features Overview

## Authentication & Onboarding

| Feature | Description |
|---------|-------------|
| **Registration** | Sign up with email, password, first name, last name, date of birth. Optional: avatar (JPG/PNG/GIF), nickname, about me. |
| **Login** | Sign in with email or nickname. Session persists across page refreshes and browser tabs. |
| **Logout** | Available from any page via the sidebar link. |
| **Session Management** | Automatic session renewal. Stay logged in until explicitly logging out. |

## Profiles

| Feature | Description |
|---------|-------------|
| **Profile View** | See user info, avatar, bio, all posts, followers count, following count |
| **Public Profile** | Visible to everyone on the network |
| **Private Profile** | Visible only to accepted followers |
| **Privacy Toggle** | Switch between public/private anytime with confirmation pop-up |
| **Profile Editing** | Change nickname, about me, avatar from settings |

## Social Feed

| Feature | Description |
|---------|-------------|
| **Feed View** | Paginated feed with infinite scroll showing all posts you're allowed to see |
| **Post Creation** | Text + optional image/GIF + categories + privacy selection |
| **Post Privacy** | Public (everyone), Almost Private (followers only), Private (specific users) |
| **Post Reactions** | Like and dislike on posts |
| **Comments** | Comment on posts with optional images, reactions on comments |

## Follow System

| Feature | Description |
|---------|-------------|
| **Follow** | Follow a user via their profile page |
| **Unfollow** | Unfollow with confirmation pop-up |
| **Follow Requests** | Public users are followed automatically; private users receive a request |
| **Request Management** | Accept/reject follow requests from notifications |

## Groups

| Feature | Description |
|---------|-------------|
| **Create Group** | Set title, description, logo, background, optionally invite followers |
| **Group Discovery** | Browse all groups, search by name |
| **Join Group** | Public: instant join. Private: request to join. |
| **Group Posts** | Create and comment on posts visible only to members |
| **Group Events** | Create events with title, description, date/time, going/not going options |
| **Group Chat** | Real-time group chat room accessible from messages |
| **Member Management** | Admins can approve/reject requests and kick members |
| **Invitations** | Invite followers to join; members can also invite others |

## Messaging

| Feature | Description |
|---------|-------------|
| **Direct Messages** | Real-time 1-on-1 chat between users who follow each other |
| **Group Chat** | Real-time chat for all group members |
| **Emoji Support** | Full emoji picker in chat |
| **Online Status** | Green/gray dots show who's online |
| **Typing Indicators** | See when someone is typing |
| **New Conversations** | Start chats from the "discover" section or user profiles |

## Notifications

| Feature | Description |
|---------|-------------|
| **Real-time Alerts** | Instant notification delivery via WebSocket |
| **Notification Center** | `/notifications` page with all history |
| **Badge Count** | Unread count shown in header |
| **Inline Actions** | Accept/reject follow requests and group invites directly from notifications |
| **Bulk Actions** | Mark all as read, delete all |
| **Notification Types** | Follow requests, new followers, group invites, join requests, events, post reactions, comments |

## Search

| Feature | Description |
|---------|-------------|
| **User Search** | Search by name or nickname with avatar preview |
| **Group Search** | Search groups by title |
| **Live Results** | Results update as you type |
| **Quick Links** | Click results to navigate directly |

## UI/UX

| Feature | Description |
|---------|-------------|
| **Responsive Design** | Works on desktop and mobile (bottom nav on mobile) |
| **Dark Theme** | Full dark mode design |
| **Avatars** | User profile pictures shown throughout the app (sidebar, search, posts, comments, chat) |
| **Group Avatars** | Group logos displayed in sidebar, messages, and search |
| **Smooth Scrolling** | Animated scrolling effects in feeds |
| **Emoji Picker** | Consistent emoji picker across posts and chat |
| **Toast Notifications** | Non-intrusive toast messages for actions |
