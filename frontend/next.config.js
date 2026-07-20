/**
 * Run `build` or `dev` with `SKIP_ENV_VALIDATION` to skip env validation. This is especially useful
 * for Docker builds.
 */
import "./src/env.js";

/** @type {import("next").NextConfig} */
const config = {
    output: process.env.NODE_ENV === "production" ? "standalone" : undefined,
    async rewrites() {
    return [
      {
        source: "/uploads/:path*",
        destination: `${process.env.GO_BACKEND_URL}/uploads/:path*`,
      },
    ];
  },
};

export default config;
