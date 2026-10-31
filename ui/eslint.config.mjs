import vitest from "@vitest/eslint-plugin";
import { defineConfig, globalIgnores } from "eslint/config";
import nextVitals from "eslint-config-next/core-web-vitals";
import nextTs from "eslint-config-next/typescript";
import prettier from "eslint-config-prettier";
import noUnneededEffect from "eslint-plugin-react-you-might-not-need-an-effect";
import simpleImportSort from "eslint-plugin-simple-import-sort";

const eslintConfig = defineConfig([
  ...nextVitals,
  ...nextTs,
  prettier,
  // Effects used for derived state, event handling or prop sync. Warnings only:
  // post-hydration init and Aladin/scroll sync are legitimate and get flagged.
  noUnneededEffect.configs.recommended,
  // Tests: no stray .only, no test without an assertion, no conditional expect.
  {
    files: ["app/**/*.test.ts"],
    plugins: { vitest },
    rules: vitest.configs.recommended.rules,
  },
  // Custom rules
  {
    plugins: {
      "simple-import-sort": simpleImportSort,
    },
    rules: {
      // Forbid inline styles - use Tailwind instead
      "react/forbid-dom-props": ["error", { forbid: ["style"] }],
      "react/forbid-component-props": ["error", { forbid: ["style"] }],
      // No console.log (allow console.warn/error)
      "no-console": ["error", { allow: ["warn", "error"] }],
      // Handle the undefined case instead of asserting it away. Warning until
      // the existing ones are cleaned up.
      "@typescript-eslint/no-non-null-assertion": "warn",
      // Import sorting
      "simple-import-sort/imports": "error",
      "simple-import-sort/exports": "error",
    },
  },
  // Override default ignores of eslint-config-next.
  globalIgnores([
    // Default ignores of eslint-config-next:
    ".next/**",
    "out/**",
    "build/**",
    "next-env.d.ts",
    // Utility scripts using CommonJS
    "scripts/**",
  ]),
]);

export default eslintConfig;
