# Creating depth

Depth makes a flat screen feel tactile and helps the eye rank elements on a z-axis.

## Emulate a light source

Our brains read depth from how light hits an object, and on screens **light comes from above**. To
make an element look raised or inset, decide its _profile_, then light it consistently.

- **Raised** (button, card): the top edge faces the sky, so make it **slightly lighter** than the
  face (a subtle top highlight via a top border or inset light shadow). The element blocks light
  below it, so add a **small dark shadow** just beneath.
- **Inset** (well, text input, pressed state): the top is in shadow (**darker top edge**) and the
  bottom lip catches light (**lighter bottom edge**).
- **Hand-pick the highlight color** — overlaying translucent white sucks the saturation out of the
  surface underneath. Keep shadows fairly **sharp** (small blur); a few px is plenty.

## Use shadows to convey elevation

A shadow positions an element on a virtual z-axis. **Tight, small** shadow = barely lifted; **large,
blurry** shadow = floating close to the user. Closer-feeling things grab more attention, so match the
shadow to the element's importance:

- Button: small. Dropdown: medium. Modal: large.
- Define a **5-step elevation scale** (small → large, increasing roughly linearly), e.g.

  ```css
  box-shadow: 0 1px 3px hsla(0, 0%, 0%, 0.2); /* 1 — resting */
  box-shadow: 0 4px 6px hsla(0, 0%, 0%, 0.2); /* 2 */
  box-shadow: 0 5px 15px hsla(0, 0%, 0%, 0.2); /* 3 */
  box-shadow: 0 10px 24px hsla(0, 0%, 0%, 0.2); /* 4 */
  box-shadow: 0 15px 35px hsla(0, 0%, 0%, 0.2); /* 5 — modal */
  ```

- **Shadows on interaction**: raise the shadow when an item is picked up/dragged; shrink or remove it
  to make a button feel pressed on click. Think "where on the z-axis", then assign the shadow.

## Shadows can have two parts

A great shadow is often **two** layered shadows doing different jobs:

- A **larger, softer** shadow with a bigger vertical offset and blur — the _direct-light_ cast shadow.
- A **tighter, darker** shadow with little offset — the _ambient_ shadow hugging the element's edge.

```css
box-shadow:
  0 4px 6px hsla(0, 0%, 0%, 0.7),
  0 5px 15px hsla(0, 0%, 0%, 0.1);
```

As an element rises, the **ambient (tight, dark) part fades** while the soft part grows — distinct at
low elevation, almost gone at high elevation.

## Even flat designs have depth

"Flat" doesn't mean depthless. Convey layers with **background color and spacing** instead of
shadows — a slightly different panel color reads as "above" the page.

## Overlap elements to create layers

Letting one element overlap another (a card straddling a colored header, an avatar crossing a banner
edge) creates depth without any shadow at all.

## Don't get carried away

Borrow cues from the real world for a hint of depth — don't chase photorealism; it gets busy.

## Apply it in Vaquita

- Radix `<Card>` and Themes components already carry sensible elevation — prefer them over bespoke
  shadows (see the `radix-ui` skill).
- If you need a custom shadow, define a **small fixed set** (don't one-off per component) and express
  shadow colors with `hsla()`, matching the elevation scale above.
- Use `--vaquita-gray-subtle` / `-light` as "raised panel" background colors for the flat-depth
  approach instead of reaching for a shadow.
