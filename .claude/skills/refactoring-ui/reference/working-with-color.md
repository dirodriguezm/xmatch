# Working with color

A real UI needs far more color than a "5-color palette" generator gives you. Build a system.

## Ditch hex for HSL

Hex/RGB hide the relationships between colors — two visually-similar blues look nothing alike in code.
**HSL** describes color the way the eye reads it:

- **Hue** (0–360°): position on the wheel (0 red, 120 green, 240 blue).
- **Saturation** (0–100%): how vivid (0 = grey).
- **Lightness** (0–100%): 0 black, 50 the pure hue, 100 white.

(Don't confuse HSL with HSB/HSV — browsers speak HSL.)

## You need more colors than you think

Plan for three categories, each with multiple shades:

- **Greys** — text, backgrounds, panels, borders. You'll want **8–10**. Start from a very dark grey
  (true black looks unnatural) up to near-white.
- **Primary** — one, maybe two brand colors for primary actions and active states; **5–10 shades**
  (light tints for backgrounds, dark shades for text).
- **Accent / semantic** — red (destructive), yellow (warning), green (success), plus any
  highlight/category colors; multiple shades each, used sparingly.

A complex UI can need ~10 colors × 5–10 shades.

## Define your shades up front

Don't generate shades on the fly with `lighten()`/`darken()` — you end up with 35 almost-identical
blues. Define a fixed ladder you choose from.

- **9 steps** is a sweet spot: name them `100` (lightest) → `900` (darkest).
- Pick the **base `500`** first — a shade that works as a **button background**.
- Pick the **edges**: darkest `900` (reserved for text) and lightest `100` (a tinted background).
- Fill the gaps by halving: `300` and `700` are the midpoints, then `200/400/600/800`.

## Don't let lightness kill your saturation

In HSL, the same saturation looks **less colorful** near 0% or 100% lightness. To keep light/dark
shades from washing out, **increase saturation as lightness moves away from 50%**.

You can also change _perceived_ brightness by **rotating the hue** (each hue has its own perceived
brightness):

- To **lighten**: rotate toward the nearest bright hue — **60° (yellow), 180° (cyan), 300° (magenta)**.
- To **darken**: rotate toward the nearest dark hue — **0° (red), 120° (green), 240° (blue)**.
- Keep rotations **≤20–30°** or it becomes a different color (e.g. yellow → warm orange darks instead
  of dull brown).

## Two more rules

- **Greys don't have to be grey** — tint them slightly (cool or warm) for a less clinical feel.
- **Accessible ≠ ugly** — large text can use lower contrast than the 4.5:1 body-text minimum, so you
  don't have to make everything dark. And **don't rely on color alone** to carry meaning (add an
  icon/label).

## Apply it in Vaquita

The repo **already is** this system — consume it, don't rebuild it:

- Primary ladder: `--vaquita-blue` (base), `--vaquita-blue-hover`, `--vaquita-blue-active` (darker
  shades), `--vaquita-blue-light` (the `100` tint for backgrounds/alerts).
- Greys: `--vaquita-gray-dark` / `-medium` / `-light` / `-subtle` in `app/globals.css`.
- **Consume the CSS vars; never hardcode hex** (`design-system` rule). Need a new shade or accent?
  Add it to `app/globals.css` (and `tailwind.config.js` for a utility) first.
- Radix Themes' own color scales cover semantic states — prefer them over inventing red/green hexes.
