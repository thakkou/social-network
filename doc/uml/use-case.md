# Use Case Diagram — Social Network

```mermaid
%%{init: {'theme':'dark', 'gitGraph': {'rotateCommitLabel': true}} }%%
usecaseDiagram
    title Social Network — Use Case Diagram

    actor "Guest\n(unauthenticated)" as Guest
    actor "User\n(authenticated)" as User
    actor "Group Admin" as Admin

    rectangle "Authentication" {
        usecase "UC1\nRegister" as UC1
        usecase "UC2\nLogin" as UC2
        usecase "UC3\nLogout" as UC3
    }

    rectangle "Profile" {
        usecase "UC4\nView Profile" as UC4
        usecase "UC5\nUpdate Profile" as UC5
        usecase "UC6\nToggle Privacy" as UC6
    }

    rectangle "Posts" {
        usecase "UC7\nCreate Post" as UC7
        usecase "UC8\nView Feed" as UC8
        usecase "UC9\nView Single Post" as UC9
        usecase "UC10\nDelete Post" as UC10
        usecase "UC11\nLike / Dislike Post" as UC11
    }

    rectangle "Comments" {
        usecase "UC12\nCreate Comment" as UC12
        usecase "UC13\nDelete Comment" as UC13
        usecase "UC14\nLike / Dislike Comment" as UC14
    }

    rectangle "Follow System" {
        usecase "UC15\nFollow User" as UC15
        usecase "UC16\nAccept / Reject Follow" as UC16
        usecase "UC17\nUnfollow User" as UC17
        usecase "UC18\nView Followers / Following" as UC18
    }

    rectangle "Groups" {
        usecase "UC19\nCreate Group" as UC19
        usecase "UC20\nList Groups" as UC20
        usecase "UC21\nJoin Group (request)" as UC21
        usecase "UC22\nLeave Group" as UC22
        usecase "UC23\nInvite Members" as UC23
        usecase "UC24\nAccept / Reject Invite" as UC24
        usecase "UC25\nApprove Join Requests" as UC25
        usecase "UC26\nKick Member" as UC26
        usecase "UC27\nPost in Group" as UC27
        usecase "UC28\nCreate Event" as UC28
        usecase "UC29\nRespond to Event" as UC29
        usecase "UC30\nGroup Chat" as UC30
    }

    rectangle "Messaging" {
        usecase "UC31\nSend Direct Message" as UC31
        usecase "UC32\nView Conversations" as UC32
        usecase "UC33\nTyping Indicator" as UC33
    }

    rectangle "Notifications" {
        usecase "UC34\nView Notifications" as UC34
        usecase "UC35\nMark as Read" as UC35
        usecase "UC36\nDelete Notification" as UC36
    }

    rectangle "Search" {
        usecase "UC37\nSearch Users" as UC37
        usecase "UC38\nSearch Groups" as UC38
    }

    Guest --> UC1
    Guest --> UC2

    User --> UC3
    User --> UC4
    User --> UC5
    User --> UC6
    User --> UC7
    User --> UC8
    User --> UC9
    User --> UC10
    User --> UC11
    User --> UC12
    User --> UC13
    User --> UC14
    User --> UC15
    User --> UC16
    User --> UC17
    User --> UC18
    User --> UC19
    User --> UC20
    User --> UC21
    User --> UC22
    User --> UC23
    User --> UC24
    User --> UC25
    User --> UC27
    User --> UC28
    User --> UC29
    User --> UC30
    User --> UC31
    User --> UC32
    User --> UC33
    User --> UC34
    User --> UC35
    User --> UC36
    User --> UC37
    User --> UC38

    Admin --> UC25
    Admin --> UC26
    Admin --> UC19
```
