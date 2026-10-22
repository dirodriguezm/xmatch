---
name: refactoring-ui
description: Tactical visual-design rules from Adam Wathan & Steve Schoger's "Refactoring UI", expressed in Vaquita's design system. Use when a screen looks plain, busy, or "undesigned" and you need concrete fixes for hierarchy, spacing, type, color, depth, or finishing — not when wiring data or behavior. Trigger on "make this look better", "make it look designed", "polish the UI", "looks plain/bland/flat/cluttered/off", "visual hierarchy", "emphasize/de-emphasize", "white space", "spacing system", "type scale", "font size", "font weight", "line length", "color palette", "HSL", "shades", "grey text on color", "shadow", "elevation", "depth", "raised/inset", "empty state", "accent border", "decorate background", "fewer borders", "icon bullets".
---

# Refactoring UI

Tactical rules for making interfaces look _designed_ — distilled from **Refactoring UI** (Adam
Wathan & Steve Schoger) and mapped onto Vaquita's tokens, Radix usage, and conventions. This is the
"why does this look off and what knob do I turn" skill. It is **not** about data, routing, or
behavior (those live in the `vaquita-frontend` skill).

> Brand tokens, Tailwind composition, the no-`dark:` rule, and Spanish-copy rules are
> **design-system** concerns — they live in `vaquita-frontend/reference/design-system.md`. Radix
> component specifics (valid `weight`s, `Skeleton`, `Avatar.Fallback`, sizing scale) live in the
> **`radix-ui`** skill. This skill owns the _visual-design judgment_ and points at those for the
> how-to-implement-in-Vaquita details.

## Contents

Each chapter has a `reference/` file with the full rules + Vaquita "apply it" notes. Read the one
that matches the problem.

- Starting from scratch → `reference/starting-from-scratch.md`
- Hierarchy is everything → `reference/hierarchy.md`
- Layout and spacing → `reference/layout-and-spacing.md`
- Designing text → `reference/designing-text.md`
- Working with color → `reference/working-with-color.md`
- Creating depth → `reference/creating-depth.md`
- Working with images → `reference/working-with-images.md`
- Finishing touches → `reference/finishing-touches.md`

## The rules at a glance

### Starting from scratch

- Design a **feature**, not the app shell. Nav comes after you know what the pages hold.
- **Work in low fidelity first** — grayscale, rough spacing — then add color. Color hides weak
  hierarchy; grayscale forces you to fix it with size, weight, and spacing.
- Don't design too much up front; work in short cycles and design in the browser.
- Pick a **personality** (font, accent, border-radius, voice) and commit to it.

### Hierarchy is everything

- **Size isn't everything.** Lean on **font weight** and **color** before reaching for a bigger
  font. ~2 weights (`medium` 400–500 + `bold` 600–700), 3 text colors (primary / secondary /
  tertiary). Never go below weight 400 to de-emphasize — use a lighter _color_ or smaller size.
- **Never grey text on a colored background** — hand-pick a same-hue color closer to the background.
  Don't use white-at-reduced-opacity (it washes out and lets the background bleed through).
- **Emphasize by de-emphasizing**: when the hero won't pop, mute everything around it.
- **Labels are a last resort.** Prefer combining label+value ("12 left in stock", "3 bedrooms") over
  `Label: value`. Let format/context identify data.
- **Visual hierarchy ≠ document hierarchy.** An `h1` need not be big; a section title is often a
  small, quiet label. Pick the tag for semantics, style it for the eye.

→ `reference/hierarchy.md`

### Layout and spacing

- **Start with too much white space, then remove.** Cramped-but-not-broken is the default failure.
- **Use a spacing/sizing system, never arbitrary pixels.** Non-linear, no two values closer than
  ~25%. The book's 16-base scale: `4 8 12 16 24 32 48 64 96 128 192 256 384 512 640 768`.
- **You don't have to fill the screen.** Give each element the width it needs; shrink the canvas
  (design mobile-first); split into columns instead of stretching.
- **Grids are overrated** — fixed widths often beat fluid percentage columns.
- **Avoid ambiguous spacing**: more space _around_ a group than _between_ its items.

→ `reference/layout-and-spacing.md`

### Designing text

- **Establish a type scale** (hand-picked, not a modular ratio):
  `12 14 16 18 20 24 30 36 48 60 72`. Use **px/rem**, never `em` (em compounds when nested).
- **Use good fonts**: a neutral sans for UI (or the system stack); favor families with **5+ weights**.
- **Line length 45–75 characters** (~`20–35em`, `max-width: ~34em`) — even when the surrounding
  column is wider.
- **Baseline-align, not center**, when mixing font sizes on one line.
- **Line-height is proportional**: taller for small/long text, tighter for large headings.

→ `reference/designing-text.md`

### Working with color

- **Use HSL**, not hex — it maps to how the eye reads color (hue / saturation / lightness).
- **You need more colors than you think**: 8–10 greys, 5–10 shades of each primary, plus accent and
  semantic (red/yellow/green) families.
- **Define shades up front** (9 steps, `100`→`900`, base `500` = a good button background; fill the
  gaps). Don't `lighten()`/`darken()` on the fly.
- **Don't let lightness kill saturation** — raise saturation as lightness moves away from 50%; nudge
  hue toward a bright hue (60/180/300) to lighten or a dark hue (0/120/240) to darken (≤20–30°).
- **Greys can be tinted**; **don't rely on color alone** to convey meaning.

→ `reference/working-with-color.md`

### Creating depth

- **Emulate a light source from above.** Raised = lighter top edge + shadow below. Inset = darker
  top edge + lighter bottom. Hand-pick the highlight color; don't overlay translucent white.
- **Shadows convey elevation** — a 5-step shadow scale: small/tight = barely raised, large/blurry =
  close to the user. Lift a card on drag; flatten a button on press.
- **Two-part shadows**: a large soft "direct light" shadow + a tight dark "ambient" shadow; the
  ambient one fades as the element rises.
- **Even flat designs have depth** — use background color and spacing instead of shadows.

→ `reference/creating-depth.md`

### Working with images

- **Text over images needs guaranteed contrast**: a scrim/overlay, a knocked-back image, a longer
  text-shadow, or a gradient — never raw text on an unmodified photo.
- **Everything has an intended size** — never upscale icons, avatars, or logos. Put a small icon on
  a colored backdrop instead of blowing it up.
- **Beware user-uploaded content**: constrain with `background-size: cover`, prevent layout breakage,
  add a subtle inner border so light images don't vanish on a light background.

→ `reference/working-with-images.md`

### Finishing touches

- **Supercharge the defaults**: icon bullets, promoted quotes, custom checkboxes/links.
- **Add color with accent borders**: top of a card, an active-nav underline, the side of an alert.
- **Decorate backgrounds** sparingly: a color block, a ≤30° gradient, a low-contrast pattern or shape.
- **Don't overlook empty states** — illustration + clear CTA; hide tabs/filters until there's content.
- **Use fewer borders** — reach for a shadow, a different background color, or extra spacing first.

→ `reference/finishing-touches.md`

## Where this overlaps Vaquita (source of truth)

This skill teaches _judgment_; the repo already encodes most of these decisions. Don't re-derive or
contradict them — consume them:

- **Color** → the `--vaquita-blue` / `--vaquita-blue-hover` / `-active` / `-light` ladder and the
  `--vaquita-gray-*` greys in `app/globals.css` already are a hand-picked, HSL-spirit shade system.
  Consume the CSS vars; **never hardcode hex** (`design-system` rule).
- **Type & weight** → use Radix `<Text>`/`<Heading>` with the 1–9 size scale; only `medium`/`bold`
  weights are valid (`weight="semibold"` fails the build — `radix-ui` skill).
- **Spacing** → use the existing Tailwind spacing scale and mobile-first `md:` breakpoints, not
  arbitrary bracket values.
- **Avatars / hashed backgrounds** → render a plain styled `<div>`, not `<Avatar.Fallback>`
  (the fallback inherits accent color in light mode — `radix-ui` skill).

## Keeping this skill current

Update a `reference/` file when a book rule is **deliberately overridden** by a Vaquita decision
(e.g. the team standardizes a shadow token set, or picks a denser spacing scale for a specific
surface). Name the file and the change in your edit — don't silently drift. If a rule here ever
contradicts `design-system.md` or the `radix-ui` skill, **those win** for implementation specifics;
fix the pointer here.
