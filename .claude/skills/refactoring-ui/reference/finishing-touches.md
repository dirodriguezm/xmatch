# Finishing touches

Small, cheap details that take a correct-but-plain screen to "designed". Reach for these once
hierarchy, spacing, type, and color are solid.

## Supercharge the defaults

Don't add new elements — upgrade the ones already there.

- Replace **bullet points with icons** (checkmarks for features, a padlock for security items).
- **Promote quotes/testimonials** into visual elements: bigger, colored, with quotation marks.
- Style **links** beyond default blue (weight, a custom underline that overlaps the text).
- Use **custom checkboxes/radios** in a brand color instead of browser defaults.

## Add color with accent borders

A thin colored border is the cheapest way to add brand flair to a bland block:

- Across the **top of a card**.
- As an **active-nav underline**.
- Down the **side of an alert** (left border).
- A short accent **under a heading**.
- Across the **top of the whole layout**.

## Decorate your backgrounds

Break monotony without redesigning — keep contrast **low** so nothing fights the content:

- **Change the background color** of a section/panel to separate it.
- A **slight gradient** — two hues **≤30° apart** for a subtle, non-garish blend.
- A **repeating pattern** (whole background or along one edge).
- A **simple shape or illustration** (geometric accents, a faint map/dot grid behind a section).

## Don't overlook empty states

The first thing a new user sees is often the empty state — design it, don't ship "No items found".

- Add an **illustration or icon** and a clear **call-to-action** to the next step.
- **Hide** supporting chrome (tabs, filters, search) until there's content to act on.

## Use fewer borders

Borders pile up fast and make a UI busy. Before adding one to separate elements, try:

- a **box-shadow**,
- two **different background colors**, or
- simply **more spacing**.

Often no divider is needed at all.

## Think outside the box

Let elements escape their rectangle — overlap a card past a section edge, let an image bleed, break
the grid intentionally for a focal point.

## Apply it in Vaquita

- Build accents from the brand tokens (`--vaquita-blue` top borders, active-nav underlines) — via CSS
  vars / Tailwind utilities, **never hardcoded hex** (`design-system` rule).
- Vaquita already has empty states in receipts/groups — keep them first-class (illustration + CTA),
  not afterthoughts.
- Prefer Radix elevation/`Card` background separation over adding borders, in line with the repo's
  "avoid carditis / use compact sections" guidance.
