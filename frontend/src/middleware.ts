import { auth } from "~/server/auth";
import { NextResponse } from "next/server";

// list every route that lives in (public)
const publicRoutes = ["/login", "/register"]; // , "/forgot-password"
const privateRoutes = ["/groups", "/messages", "/notifications", "/profile"]; // except root "/"

export default auth((req) => {
  const isLoggedIn = !!req.auth;
  const { pathname } = req.nextUrl;

  const isPublicRoute = publicRoutes.some(
    (route) => pathname === route // || pathname.startsWith(`${route}/`)
  );

  const isPrivateRoute = pathname === "/" || privateRoutes.some(
    (route) => pathname === route || pathname.startsWith(`${route}/`)
  );

  // 1. Not logged in + trying to access a private route → send to /login
  if (isPrivateRoute && !isLoggedIn) {
    const loginUrl = new URL("/login", req.nextUrl.origin);
    loginUrl.searchParams.set("callbackUrl", pathname);
    return NextResponse.redirect(loginUrl);
  }

  // 2. Logged in + trying to access a public route (login/register) → send to dashboard
  if (isPublicRoute && isLoggedIn) {
    return NextResponse.redirect(new URL("/", req.nextUrl.origin));
  }

  return NextResponse.next();
});

export const config = {
  // run on every route except static assets/api — so we can catch both groups
  matcher: ["/((?!api|_next/static|_next/image|favicon.ico).*)"],
};