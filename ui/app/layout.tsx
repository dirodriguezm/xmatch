import "./globals.css";

import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";

import { SITE_URL } from "@/app/lib/constants/site";

import { Providers } from "./providers";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  metadataBase: new URL(SITE_URL),
  title: {
    default: "XWave",
    template: "XWave | %s",
  },
  description:
    "Cross-match any sky position against Gaia DR3, AllWISE and eROSITA in milliseconds — with light curves, SEDs, code snippets and citations.",
  alternates: {
    types: { "application/rss+xml": "/changelog/rss.xml" },
  },
  icons: {
    icon: "/xwave-icon.svg",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="dark">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
