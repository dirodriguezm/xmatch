# Designing text

Typography decisions repeat on every screen, so systematize them once.

## Establish a type scale

Picking sizes ad hoc gives you ten near-identical values (10–24px all used somewhere). Define a fixed
scale and choose from it.

- A **modular scale** (a ratio like 4:5 or the golden ratio) produces ugly fractional pixels
  (31.25px) and the wrong granularity for UI. **Hand-pick** instead. A scale that works for most
  projects:

  ```
  12  14  16  18  20  24  30  36  48  60  72
  ```

- **Use `px` or `rem`, never `em`.** `em` is relative to the current font size, so a nested element's
  computed size drifts off your scale (1.25em inside 1.25em ≠ a scale value).

## Use good fonts

- For UI, a **neutral sans-serif** is the safe bet. The system font stack is a solid default:
  `-apple-system, Segoe UI, Roboto, Noto Sans, Ubuntu, Cantarell, Helvetica Neue, sans-serif`.
- Heuristics for picking: favor families with **5+ weights** (more weights = more care), sort font
  directories by popularity, and steal from sites you admire (inspect their `font-family`).
- Pick fonts **optimized for the size** — UI text fonts have taller x-heights and looser spacing;
  avoid condensed display faces for body copy.

## Keep line length in check

The single biggest readability win: **45–75 characters per line**.

- In `em`, that's roughly **20–35em**; `max-width: ~34em` is a reliable target.
- Limit the **paragraph** width even when the surrounding content area is wider (e.g. text beside an
  image). Different widths in one block looks more polished, not less.

## Baseline, not center

When you mix font sizes on one line (a big title and a small action beside it), **baseline-align**
them, don't vertically center — centering makes the sizes look like a mistake.

## Line-height is proportional

Line-height should scale with both font size and measure: **tall** (1.5–2) for small or long-line
body text, **tight** (~1–1.2) for large headings. A single global line-height fits nothing well.

## Other levers

- **Not every link needs a color** — in link-dense UIs, weight or a subtle treatment can mark a link
  without a sea of blue.
- **Align for readability** — right-align numeric columns; mind ragged alignment in lists.
- **Letter-spacing** — tighten large headlines slightly; widen all-caps labels for legibility.

## Apply it in Vaquita

- Drive type through Radix `<Text>` / `<Heading>` with the **1–9 `size` scale** and the `radix-ui`
  weight rules (`medium` / `bold` only) — that's the project's type system; don't set raw `font-size`.
- **Inter** is the app font. Keep all user-facing copy in **Spanish** (`design-system` rule).
- Constrain long-form paragraphs with a `max-w-[34em]` (or the nearest Tailwind `max-w-*`) so line
  length stays in the 45–75 char band.
