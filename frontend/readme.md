what is stopping me from writing :
// import Sidebar from '@/components/Sidebar';
?

➜  social-network npm create t3-app@latest web

> npx
> create-t3-app web

   ___ ___ ___   __ _____ ___   _____ ____    __   ___ ___
  / __| _ \ __| /  \_   _| __| |_   _|__ /   /  \ | _ \ _ \
 | (__|   / _| / /\ \| | | _|    | |  |_ \  / /\ \|  _/  _/
  \___|_|_\___|_/‾‾\_\_| |___|   |_| |___/ /_/‾‾\_\_| |_|


│
◇  Will you be using TypeScript or JavaScript?
│  TypeScript
│
◇  Will you be using Tailwind CSS for styling?
│  Yes
│
◇  Would you like to use tRPC?
│  No
│
◇  What authentication provider would you like to use?
│  NextAuth.js
│
◇  What database ORM would you like to use?
│  None
│
◇  Would you like to use Next.js App Router?
│  Yes
│
◇  Would you like to use ESLint and Prettier or Biome for linting and
formatting?
│  ESLint/Prettier
│
◇  Should we initialize a Git repository and stage the changes?
│  No
│
◇  Should we run 'npm install' for you?
│  No
│
◇  What import alias would you like to use?
│  ~/

✔ web scaffolded successfully!

Adding boilerplate...
✔ Successfully setup boilerplate for nextAuth
✔ Successfully setup boilerplate for tailwind
✔ Successfully setup boilerplate for envVariables
✔ Successfully setup boilerplate for eslint

Next steps:
  cd web
  npm install
  Fill in your .env with necessary values. See https://create.t3.gg/en/usage/first-steps for more info.
  npm run dev
  git init
  git commit -m "initial commit"


===========================================================

Required

✅ TypeScript
✅ Tailwind CSS

Recommended

✅ ESLint
✅ Prettier

Optional (depends on your needs)
Auth.js
Choose Yes if users need accounts/login.
(default with App Router).

================================================

go run ./cmd/server


==================================================

1. add types/css.d.ts
2. add icons library:
  @import "@tabler/icons-webfont/dist/tabler-icons.min.css";
  npm install @tabler/icons-webfont