# Working with images

Images make or break a design's perceived quality — and they're where unpredictable real-world
content breaks layouts.

## Use good photos

A mediocre layout with great photography looks professional; a great layout with stock-photo clichés
looks cheap. Invest in (or carefully source) real, high-quality imagery.

## Text needs consistent contrast

Text placed directly on a photo is unreadable wherever the image is light (or dark). The image
changes per upload, so you can't rely on it. Guarantee contrast with one of:

- **A scrim / overlay** — a semi-transparent dark (or brand) layer between image and text.
- **Lower the image's contrast** or darken/brighten it uniformly so text always pops.
- **A gradient overlay** — fades the image behind the text only (great for captions at the bottom).
- **A text-shadow** — a soft, large-radius shadow lifts light text off a busy area.

Pick the approach that keeps contrast constant **regardless of which image loads**.

## Everything has an intended size

Assets are drawn for a size — respect it.

- ❌ **Never upscale** icons, avatars, logos, or favicons; they go blurry/pixelated.
- ✅ Need a small icon to fill a larger space? **Put it on a colored/rounded backdrop** at its native
  size instead of stretching it.
- Embedded logos and screenshots have a natural resolution — give them a container that frames them
  rather than blowing them up.

## Beware user-uploaded content

You don't control aspect ratio, dimensions, or colors of what users upload. Defend the layout:

- Constrain to a fixed box with **`background-size: cover`** (or `object-fit: cover`) so any aspect
  ratio fills the frame without distortion or layout shift.
- Cap dimensions and **prevent long/oversized content from breaking** the surrounding layout.
- Add a **subtle inner border** (a low-opacity dark inset) so a light image doesn't disappear against
  a light background.

## Apply it in Vaquita

- Receipt photos and user avatars are exactly this "uncontrolled content" case — always frame them in
  a fixed box with `object-fit: cover` and a subtle border; never trust the upload's dimensions.
- For **hashed-background avatars**, render a plain styled `<div>`, **not** `<Avatar.Fallback>` (the
  fallback inherits the accent color in light mode — see the `radix-ui` skill and the design memory).
- Use `next/image` for framing/optimization where applicable, with explicit width/height to avoid
  layout shift.
