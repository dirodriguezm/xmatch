# Layout and spacing

Spacing is where "clean" comes from. Most plain-looking UIs are just cramped and inconsistent.

## Start with too much white space, then remove

On the web, space is almost always _added_ until something stops looking actively bad — which leaves
it at the bare minimum. Invert it: give an element **way too much** room, then dial back until you're
happy. What feels like too much on one element reads as "just right" in the full UI.

Dense UIs are a deliberate choice (dashboards packing a lot on one screen), not a default. Make
density a decision, not an accident.

## Establish a spacing and sizing system

Never nitpick 120px vs 125px. Choose from a constrained, **pre-defined** scale.

- A **linear** "multiple of 4" scale is too weak — a few px matters enormously at small sizes and is
  invisible at large ones. Build a scale where **no two adjacent values are closer than ~25%**.
- Start from a `16px` base and use factors/multiples. The book's practical scale:

  ```
  4  8  12  16  24  32  48  64  96  128  192  256  384  512  640  768
  ```

  Packed at the small end, spreading out as it grows. Need more space? Grab the next value up.

## You don't have to fill the screen

A big canvas doesn't mean you must use it. If a form only needs 600px, use 600px.

- **Shrink the canvas** — design mobile-first on a ~400px frame; real constraints make small layouts
  easier, and they scale up with little change.
- **Think in columns** — if something works at a narrow width but feels lost in a wide UI, split it
  into columns (e.g. break supporting text into its own column beside a form) instead of stretching
  it. Don't make an element worse just to match a neighbor's width.

## Grids are overrated

A 12-column grid gives every element a _fluid, percentage_ width. But many elements should have a
**fixed** width (a sidebar at 240px, not 25%) and stop growing past a point. Don't outsource every
width decision to a grid.

## Relative sizing doesn't scale

Don't scale a whole component proportionally for big screens — headlines that grow 1:1 with the
viewport get absurd while body text balloons. Large elements should shrink _more_ than small ones on
small screens; let each thing size on its own curve.

## Avoid ambiguous spacing

When elements are separated only by space, there must be **more space around a group than within it**.
A heading sitting equidistant between its own paragraph and the one above is ambiguous — tighten the
gap to the content it belongs to.

## Apply it in Vaquita

- Use the **Tailwind spacing scale** (`p-4`, `gap-6`, `mt-8`) and brand-mapped values rather than
  arbitrary `p-[13px]` brackets — the scale is the system. Reach for brackets only for genuine
  non-token needs (`design-system` rule).
- Mobile-first: unprefixed utilities are the base; layer `md:` / `lg:` upward.
- Prefer Radix `<Flex>` / `<Grid>` with `gap` for grouped spacing so "more around than within"
  comes for free.
