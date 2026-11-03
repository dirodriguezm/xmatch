# Starting from scratch

How to begin a screen so it has good bones before you polish it.

## Start with a feature, not a layout

Don't open by designing the app shell (top nav vs sidebar, where the logo goes). You don't have the
information to make those calls yet. Pick one concrete piece of functionality — "create a receipt",
"join a split", "see who paid" — and design just that. The shell falls out once a few features exist.

## Detail comes later — work in grayscale first

In the earliest passes, **don't** make low-level decisions about typefaces, shadows, icons, or color.

- Sketch rough (paper or low-fi) to explore layouts fast.
- **Hold the color.** Design in grayscale and force size, weight, contrast, and spacing to do the
  hierarchy work. If it reads well in grayscale, color only makes it better. Color added early
  papers over weak hierarchy.

## Don't over-invest, work in short cycles

Wireframes are disposable — users can't use a static mockup. Use them to make a decision, then move
to the real thing. Design in the browser early so the constraints are real (✅ a 400px canvas beats
a 1440px artboard for a mobile-first screen).

## Choose a personality

Decide the voice up front and commit: the **font** (neutral vs characterful), the **accent color**,
the **border-radius** (sharp = serious, round = friendly), and the **language/tone** of copy. These
four choices set the feel; keep them consistent across the surface.

## Apply it in Vaquita

The personality is already chosen — **don't reinvent it.** Vaquita is mobile-first, Inter, the
`--vaquita-blue` accent, Radix's default radii, and a warm Spanish voice. Start a new screen by
composing the existing Radix Themes components in grayscale-equivalent (default tokens) first, then
layer the brand accent. The personality decisions live in `vaquita-frontend/reference/design-system.md`
and the `radix-ui` skill — consume them rather than picking new ones.
