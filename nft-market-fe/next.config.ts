import type { NextConfig } from "next";
import path from "path";

const nextConfig: NextConfig = {
  /* config options here */
  async rewrites() {
    // if (!isProd) {
    return [
      {
        source: "/api/v1/:path*",
        destination: "http://localhost:80/api/v1/:path*", // 开发环境
      },
    ];
    // }
    // return [];
  },
  images: {
    remotePatterns: [
      {
        protocol: "https",
        hostname: "**",
      },
      {
        protocol: "http",
        hostname: "**",
      },
    ],
  },
  i18n: {
    locales: ["en-US", "zh-CN"],
    defaultLocale: "en-US",
    localeDetection: false,
  },
  typescript: {
    ignoreBuildErrors: true,
  },

  // 关闭 ESLint
  eslint: {
    ignoreDuringBuilds: true,
  },

  output: "standalone",
};

export default nextConfig;
