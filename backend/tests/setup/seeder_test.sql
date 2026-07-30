-- SEEDER MIGRATION

INSERT INTO USERS (id, firstname, lastname, email, password, birthdate, nickname, aboutme, avatar, is_private) VALUES
(1,  'Alice',   'Martin',     'alice@example.com',    '$2a$10$WKwmrHpemaSqsTWBATFNDO3VWX23TS9bEqA0FRkJVsYPS9OKRebV.', '1996-04-12', 'ali_m',   'Coffee & code.',                     '/uploads/seeder/avatars/avatar.jpeg', 0),
(2,  'Bob',     'Nguyen',     'bob@example.com',      '$2a$10$WKwmrHpemaSqsTWBATFNDO3VWX23TS9bEqA0FRkJVsYPS9OKRebV.', '1994-08-23', NULL,      'Traveling the world.',               '/uploads/seeder/avatars/lofi.jpeg',   0),
(3,  'Chloe',   'Dubois',     'chloe@example.com',    '$2a$10$WKwmrHpemaSqsTWBATFNDO3VWX23TS9bEqA0FRkJVsYPS9OKRebV.', '1999-01-05', 'Chloe',   NULL,                                 '/uploads/seeder/avatars/goat.jpg',    1),
(4,  'David',   'Smith',      'david@example.com',    '$2a$10$WKwmrHpemaSqsTWBATFNDO3VWX23TS9bEqA0FRkJVsYPS9OKRebV.', '1990-11-30', 'david',   'Full-stack dev.',                    NULL,                                   0),
(5,  'Emma',    'Wilson',     'emma@example.com',     '$2a$10$WKwmrHpemaSqsTWBATFNDO3VWX23TS9bEqA0FRkJVsYPS9OKRebV.', '1997-06-18', NULL,      'Photography enthusiast.',             NULL,                                   1),
(6,  'Farid',   'El Amrani',  'farid@example.com',    '$2a$10$WKwmrHpemaSqsTWBATFNDO3VWX23TS9bEqA0FRkJVsYPS9OKRebV.', '1993-03-09', NULL,      'Backend > frontend, fight me.',       '/uploads/seeder/avatars/lofi.jpeg',   0),
(7,  'Grace',   'Lee',        'grace@example.com',    '$2a$10$WKwmrHpemaSqsTWBATFNDO3VWX23TS9bEqA0FRkJVsYPS9OKRebV.', '2000-09-27', NULL,      NULL,                                 NULL,                                   0),
(8,  'Hugo',    'Costa',      'hugo@example.com',     '$2a$10$WKwmrHpemaSqsTWBATFNDO3VWX23TS9bEqA0FRkJVsYPS9OKRebV.', '1995-12-14', 'costa77', 'Music producer.',                    NULL,                                   0),
(9,  'Isabella','Rossi',      'isabella@example.com', '$2a$10$WKwmrHpemaSqsTWBATFNDO3VWX23TS9bEqA0FRkJVsYPS9OKRebV.', '1998-07-22', 'bella',   'Art & design lover.',                NULL,                                   0),
(10, 'Jack',    'Thompson',   'jack@example.com',     '$2a$10$WKwmrHpemaSqsTWBATFNDO3VWX23TS9bEqA0FRkJVsYPS9OKRebV.', '1992-10-05', NULL,      'Open source contributor.',            NULL,                                   0);

-- ============================================================
-- FOLLOWS (14 follows)
-- ============================================================
INSERT INTO FOLLOWS (follower_id, following_id, status) VALUES
(1, 2,  'accepted'),
(2, 1,  'accepted'),
(1, 4,  'accepted'),
(4, 6,  'accepted'),
(6, 1,  'accepted'),
(7, 3,  'pending'),
(8, 5,  'pending'),
(2, 5,  'accepted'),
(3, 1,  'accepted'),
(1, 7,  'accepted'),
(2, 4,  'accepted'),
(8, 1,  'accepted'),
(4, 1,  'accepted'),
(6, 8,  'accepted');

-- ============================================================
-- POSTS (63 posts)
-- ============================================================
INSERT INTO POSTS (id, user_id, created_at, title, text, image, privacy) VALUES
-- Original posts
(1,  1,  datetime('now', '-63 hours'),  'Hello world',                                     'My very first post on this network!',                                             '/uploads/seeder/posts/dev.jpeg',      'public'),
(2,  2,  datetime('now', '-62 hours'),  'Weekend trip',                                    'Just got back from the mountains 🏔️',                                              '/uploads/seeder/posts/travel.jpeg',   'public'),
(3,  4,  datetime('now', '-61 hours'),  NULL,                                              'Debugging is 90% of my job today.',                                               '/uploads/seeder/posts/dev2.jpeg',     'almost_private'),
(4,  6,  datetime('now', '-60 hours'),  'Go tip',                                          'context.Context should be your first param, always.',                              '/uploads/seeder/posts/dev3.jpeg',     'public'),
(5,  2,  datetime('now', '-59 hours'),  'Private thoughts',                                'Only a few people should see this.',                                               '/uploads/seeder/posts/learning.jpeg', 'private'),
(6,  2,  datetime('now', '-58 hours'),  'Morning routine',                                 'Coffee, reading, then coding.',                                                    '/uploads/seeder/posts/cooking1.jpeg', 'public'),
(7,  2,  datetime('now', '-57 hours'),  'Photography',                                     'Took some beautiful sunset photos today.',                                         '/uploads/seeder/posts/travel2.jpeg',  'public'),
(8,  2,  datetime('now', '-56 hours'),  'Weekend plans',                                   'Thinking about hiking this weekend.',                                              '/uploads/seeder/posts/cooking2.jpeg', 'almost_private'),
(9,  3,  datetime('now', '-55 hours'),  'Learning Go',                                     'Interfaces finally clicked today!',                                                '/uploads/seeder/posts/learning2.jpeg', 'public'),
(10, 5,  datetime('now', '-54 hours'),  'New camera',                                      'Testing out my new lens today.',                                                   '/uploads/seeder/posts/travel.jpeg',   'almost_private'),
(11, 8,  datetime('now', '-53 hours'),  'New track',                                       'Dropping a new beat this Friday 🎵',                                                NULL,                                  'public'),
(12, 1,  datetime('now', '-52 hours'),  'Distributed Systems',                             'Reading ''Designing Data-Intensive Applications''. Highly recommend!',               '/uploads/seeder/posts/dev3.jpeg',     'public'),
(13, 6,  datetime('now', '-51 hours'),  'Rust vs Go',                                      'Trying out Rust for a CLI tool. The borrow checker is something else.',             '/uploads/seeder/posts/dev2.jpeg',     'public'),
(14, 4,  datetime('now', '-50 hours'),  'Weekend project',                                 'Built a tiny load balancer in Go over the weekend.',                                NULL,                                  'almost_private'),
(15, 7,  datetime('now', '-49 hours'),  'New hobby',                                       'Started learning the piano! Any tips for beginners?',                              NULL,                                  'public'),
(16, 8,  datetime('now', '-48 hours'),  'Studio update',                                   'New soundproof panels in the studio. Check it out!',                               '/uploads/seeder/posts/dev.jpeg',      'public'),
(17, 3,  datetime('now', '-47 hours'),  'Book club',                                       'Reading ''The Pragmatic Programmer'' with the Go group. Join us!',                    NULL,                                  'public'),
(18, 5,  datetime('now', '-46 hours'),  'Sunset shots',                                    'Golden hour at the lake today. Nature is healing.',                                '/uploads/seeder/posts/travel2.jpeg',  'public'),
(19, 9,  datetime('now', '-45 hours'),  'Art Exhibition',                                  'Visited the modern art museum today. Absolutely inspiring pieces!',                '/uploads/seeder/posts/travel.jpeg',   'public'),
(20, 10, datetime('now', '-44 hours'),  'Open Source Saturday',                            'Contributed to my first open source project today! Feeling great.',                '/uploads/seeder/posts/dev2.jpeg',     'public'),
(21, 7,  datetime('now', '-43 hours'),  'Baking bread',                                    'Made sourdough from scratch for the first time. Turned out amazing!',              '/uploads/seeder/posts/cooking1.jpeg', 'public'),
(22, 1,  datetime('now', '-42 hours'),  'Microservices talk',                              'Went to a great meetup about microservices patterns. Event sourcing is fascinating.','/uploads/seeder/posts/dev3.jpeg',     'public'),
(23, 3,  datetime('now', '-41 hours'),  'Coding playlist',                                 'Share your favorite coding playlist! I need new recommendations.',                 NULL,                                  'public'),
(24, 6,  datetime('now', '-40 hours'),  'Docker tips',                                     'Multi-stage builds are a game changer for reducing image size. Here''s how...',     '/uploads/seeder/posts/dev.jpeg',      'public'),
(25, 4,  datetime('now', '-39 hours'),  'Night sky',                                       'Captured the Milky Way with my telescope last night. Astro-photography is hard!',  '/uploads/seeder/posts/travel2.jpeg',  'almost_private'),
(26, 8,  datetime('now', '-38 hours'),  'New single out now',                              'My new single ''Midnight Code'' is out on all platforms! Link in bio 🎶',            NULL,                                  'public'),
(27, 9,  datetime('now', '-37 hours'),  'Digital art',                                     'Experimenting with procedural art generation using Go. Unexpectedly beautiful!',   '/uploads/seeder/posts/learning.jpeg', 'public'),
(28, 2,  datetime('now', '-36 hours'),  'Travel tips',                                     'Best travel destinations for digital nomads in 2026: full guide in thread.',       '/uploads/seeder/posts/travel.jpeg',   'public'),
(29, 5,  datetime('now', '-35 hours'),  'Film photography',                                'Got my first film roll developed. There''s something magical about analog.',        '/uploads/seeder/posts/learning2.jpeg','public'),
(30, 10, datetime('now', '-34 hours'),  'Terminal tools',                                  'My favorite terminal tools: fzf, ripgrep, bat, and lazygit. What are yours?',      '/uploads/seeder/posts/dev3.jpeg',     'public'),
(31, 7,  datetime('now', '-33 hours'),  'Hiking adventure',                                'Summited Mount Toubkal! Toughest hike of my life but the view was worth it.',      '/uploads/seeder/posts/travel2.jpeg',  'public'),
(32, 1,  datetime('now', '-32 hours'),  'AI pair programming',                             'Been using AI tools for code reviews. They catch things I miss!',                  '/uploads/seeder/posts/dev2.jpeg',     'public'),
(33, 4,  datetime('now', '-31 hours'),  'Minimalist setup',                                'My new WFH setup: minimal, clean, productive. Less is more.',                      '/uploads/seeder/posts/dev.jpeg',      'public'),
(34, 3,  datetime('now', '-30 hours'),  'Pointers in C',                                   'Finally understanding pointers in C after years of avoiding them. It''s just references!','/uploads/seeder/posts/learning2.jpeg','public'),
(35, 7,  datetime('now', '-29 hours'),  'Morning yoga',                                    '30 days of morning yoga done. Flexibility improved, stress gone. Game changer.',   '/uploads/seeder/posts/lofi.jpeg',     'public'),
(36, 9,  datetime('now', '-28 hours'),  'New painting',                                    'Finally finished my oil painting after 3 weeks. Titled ''Digital Dreams''.',         NULL,                                  'public'),
(37, 10, datetime('now', '-27 hours'),  'Rust CLI tool',                                   'Built a CLI file organizer in Rust this weekend. 10x faster than my Python version.','/uploads/seeder/posts/dev2.jpeg',     'public'),
(38, 5,  datetime('now', '-26 hours'),  'Street photography',                              'Black and white street photography is my new obsession. Capturing raw moments.',    '/uploads/seeder/posts/travel.jpeg',   'public'),
(39, 1,  datetime('now', '-25 hours'),  'Promotion news',                                  'Getting promoted to Senior Engineer next month! All those late nights paid off.',  NULL,                                  'public'),
(40, 8,  datetime('now', '-24 hours'),  'New EP done',                                     'My debut EP ''Electric Dreams'' is finally mastered and ready. Dropping next Friday!','/uploads/seeder/posts/dev.jpeg',      'public'),
(41, 4,  datetime('now', '-23 hours'),  'Homemade sushi',                                  'Made sushi from scratch for the first time. Rice to fish ratio is an art.',         '/uploads/seeder/posts/cooking1.jpeg', 'public'),
(42, 2,  datetime('now', '-22 hours'),  'Sahara expedition',                               '3 days trekking the Sahara desert. Sand dunes, starry nights, and pure silence.',  '/uploads/seeder/posts/travel2.jpeg',  'public'),
(43, 6,  datetime('now', '-21 hours'),  'Kubernetes thoughts',                             'Kubernetes is both the best and most frustrating thing I''ve ever worked with.',    '/uploads/seeder/posts/dev3.jpeg',     'public'),
(44, 10, datetime('now', '-20 hours'),  'PR merged!',                                      'My first open source PR got merged into a major project! Contributing feels amazing.','/uploads/seeder/posts/dev2.jpeg',     'public'),
(45, 7,  datetime('now', '-19 hours'),  'First marathon',                                  'Ran my first marathon! 4h 23m. Never thought I could do it but here we are.',      '/uploads/seeder/posts/travel2.jpeg',  'public'),
(46, 3,  datetime('now', '-18 hours'),  'Clean Code review',                               'Just finished ''Clean Code'' by Uncle Bob. Every developer should read this.',        '/uploads/seeder/posts/learning.jpeg', 'public'),
(47, 9,  datetime('now', '-17 hours'),  'UI design tips',                                  '5 UI/UX tips for beginners: whitespace is your friend, consistency matters, accessibility first.','/uploads/seeder/posts/dev.jpeg','public'),
(48, 1,  datetime('now', '-16 hours'),  'Legacy code',                                     'Refactoring a 10-year-old codebase. It''s like archeology but with more existential dread.','/uploads/seeder/posts/dev3.jpeg','public'),
(49, 8,  datetime('now', '-15 hours'),  'Studio tour',                                     'My home recording studio setup: Focusrite Scarlett, SM7B, and lots of patience.',   '/uploads/seeder/posts/cooking2.jpeg', 'public'),
(50, 5,  datetime('now', '-14 hours'),  'Tokyo guide',                                     'Just got back from Tokyo! Here''s my guide: Shibuya at night, Tsukiji for breakfast, Akihabara for tech.','/uploads/seeder/posts/travel.jpeg','public'),
(51, 4,  datetime('now', '-13 hours'),  'Weather app',                                     'Building a weather app with Go and HTMX. No JavaScript, minimal CSS, maximum fun.','/uploads/seeder/posts/learning2.jpeg', 'public'),
(52, 2,  datetime('now', '-12 hours'),  'Alps hiking',                                     'Best hiking trails in the Swiss Alps: Eiger Trail, Haute Route, and Jungfrau region.','/uploads/seeder/posts/travel2.jpeg','public'),
(53, 6,  datetime('now', '-11 hours'),  'gRPC vs REST',                                    'After building APIs with both: gRPC for internal services, REST for public APIs. Different tools.','/uploads/seeder/posts/dev2.jpeg','public'),
(54, 7,  datetime('now', '-10 hours'),  'Vegan banana bread',                              'Vegan banana bread recipe that even non-vegans love. Secret ingredient: coconut oil.','/uploads/seeder/posts/cooking1.jpeg','public'),
(55, 10, datetime('now', '-9 hours'),   'Custom keyboard',                                 'Built my first mechanical keyboard! GMK keycaps, Gateron Black switches, aluminum case.','/uploads/seeder/posts/lofi.jpeg','public'),
(56, 3,  datetime('now', '-8 hours'),   'Linear algebra',                                  'Linear algebra is like the physics of programming. Everything makes sense with matrices.','/uploads/seeder/posts/learning2.jpeg','public'),
(57, 9,  datetime('now', '-7 hours'),   'Forest photography',                              'Spent the weekend in the forest with my camera. Nothing beats natural light through leaves.','/uploads/seeder/posts/travel.jpeg','public'),
(58, 1,  datetime('now', '-6 hours'),   'SQLite vs PostgreSQL',                            'For side projects: SQLite for simplicity, PostgreSQL when you need features. Both are amazing.','/uploads/seeder/posts/dev3.jpeg','public'),
(59, 8,  datetime('now', '-5 hours'),   'Concert night',                                   'Went to see a jazz fusion band last night. Live music hits different. 🎷',           NULL,                                  'public'),
(60, 5,  datetime('now', '-4 hours'),   'Film photography',                                'First roll of Fujifilm Superia 400 developed. There''s magic in the imperfections.','/uploads/seeder/posts/learning.jpeg', 'public'),
(61, 4,  datetime('now', '-3 hours'),   'Home server',                                     'My home server setup: Raspberry Pi 5, 4TB SSD, Pi-hole, Jellyfin, and Grafana.',    '/uploads/seeder/posts/dev.jpeg',      'public'),
(62, 2,  datetime('now', '-2 hours'),   'Budget travel',                                   'How to travel Europe on €50/day: hostels over hotels, street food, free walking tours.','/uploads/seeder/posts/travel2.jpeg','public'),
(63, 6,  datetime('now', '-1 hour'),    '.vimrc secrets',                                  'My .vimrc settings after 5 years of Vim: relative numbers, easy motion, and snippets.','/uploads/seeder/posts/dev2.jpeg',    'public');

-- ============================================================
-- POST ALLOWED USERS (for private post #5)
-- ============================================================
INSERT INTO POST_ALLOWED_USERS (post_id, user_id) VALUES
(5, 1),
(5, 7);

-- ============================================================
-- POST CATEGORIES
-- Category IDs: 1=General, 2=Lifestyle, 3=Health&Fitness, 4=Travel,
--               5=Food&Cooking, 6=Education, 7=Business, 8=Finance,
--               9=Entertainment, 10=Sports, 11=Personal Dev, 12=Culture
-- Matches seedPostCategories() in posts.go
-- ============================================================
INSERT INTO POST_CATEGORY (post_id, category_id) VALUES
-- Original post categories (kept)
(1,  1),
(2,  4),
(4,  6),
(7,  9),
(12, 6),
(13, 6),
(19, 9),
(21, 1),
(22, 6),
(24, 6),
(25, 6),
(26, 9),
(27, 4),
(28, 4),    (28, 10),
(29, 9),
-- Categories for previously uncategorized posts
(3,  6),
(6,  2),
(8,  10),   (8,  3),
(9,  6),
(10, 9),
(11, 9),
(14, 6),
(15, 11),
(16, 9),
(17, 6),    (17, 12),
(18, 4),
(20, 1),    (20, 6),
(23, 9),
(30, 6),
(31, 4),    (31, 10),
(32, 6),
(33, 2),
-- Categories for 30 new posts (indices 33-62)
(34, 6),
(35, 3),
(36, 12),
(37, 6),
(38, 9),
(39, 7),
(40, 9),
(41, 5),
(42, 4),
(43, 6),
(44, 1),
(45, 10),   (45, 3),
(46, 11),
(47, 6),
(48, 6),
(49, 9),
(50, 4),
(51, 6),
(52, 10),   (52, 4),
(53, 6),
(54, 5),
(55, 2),
(56, 6),
(57, 4),    (57, 12),
(58, 6),
(59, 9),    (59, 12),
(60, 12),
(61, 6),
(62, 4),    (62, 8),
(63, 6);

-- ============================================================
-- COMMENTS (88 comments)
-- ============================================================
INSERT INTO COMMENTS (id, user_id, post_id, created_at, text) VALUES
(1,  2,  1,  datetime('now', '-2640 minutes'), 'Welcome! Great to have you here.'),
(2,  4,  1,  datetime('now', '-2610 minutes'), 'Nice first post 🎉'),
(3,  1,  2,  datetime('now', '-2580 minutes'), 'Looks amazing, take me next time!'),
(4,  6,  4,  datetime('now', '-2550 minutes'), 'Solid tip, saved me a bug last week.'),
(5,  1,  5,  datetime('now', '-2520 minutes'), 'Thanks for sharing this with me.'),
(6,  2,  11, datetime('now', '-2490 minutes'), 'Can''t wait to hear it!'),
(7,  1,  6,  datetime('now', '-2460 minutes'), 'That sounds like a productive morning!'),
(8,  2,  6,  datetime('now', '-2430 minutes'), 'Coffee first is always the right choice ☕'),
(9,  4,  7,  datetime('now', '-2400 minutes'), 'Would love to see those photos!'),
(10, 6,  7,  datetime('now', '-2370 minutes'), 'Sunsets are the best 🌅'),
(11, 1,  8,  datetime('now', '-2340 minutes'), 'Hope the weather stays nice!'),
(12, 8,  8,  datetime('now', '-2310 minutes'), 'Enjoy your hike!'),
(13, 2,  9,  datetime('now', '-2280 minutes'), 'Go interfaces are confusing at first 😄'),
(14, 6,  9,  datetime('now', '-2250 minutes'), 'Wait until you discover generics!'),
(15, 2,  12, datetime('now', '-2220 minutes'), 'That book is a masterpiece! Chapter 5 on replication is gold.'),
(16, 6,  12, datetime('now', '-2190 minutes'), 'I should re-read this. Great recommendation.'),
(17, 4,  13, datetime('now', '-2160 minutes'), 'Borrow checker is tough but rewarding.'),
(18, 1,  13, datetime('now', '-2130 minutes'), 'Stick with it! Rust''s ecosystem is growing fast.'),
(19, 6,  15, datetime('now', '-2100 minutes'), 'Consistency is key! 10 minutes a day makes a big difference.'),
(20, 1,  16, datetime('now', '-2070 minutes'), 'Nice setup! What DAW are you using?'),
(21, 2,  17, datetime('now', '-2040 minutes'), 'Count me in! I''m halfway through the book.'),
(22, 7,  18, datetime('now', '-2010 minutes'), 'Absolutely stunning 😍'),
(23, 1,  19, datetime('now', '-1980 minutes'), 'Modern art is so thought-provoking. Which exhibit was your favorite?'),
(24, 6,  19, datetime('now', '-1950 minutes'), 'Wish I could visit! Send pictures.'),
(25, 2,  20, datetime('now', '-1920 minutes'), 'That''s awesome! Which project did you contribute to?'),
(26, 8,  20, datetime('now', '-1890 minutes'), 'Open source is the way to go! Keep it up 💪'),
(27, 4,  21, datetime('now', '-1860 minutes'), 'That looks delicious! Drop the recipe 🍞'),
(28, 9,  21, datetime('now', '-1830 minutes'), 'Sourdough is the best. Try adding rosemary next time!'),
(29, 6,  22, datetime('now', '-1800 minutes'), 'Event sourcing is great for audit logs. Check out EventStore!'),
(30, 2,  22, datetime('now', '-1770 minutes'), 'I was at that meetup too! Great talk.'),
(31, 10, 23, datetime('now', '-1740 minutes'), 'Lo-fi hip hop beats to code to, always.'),
(32, 1,  23, datetime('now', '-1710 minutes'), 'I like ambient electronic. Brian Eno is great.'),
(33, 4,  24, datetime('now', '-1680 minutes'), 'Multi-stage builds saved me 60% image size. Great tip!'),
(34, 8,  24, datetime('now', '-1650 minutes'), 'Don''t forget to use .dockerignore too!'),
(35, 1,  25, datetime('now', '-1620 minutes'), 'Absolutely breathtaking! What equipment did you use?'),
(36, 7,  25, datetime('now', '-1590 minutes'), 'I tried astrophotography once, it''s so difficult. Great shot!'),
(37, 2,  26, datetime('now', '-1560 minutes'), 'Congrats on the release! Adding to my playlist now 🎧'),
(38, 6,  26, datetime('now', '-1530 minutes'), 'Love the title! ''Midnight Code'' hits different.'),
(39, 4,  27, datetime('now', '-1500 minutes'), 'Procedural generation is so cool. What library did you use?'),
(40, 10, 28, datetime('now', '-1470 minutes'), 'Portugal is amazing for digital nomads. Highly recommend!'),
(41, 1,  28, datetime('now', '-1440 minutes'), 'Thailand is also incredible. Great community there.'),
(42, 2,  29, datetime('now', '-1410 minutes'), 'Analog has such a unique feel. What camera are you using?'),
(43, 8,  30, datetime('now', '-1380 minutes'), 'ripgrep is a lifesaver! I also love jq for JSON processing.'),
(44, 4,  30, datetime('now', '-1350 minutes'), 'lazygit changed my workflow completely.'),
(45, 6,  31, datetime('now', '-1320 minutes'), 'Wow that''s impressive! How long did the hike take?'),
(46, 9,  31, datetime('now', '-1290 minutes'), 'Toubkal is on my bucket list! Any tips for beginners?'),
(47, 2,  32, datetime('now', '-1260 minutes'), 'AI code review is surprisingly good. Which tool do you use?'),
(48, 7,  32, datetime('now', '-1230 minutes'), 'Just be careful about sensitive data in the prompts!'),
(49, 1,  33, datetime('now', '-1200 minutes'), 'Love the clean setup! What monitor is that?'),
(50, 6,  33, datetime('now', '-1170 minutes'), 'Minimalist setups are the best for focus.'),
(51, 1,  34, datetime('now', '-1140 minutes'), 'Pointers are tricky but so powerful once they click!'),
(52, 6,  34, datetime('now', '-1110 minutes'), 'Wait until you try double pointers 😄'),
(53, 9,  35, datetime('now', '-1080 minutes'), 'Yoga has been life-changing for my posture too!'),
(54, 2,  36, datetime('now', '-1050 minutes'), 'Would love to see it! Post a photo!'),
(55, 4,  37, datetime('now', '-1020 minutes'), 'Rust is amazing for CLIs. Did you use clap for argument parsing?'),
(56, 1,  38, datetime('now', '-990 minutes'),  'Street photography is an art. Henri Cartier-Bresson is my inspiration.'),
(57, 6,  39, datetime('now', '-960 minutes'),  'Congrats on the promotion! Well deserved! 🎉'),
(58, 2,  39, datetime('now', '-930 minutes'),  'Senior Engineer! Amazing achievement.'),
(59, 4,  40, datetime('now', '-900 minutes'),  'Can''t wait to hear it! What genre?'),
(60, 6,  41, datetime('now', '-870 minutes'),  'Sushi is all about the rice. Great job!'),
(61, 7,  42, datetime('now', '-840 minutes'),  'The Sahara at night must be incredible. Did you see the stars?'),
(62, 1,  42, datetime('now', '-810 minutes'),  'Desert camping is on my bucket list!'),
(63, 4,  43, datetime('now', '-780 minutes'),  'Kubernetes in a nutshell: works great until it doesn''t 😂'),
(64, 8,  43, datetime('now', '-750 minutes'),  'K8s is overkill for small projects but essential at scale.'),
(65, 1,  44, datetime('now', '-720 minutes'),  'Congrats on the PR! Which project?'),
(66, 6,  45, datetime('now', '-690 minutes'),  '4h 23m is impressive for a first marathon!'),
(67, 1,  46, datetime('now', '-660 minutes'),  'Clean Code changed how I think about naming variables.'),
(68, 6,  47, datetime('now', '-630 minutes'),  'Accessibility is still underrated in web design. Great tips!'),
(69, 4,  48, datetime('now', '-600 minutes'),  'Legacy code: where ''works on my machine'' was written in 2012.'),
(70, 9,  48, datetime('now', '-570 minutes'),  'Good luck! I''ve been there. Incremental refactoring is key.'),
(71, 2,  49, datetime('now', '-540 minutes'),  'SM7B is a classic choice! Great for vocals.'),
(72, 6,  50, datetime('now', '-510 minutes'),  'Tokyo is my favorite city! Did you visit the Ghibli Museum?'),
(73, 8,  50, datetime('now', '-480 minutes'),  'Shibuya crossing at night is magical.'),
(74, 1,  51, datetime('now', '-450 minutes'),  'HTMX + Go is such a productive stack!'),
(75, 6,  52, datetime('now', '-420 minutes'),  'The Eiger Trail is stunning! Did you see the Eiger north face?'),
(76, 9,  53, datetime('now', '-390 minutes'),  'Great summary! I''d add that gRPC streaming is a game changer.'),
(77, 2,  54, datetime('now', '-360 minutes'),  'Coconut oil in baking is genius! Trying this weekend.'),
(78, 4,  55, datetime('now', '-330 minutes'),  'Gateron Blacks are smooth! What case did you get?'),
(79, 7,  56, datetime('now', '-300 minutes'),  'Linear algebra finally clicked for me with 3Blue1Brown''s series.'),
(80, 2,  57, datetime('now', '-270 minutes'),  'Forest photography in autumn is unbeatable.'),
(81, 6,  58, datetime('now', '-240 minutes'),  'SQLite for MVPs, PostgreSQL for production. Perfect combo.'),
(82, 10, 58, datetime('now', '-210 minutes'),  'Don''t forget about SQLite''s WAL mode for concurrent reads!'),
(83, 1,  59, datetime('now', '-180 minutes'),  'Jazz fusion is incredible! Which band was it?'),
(84, 6,  60, datetime('now', '-150 minutes'),  'Film has such a unique character. Digital can''t replicate it.'),
(85, 8,  61, datetime('now', '-120 minutes'),  'Raspberry Pi 5 is a beast for a home server!'),
(86, 1,  62, datetime('now', '-90 minutes'),   'Great tips! I''d add: book hostels with kitchen access.'),
(87, 4,  63, datetime('now', '-60 minutes'),   '.vimrc is a personal diary of a developer''s journey.'),
(88, 10, 63, datetime('now', '-30 minutes'),   'Easy motion is a must-have! Nice config.');

-- ============================================================
-- POST REACTIONS
-- ============================================================
INSERT INTO POST_REACTIONS (user_id, post_id, is_like) VALUES
(2,  1,  1),  (4,  1,  1),  (6,  1,  1),
(6,  2,  1),
(1,  4,  1),  (8,  4,  -1),
(1,  7,  1),
(2,  8,  1),
(2,  9,  1),
(1,  11, 1),
(4,  12, 1),
(2,  13, -1),
(8,  15, 1),
(3,  17, 1),
(1,  19, 1),  (6,  19, 1),
(2,  20, 1),  (8,  20, 1),
(4,  21, 1),  (9,  21, -1),
(6,  22, 1),  (2,  22, 1),
(1,  23, 1),  (10, 23, 1),
(4,  24, 1),  (8,  24, 1),
(1,  25, 1),  (7,  25, 1),
(2,  26, 1),  (6,  26, 1),
(4,  27, 1),
(1,  28, 1),
(2,  29, 1),
(8,  30, 1),
(6,  31, 1),  (9,  31, 1),
(2,  32, 1),  (7,  32, -1),
(1,  33, 1),  (6,  33, 1),
(3,  34, 1),  (6,  34, 1),
(7,  35, 1),  (9,  35, 1),
(9,  36, 1),  (2,  36, 1),
(10, 37, 1),  (4,  37, 1),
(5,  38, 1),  (1,  38, 1),
(1,  39, 1),  (2,  39, 1),  (6,  39, 1),
(8,  40, 1),  (4,  40, 1),
(4,  41, 1),  (6,  41, 1),
(2,  42, 1),  (7,  42, 1),  (1,  42, 1),
(6,  43, 1),  (4,  43, -1),
(10, 44, 1),  (1,  44, 1),
(7,  45, 1),  (6,  45, 1),
(3,  46, 1),  (1,  46, 1),
(9,  47, 1),  (6,  47, 1),
(1,  48, 1),  (9,  48, 1),
(8,  49, 1),  (2,  49, 1),
(5,  50, 1),  (6,  50, 1),  (8,  50, 1),
(4,  51, 1),  (1,  51, 1),
(2,  52, 1),  (6,  52, 1),
(6,  53, 1),  (9,  53, 1),
(7,  54, 1),  (2,  54, 1),
(10, 55, 1),  (4,  55, 1),
(3,  56, 1),  (7,  56, 1),
(9,  57, 1),  (2,  57, 1),
(1,  58, 1),  (6,  58, 1),  (10, 58, 1),
(8,  59, 1),  (1,  59, 1),
(5,  60, 1),  (6,  60, 1),
(4,  61, 1),  (8,  61, 1),
(2,  62, 1),  (1,  62, 1),
(6,  63, 1),  (4,  63, 1),  (10, 63, 1);

-- ============================================================
-- COMMENT REACTIONS
-- ============================================================
INSERT INTO COMMENT_REACTIONS (user_id, comment_id, is_like) VALUES
(1,  1,  1),
(6,  4,  1),
(2,  5,  1),
(1,  6,  1),
(4,  7,  1),
(6,  9,  -1),
(2,  11, 1),
(1,  13, 1),
(8,  15, 1),
(2,  16, 1),
(6,  19, 1),
(1,  21, 1),
(4,  23, 1),
(6,  25, 1),
(8,  27, 1),
(2,  29, 1),
(1,  31, 1),
(6,  33, 1),
(1,  52, 1),
(6,  54, 1),
(4,  56, 1),
(1,  58, 1),
(6,  60, 1),
(7,  62, 1),
(4,  64, 1),
(6,  67, 1),
(1,  70, 1),
(6,  73, 1),
(2,  76, 1),
(9,  79, 1),
(6,  82, 1),
(1,  85, 1),
(4,  87, 1);

-- ============================================================
-- CONVERSATIONS (2 conversations)
-- ============================================================
INSERT INTO CONVERSATIONS (id, user1_id, user2_id, last_message, last_message_at) VALUES
(1, 1, 2, 'See you tomorrow!',   datetime('now', '-2 hours')),
(2, 1, 6, 'Sounds good, thanks!', datetime('now', '-45 minutes'));

-- ============================================================
-- MESSAGES (4 messages)
-- ============================================================
INSERT INTO MESSAGES (conversation_id, sender_id, text, is_read) VALUES
(1, 1, 'Hey Bob, are we still on for tomorrow?', 1),
(1, 2, 'Yep, see you tomorrow!', 0),
(2, 1, 'Can you send me the trail map?', 1),
(2, 6, 'Sounds good, thanks!', 0);

-- ============================================================
-- GROUPS (4 groups)
-- ============================================================
INSERT INTO GROUPS (id, logo, background, creator_id, title, description) VALUES
(1, '/uploads/seeder/groups/golage-log.jpeg', '/uploads/seeder/groups/golang-bg.jpeg',   1, 'Gophers United',  'Everything about Go programming, backend development, concurrency and open source.'),
(2, '/uploads/seeder/groups/sport-logo.jpeg', '/uploads/seeder/groups/sport-bg.jpeg',     6, 'Sports Club',     'Football, basketball, running, fitness and every kind of sport.'),
(3, '/uploads/seeder/groups/cultur-log.jpeg', '/uploads/seeder/groups/culture-bg.jpeg',   8, 'Culture & Arts',  'Books, music, cinema, painting and cultural events.'),
(4, '/uploads/seeder/groups/travel-logo.jpeg','/uploads/seeder/groups/travel-bg.jpeg',    2, 'Travel Explorers','Share destinations, travel tips, hiking adventures and unforgettable experiences.');

-- ============================================================
-- GROUP MEMBERS
-- ============================================================
INSERT INTO GROUP_MEMBERS (group_id, user_id, role) VALUES
(1, 1, 'admin'),  (1, 2, 'member'), (1, 4, 'member'), (1, 6, 'member'),
(2, 6, 'admin'),  (2, 1, 'member'), (2, 7, 'member'), (2, 8, 'member'),
(3, 8, 'admin'),  (3, 3, 'member'), (3, 5, 'member'), (3, 2, 'member'),
(4, 2, 'admin'),  (4, 1, 'member'), (4, 5, 'member'), (4, 7, 'member');

-- ============================================================
-- GROUP INVITES
-- ============================================================
INSERT INTO GROUP_INVITES (group_id, inviter_id, invited_user_id, status) VALUES
(1, 1, 7, 'pending'),
(2, 6, 5, 'pending'),
(3, 8, 4, 'pending'),
(4, 2, 3, 'pending');

-- ============================================================
-- GROUP REQUESTS
-- ============================================================
INSERT INTO GROUP_REQUESTS (group_id, user_id, status) VALUES
(1, 8, 'pending'),
(2, 4, 'pending');

-- ============================================================
-- GROUP MESSAGES
-- ============================================================
INSERT INTO GROUP_MESSAGES (group_id, sender_id, text) VALUES
(1, 1, 'Welcome to Gophers United!'),
(1, 6, 'Anyone using Go 1.25?'),
(1, 4, 'Concurrency is amazing.'),
(2, 6, 'Football match this Friday?'),
(2, 8, 'I''m in! ⚽'),
(2, 1, 'See you at 7 PM.'),
(3, 8, 'Movie night this weekend?'),
(3, 2, 'Interstellar gets my vote.'),
(3, 3, 'I''d rather visit a museum.'),
(4, 2, 'Best destination for summer?'),
(4, 7, 'I recommend Morocco 🇲🇦'),
(4, 5, 'Japan is on my bucket list.');

-- ============================================================
-- GROUP POSTS
-- ============================================================
INSERT INTO GROUP_POSTS (id, group_id, user_id, title, text) VALUES
(1, 1, 1, 'Favorite Go Feature', 'Mine is goroutines.'),
(2, 1, 6, NULL,                'Who''s using generics?'),
(3, 2, 6, 'Weekend Match',     'Who''s available this Saturday?'),
(4, 2, 8, NULL,                'I''ll reserve the field.'),
(5, 3, 8, 'Best Movie',        'Recommend your favorite movie.'),
(6, 3, 2, NULL,                'I''m reading Dune right now.'),
(7, 4, 2, 'Dream Destination', 'Where do you want to travel next?'),
(8, 4, 5, NULL,                'I want to visit Iceland.');

-- ============================================================
-- GROUP POST COMMENTS
-- ============================================================
INSERT INTO GROUP_POST_COMMENTS (id, group_post_id, user_id, text) VALUES
(1, 1, 4, 'Goroutines changed the way I write code.'),
(2, 2, 1, 'Generics are finally mature.'),
(3, 3, 7, 'Count me in!'),
(4, 4, 6, 'Perfect.'),
(5, 5, 3, 'The Godfather never gets old.'),
(6, 6, 8, 'Great choice!'),
(7, 7, 7, 'Japan is also my dream destination.'),
(8, 8, 2, 'Iceland looks incredible in winter.');

-- ============================================================
-- GROUP POST REACTIONS
-- ============================================================
INSERT INTO GROUP_POST_REACTIONS (group_post_id, user_id, is_like) VALUES
(1, 2, 1),  (1, 4, 1),  (1, 6, -1),
(2, 1, 1),  (2, 4, 1),
(3, 1, 1),  (3, 8, 1),
(4, 6, 1),
(5, 3, 1),
(6, 8, -1),
(7, 7, 1),
(8, 2, 1);

-- ============================================================
-- GROUP POST COMMENT REACTIONS
-- ============================================================
INSERT INTO GROUP_POST_COMMENT_REACTIONS (group_post_comment_id, user_id, is_like) VALUES
(1, 2, 1),  (1, 6, -1),
(2, 4, 1),
(5, 1, 1),
(6, 5, 1),
(7, 1, 1);

-- ============================================================
-- GROUP EVENTS (4 events)
-- ============================================================
INSERT INTO GROUP_EVENTS (id, group_id, creator_id, title, description, image, event_time) VALUES
(1, 1, 1, 'Go Meetup',     'Monthly Go developers meetup.',   '', datetime('now', '+7 days')),
(2, 2, 6, 'Football Match', 'Friendly football game.',         '', datetime('now', '+3 days')),
(3, 3, 8, 'Museum Visit',   'Visit the modern art museum together.', '', datetime('now', '+10 days')),
(4, 4, 2, 'Weekend Road Trip', 'Two-day trip to the mountains.', '', datetime('now', '+14 days'));

-- ============================================================
-- EVENT RESPONSES
-- ============================================================
INSERT INTO EVENT_RESPONSES (event_id, user_id, status) VALUES
(1, 1, 'going'),       (1, 4, 'going'),       (1, 6, 'not_going'),
(2, 6, 'going'),       (2, 1, 'going'),       (2, 8, 'not_going');

-- ============================================================
-- NOTIFICATIONS (22 notifications)
-- ============================================================
INSERT INTO NOTIFICATIONS (user_id, actor_id, type, object_type, object_id, is_read) VALUES
(3,  7,  'follow_request',   'follow',          7, 0),
(1,  2,  'post_reaction',    'post',            1, 0),
(1,  4,  'comment',          'comment',         2, 0),
(1,  3,  'follow_accepted',  'follow',          3, 1),
(1,  6,  'group_invite',     'group_invite',    2, 1),
(2,  1,  'post_reaction',    'post',            2, 0),
(2,  3,  'comment',          'comment',         1, 1),
(2,  7,  'follow_request',   'follow',          7, 0),
(2,  1,  'message',          'conversation',    1, 0),
(3,  2,  'post_reaction',    'post',            9, 0),
(3,  1,  'comment',          'comment',         5, 0),
(3,  1,  'follow_accepted',  'follow',          1, 1),
(3,  6,  'group_event',      'event',           2, 0),
(1,  8,  'follow_request',   'follow',          8, 0),
(2,  1,  'group_event',      'event',           1, 0),
(5,  8,  'follow_request',   'follow',          8, 0),
(7,  1,  'group_invite',     'group_invite',    1, 0),
(5,  6,  'group_invite',     'group_invite',    2, 0),
(1,  8,  'group_join_request','group_request',  1, 0),
(6,  4,  'group_join_request','group_request',  2, 0),
(4,  1,  'group_event',      'event',           1, 1),
(6,  1,  'group_event',      'event',           1, 0);
