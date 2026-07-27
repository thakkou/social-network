# Requirements Analysis

## Functional Requirements

### ✅ Authentication
| ID | Requirement | Status |
|----|-------------|--------|
| FR-01 | Users can register with email, password, first name, last name, date of birth | ✅ |
| FR-02 | Optional fields: avatar (JPEG/PNG/GIF), nickname, about me | ✅ |
| FR-03 | Users can log in with email/nickname and password | ✅ |
| FR-04 | Invalid credentials show appropriate error message | ✅ |
| FR-05 | Duplicate email/nickname is detected during registration | ✅ |
| FR-06 | Sessions persist across page refreshes and browser restarts | ✅ |
| FR-07 | Logout option is always accessible | ✅ |

### ✅ Followers
| ID | Requirement | Status |
|----|-------------|--------|
| FR-08 | Users can follow and unfollow other users | ✅ |
| FR-09 | Following a public user is instantaneous | ✅ |
| FR-10 | Following a private user sends a follow request | ✅ |
| FR-11 | Recipient can accept or decline follow requests | ✅ |
| FR-12 | Unfollow requires confirmation | ✅ |

### ✅ Profile
| ID | Requirement | Status |
|----|-------------|--------|
| FR-13 | Profile displays all registration info (except password) | ✅ |
| FR-14 | Profile displays all posts by the user | ✅ |
| FR-15 | Profile displays followers and following lists | ✅ |
| FR-16 | Users can toggle between public and private profile | ✅ |
| FR-17 | Private profiles are hidden from non-followers | ✅ |
| FR-18 | Public profiles are visible to everyone | ✅ |
| FR-19 | Followed users can see each other's private profiles | ✅ |

### ✅ Posts
| ID | Requirement | Status |
|----|-------------|--------|
| FR-20 | Authenticated users can create posts | ✅ |
| FR-21 | Posts can include images (JPG/PNG) or GIFs | ✅ |
| FR-22 | Posts have privacy levels: public, almost private, private | ✅ |
| FR-23 | Private posts allow specifying which users can see them | ✅ |
| FR-24 | Authenticated users can comment on posts | ✅ |
| FR-25 | Comments can include images or GIFs | ✅ |
| FR-26 | Users can like/dislike posts and comments | ✅ |
| FR-27 | Users can delete their own posts and comments | ✅ |

### ✅ Groups
| ID | Requirement | Status |
|----|-------------|--------|
| FR-28 | Users can create groups with title and description | ✅ |
| FR-29 | Group creators can invite followers to join | ✅ |
| FR-30 | Invited users receive a notification and can accept/reject | ✅ |
| FR-31 | Users can request to join groups | ✅ |
| FR-32 | Group creators approve/reject join requests | ✅ |
| FR-33 | Group members can invite other users | ✅ |
| FR-34 | Members can create posts and comments within groups | ✅ |
| FR-35 | Group posts/comments are only visible to members | ✅ |
| FR-36 | Members can create events with title, description, date/time | ✅ |
| FR-37 | Events have at least 2 options: going, not going | ✅ |
| FR-38 | Users can vote on event options | ✅ |
| FR-39 | Group admins can kick members and delete content | ✅ |

### ✅ Chat
| ID | Requirement | Status |
|----|-------------|--------|
| FR-40 | Users can send private messages to followers/following | ✅ |
| FR-41 | Messages are delivered in real-time via WebSocket | ✅ |
| FR-42 | Non-following users cannot message each other | ✅ |
| FR-43 | Messages only reach the intended recipient | ✅ |
| FR-44 | Group chat delivers messages to all online members | ✅ |
| FR-45 | Emojis can be sent in chat messages | ✅ |
| FR-46 | Typing indicators show when a user is typing | ✅ |

### ✅ Notifications
| ID | Requirement | Status |
|----|-------------|--------|
| FR-47 | Notifications are visible on every page | ✅ |
| FR-48 | Follow request notifications are sent to private users | ✅ |
| FR-49 | Group invitation notifications are sent to invited users | ✅ |
| FR-50 | Group join request notifications are sent to the creator | ✅ |
| FR-51 | Event creation notifications are sent to group members | ✅ |
| FR-52 | New notifications are distinct from new messages | ✅ |
| FR-53 | Users can mark notifications as read and delete them | ✅ |

### ✅ Docker
| ID | Requirement | Status |
|----|-------------|--------|
| FR-54 | Backend runs in a Docker container | ✅ |
| FR-55 | Frontend runs in a separate Docker container | ✅ |
| FR-56 | Both containers have non-zero sizes | ✅ |
| FR-57 | Application is accessible via web browser after Docker run | ✅ |

## Non-Functional Requirements

| ID | Requirement | Status |
|----|-------------|--------|
| NFR-01 | Passwords are stored using bcrypt hashing | ✅ |
| NFR-02 | Session tokens are stored in httpOnly cookies | ✅ |
| NFR-03 | Rate limiting is applied to API endpoints | ✅ |
| NFR-04 | Database migrations are applied on startup | ✅ |
| NFR-05 | Images are validated for type and size | ✅ |
| NFR-06 | Frontend is responsive across device sizes | ✅ |

## Allowed Packages

| Package | Status |
|---------|--------|
| Go standard library | ✅ |
| gorilla/websocket | ✅ |
| golang-migrate/migrate | ✅ |
| mattn/go-sqlite3 | ✅ |
| golang.org/x/crypto/bcrypt | ✅ |
| gofrs/uuid / google/uuid | ✅ |
| Next.js (JS framework) | ✅ |

## Bonus Features

| Feature | Status |
|---------|--------|
| GitHub OAuth login | ⚠️ Partial (backend code exists, routes may be disabled) |
| Database seeder | ✅ |
| Confirmation pop-up on unfollow | ✅ |
| Confirmation pop-up on privacy change | ✅ |
| Extra notifications beyond required | ✅ |
| Build script for Docker images | ✅ |
