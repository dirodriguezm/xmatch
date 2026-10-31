/**
 * Site-wide constants: navigation, external links and public URLs.
 */

export const SITE_NAME = "XWave";

/** Public origin of the web app, used for share links, embeds and OG cards. */
export const SITE_URL =
  process.env.NEXT_PUBLIC_SITE_URL ?? "https://xwave-rho.vercel.app";

export const API_ORIGIN = "https://xwave-astro.udp.cl";
export const SWAGGER_URL = `${API_ORIGIN}/swagger/index.html`;
export const REPO_URL = "https://github.com/dirodriguezm/xmatch";
export const ISSUES_URL = `${REPO_URL}/issues`;
export const DISCUSSIONS_URL = `${REPO_URL}/discussions`;

export interface NavItem {
  href: string;
  label: string;
}

/** Top-level navigation, in header order. */
export const NAV_ITEMS: NavItem[] = [
  { href: "/search", label: "Search" },
  { href: "/bulk", label: "Bulk" },
  { href: "/explore", label: "Explore" },
  { href: "/catalogs", label: "Catalogs" },
  { href: "/developers", label: "API" },
  { href: "/learn", label: "Learn" },
];

/** Footer link columns. */
export const FOOTER_GROUPS: { title: string; items: NavItem[] }[] = [
  {
    title: "Product",
    items: [
      { href: "/search", label: "Cone search" },
      { href: "/bulk", label: "Bulk cross-match" },
      { href: "/explore", label: "Explore" },
      { href: "/catalogs", label: "Catalogs" },
    ],
  },
  {
    title: "Developers",
    items: [
      { href: "/developers", label: "API playground" },
      { href: "/llms.txt", label: "llms.txt" },
      { href: SWAGGER_URL, label: "Swagger reference" },
      { href: "/status", label: "Status" },
    ],
  },
  {
    title: "Resources",
    items: [
      { href: "/learn", label: "Learn & FAQ" },
      { href: "/changelog", label: "Changelog" },
    ],
  },
  {
    title: "Project",
    items: [
      { href: "/about", label: "About" },
      { href: "/contact", label: "Contact" },
      { href: REPO_URL, label: "GitHub" },
    ],
  },
];

export function isExternalHref(href: string): boolean {
  return /^https?:\/\//.test(href);
}
