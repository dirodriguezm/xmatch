# Hierarchy is everything

Visual hierarchy — how important each element _looks_ relative to the others — is the highest-leverage
tool for making a UI feel designed. When everything competes for attention, it reads as noise.

## Size isn't everything

Relying on font size alone gives you oversized headings and unreadably tiny secondary text. Reach for
**weight** and **color** first.

- **Two weights** are usually enough: a normal weight (400–500) for most text and a heavier weight
  (600–700) for emphasis. A bolder-but-smaller label often beats a big light one.
- **Three text colors**: dark (primary — a heading), grey (secondary — supporting copy), lighter grey
  (tertiary — a timestamp or footnote).
- ❌ Never drop **below weight 400** to de-emphasize — thin text is hard to read. De-emphasize with a
  **lighter color** or **smaller size** instead.

## Don't use grey text on a colored background

Grey-on-white works because it reduces contrast. On a colored panel, grey looks dim and muddy.

- ✅ **Hand-pick a new color** with the **same hue** as the background, nudging saturation/lightness
  until it sits "one step back".
- ❌ Don't use white at reduced opacity — it desaturates, looks disabled, and lets a background image
  bleed through the text.

## Emphasize by de-emphasizing

When the primary element won't stand out and there's nothing left to add to it, **mute its
competition** instead. Give inactive nav items a softer color; drop the background fill on a sidebar
so the main content leads.

## Labels are a last resort

`Label: value` gives every datum equal weight and kills hierarchy.

- Let **format/context** identify data (`jane@x.com` is obviously an email; "Customer Support" under
  a name is obviously a role).
- **Combine** label and value into one phrase: "12 left in stock", "3 bedrooms", "82 BPM".
- When you _do_ need a label (scannable dashboards, spec tables), treat it as **secondary** —
  smaller/lighter than the value it describes.

## Visual hierarchy ≠ document hierarchy

Pick HTML tags for **semantics**, style them for the **eye**. An `h1` titling a settings page can be
small and quiet; a section title is often a label, not a headline. Don't let the tag dictate the size.

## Apply it in Vaquita

- Use Radix `<Text>` / `<Heading>` `weight` — **only `medium` and `bold` are valid** (`weight="semibold"`
  fails the build; see the `radix-ui` skill).
- Map the three text colors to the greys in `app/globals.css`: `--vaquita-black`/default ink for
  primary, `--vaquita-gray-dark` for secondary, `--vaquita-gray-medium` for tertiary.
- For grey-on-colored, pick a same-hue tint of the panel via a CSS var — never reduce opacity, never
  hardcode hex (`design-system` rule).
- Receipt rows are the canonical hierarchy surface: emphasize the amount, de-emphasize labels.
