import { type DefaultSession, type NextAuthConfig } from "next-auth"; // defaultsession ?!
import Credentials from "next-auth/providers/credentials";
// import DiscordProvider from "next-auth/providers/discord";
import { env } from "~/env";

/**
 * Module augmentation for `next-auth` types. Allows us to add custom properties to the `session`
 * object and keep type safety.
 *
 * @see https://next-auth.js.org/getting-started/typescript#module-augmentation
 */
declare module "next-auth" {
  interface Session extends DefaultSession {
    user: { // email already exists !
      id: string;
      firstname: string;
      lastname: string;
      birthdate: string;
      nickname: string;
      aboutme: string;
      avatar: string;
      // ...other properties
      // role: UserRole;
    } & DefaultSession["user"];
  }

  // interface User {
  //   // ...other properties
  //   // role: UserRole;
  // }
}

/**
 * Options for NextAuth.js used to configure adapters, providers, callbacks, etc.
 *
 * @see https://next-auth.js.org/configuration/options
 */
export const authConfig = {
  providers: [
    Credentials({
      name: "Credentials",
      credentials: {
        identifier: { label: "Identifier", type: "text" },
        password: { label: "Password", type: "password" },
      },
      authorize: async (credentials) => {
        if (!credentials?.identifier || !credentials?.password) return null;

        const res = await fetch(`${env.GO_BACKEND_URL}/api/login`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            identifier: credentials.identifier,
            password: credentials.password,
          }),
        });

        // console.log(res); // remove cors headers !!!

        if (!res.ok) return null; // wrong credentials, backend returned 401 etc.

        const data = (await res.json()).data;
        // Expecting your Go backend to return something like:
        // { id, email, name, token: "jwt-from-go" }
        // console.log(data)
        if (!data?.id) return null; // why id is required

        return {
          id: String(data.id),
          firstname: data.firstname,
          lastname:  data.lastname,
          email:     data.email,
          birthdate: data.birthdate,
          nickname:  data.nickname, // nickname can be empty, so should pass first and last name instead !
          aboutme:   data.aboutme,
          avatar:    data.avatar,
          accessToken: data.token, // your Go JWT, if you issue one
          // not a jwt token ?!!!
        };
      },
    }),
    // DiscordProvider,
    /**
     * ...add more providers here.
     *
     * Most other providers require a bit more work than the Discord provider. For example, the
     * GitHub provider requires you to add the `refresh_token_expires_in` field to the Account
     * model. Refer to the NextAuth.js docs for the provider you want to use. Example:
     *
     * @see https://next-auth.js.org/providers/github
     */
  ],
  session: { strategy: "jwt" }, // required for Credentials — no DB sessions
  callbacks: {
    jwt: ({ token, user }) => { // required
      if (user) {
        token.id = user.id;
        token.firstname = (user as any).firstname;
        token.lastname = (user as any).lastname;
        token.birthdate = (user as any).birthdate;
        token.nickname = (user as any).nickname;
        token.aboutme = (user as any).aboutme;
        token.avatar = (user as any).avatar;
        token.accessToken = (user as any).accessToken;
      }
      return token;
    },
    session: ({ session, token }) => ({
      // generated with ai :
      // session.user.id = token.id as string;
      // (session as any).accessToken = token.accessToken;
      // return session;
      ...session,
      user: {
        ...session.user,
        id: token.sub,
        firstname: token.firstname,
        lastname: token.lastname,
        birthdate: token.birthdate,
        nickname: token.nickname,
        aboutme: token.aboutme,
        avatar: token.avatar,
        session_id:token.accessToken
      },
    }),
  },
  pages: {
    signIn: "/login", // optional custom login page
  },
} satisfies NextAuthConfig;
