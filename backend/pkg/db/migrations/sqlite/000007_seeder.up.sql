-- needs a down file

-- $2a$10$WKwmrHpemaSqsTWBATFNDO3VWX23TS9bEqA0FRkJVsYPS9OKRebV.  ==> 'password'

-- USERS
INSERT INTO USERS (id, firstname, lastname, email, password, birthdate, nickname, aboutme, avatar, is_private) VALUES
(1,'Alice','Johnson','alice@example.com','$2a$10$WKwmrHpemaSqsTWBATFNDO3VWX23TS9bEqA0FRkJVsYPS9OKRebV.','1997-04-12','alice','Go developer','/uploads/seeder/avatars/avatar.jpeg',0),
(2,'Bob','Smith','bob@example.com','$2a$10$WKwmrHpemaSqsTWBATFNDO3VWX23TS9bEqA0FRkJVsYPS9OKRebV.','1995-09-22','bob','Traveler','/uploads/seeder/avatars/goat.jpg',0),
(3,'Charlie','Brown','charlie@example.com','$2a$10$WKwmrHpemaSqsTWBATFNDO3VWX23TS9bEqA0FRkJVsYPS9OKRebV.','1999-01-30','charlie','Coffee lover','/uploads/seeder/avatars/lofi.jpeg',1),
(4,'Diana','Wilson','diana@example.com','$2a$10$WKwmrHpemaSqsTWBATFNDO3VWX23TS9bEqA0FRkJVsYPS9OKRebV.','1998-11-15','diana','Designer','/uploads/seeder/avatars/avatar.jpeg',0),
(5,'Ethan','Miller','ethan@example.com','$2a$10$WKwmrHpemaSqsTWBATFNDO3VWX23TS9bEqA0FRkJVsYPS9OKRebV.','1994-08-09','ethan','Backend engineer','/uploads/seeder/avatars/goat.jpg',1);

INSERT INTO FOLLOWS(follower_id,following_id,status) VALUES
(1,2,'accepted'),
(2,1,'accepted'),
(1,3,'pending'),
(4,1,'accepted'),
(5,2,'accepted');

INSERT INTO POSTS(id,user_id,created_at,title,text,image,privacy) VALUES
(1,1,datetime('now','-8 day'),'Learning Go','Pointers finally clicked.','/uploads/seeder/posts/golang-example.jpeg','public'),
(2,2,datetime('now','-7 day'),'Vacation','Amazing trip!','/uploads/seeder/posts/travel.jpeg','public'),
(3,3,datetime('now','-6 day'),'Cooking','Homemade pasta.','/uploads/seeder/posts/cooking1.jpeg','almost_private'),
(4,4,datetime('now','-5 day'),'UI Design','Minimalism matters.','/uploads/seeder/posts/web-design.jpeg','public'),
(5,5,datetime('now','-4 day'),'MCP','Testing the new API.','/uploads/seeder/posts/mcp-new-api.jpeg','private');

INSERT INTO POST_CATEGORY VALUES
(1,6),(2,4),(3,5),(4,9),(5,2);

INSERT INTO COMMENTS(user_id,post_id,created_at,text) VALUES
(2,1,datetime('now','-7 day'),'Nice work!'),
(3,1,datetime('now','-7 day'),'Very helpful.'),
(1,2,datetime('now','-6 day'),'Looks amazing!'),
(5,4,datetime('now','-3 day'),'Love the colors.');

INSERT INTO POST_REACTIONS VALUES
(2,1,1),
(3,1,1),
(4,1,1),
(1,2,1),
(5,2,-1),
(2,4,1);

INSERT INTO CONVERSATIONS(id,user1_id,user2_id,last_message,last_message_at) VALUES
(1,1,2,'See you tomorrow!',datetime('now'));

INSERT INTO MESSAGES(conversation_id,sender_id,text) VALUES
(1,1,'Hey Bob!'),
(1,2,'Hi Alice!'),
(1,1,'See you tomorrow!');

INSERT INTO GROUPS(id,logo,background,creator_id,title,description) VALUES
(1,'/uploads/seeder/groups/golage-log.jpeg','/uploads/seeder/groups/golang-bg.jpeg',1,'Go Developers','Everything Go'),
(2,'/uploads/seeder/groups/travel-logo.jpeg','/uploads/seeder/groups/travel-bg.jpeg',2,'Travel Club','Share adventures');

INSERT INTO GROUP_MEMBERS(group_id,user_id,role) VALUES
(1,1,'admin'),
(1,5,'member'),
(2,2,'admin'),
(2,1,'member');

INSERT INTO GROUP_POSTS(group_id,user_id,text,image) VALUES
(1,1,'Check out this package.','/uploads/seeder/posts/go-code.jpeg'),
(2,2,'Beautiful destination!','/uploads/seeder/posts/travel2.jpeg');

INSERT INTO GROUP_MESSAGES(group_id,sender_id,text) VALUES
(1,1,'Welcome everyone!'),
(1,5,'Happy to join!'),
(2,2,'Next trip?');

INSERT INTO GROUP_EVENTS(group_id,creator_id,title,description,event_time) VALUES
(1,1,'Go Meetup','Monthly meetup',datetime('now','+7 day')),
(2,2,'Weekend Hike','Mountain trail',datetime('now','+14 day'));

INSERT INTO EVENT_RESPONSES VALUES
(1,5,'going'),
(2,1,'going');

INSERT INTO NOTIFICATIONS(user_id,actor_id,type,object_type,object_id) VALUES
(1,2,'like','post',1),
(2,1,'follow','user',1),
(5,1,'group_invite','group',1);