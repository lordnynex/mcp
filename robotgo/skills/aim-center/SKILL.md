---
name: aim-center
description: Click the geometric center of on-screen controls, not edges or the first glyph.
---

# Aim at the center

Clicks on a 1px border, a placeholder's first letter, or a Send chevron often miss. Read a bounding box from the screenshot, map it into mouse space, then click its interior.

## Map then point

Image pixels (`ix`, `iy`) are not mouse coordinates when the capture is Retina-scaled or `maxWidth`-shrunk:

`mouseX = originX + ix * mouseWidth / width`

`mouseY = originY + iy * mouseHeight / height`

Then click the geometric center of the mapped box: `(x + w/2, y + h/2)`.

Inset from chrome:

- Not the first glyph of a placeholder (`Add a follow-up`, `Search`).
- Not the Send button, model picker, or a dropdown chevron.
- Not a 1px window or field border.

Text fields: click the **empty interior** of the field — slightly right of a leading icon, vertically centered in the field height. Composer boxes (Cursor follow-up): the dark rounded rect **above** Agent/Send, not the placeholder start.

If the last click produced no focus change, inset **40–80px** into the box and click once more. Do not take a series of tighter crops.

## Select

Highlighting text is not a click on the first character. Use `MouseSelect` from one interior point to another, or `MousePath` with `hold` for a polyline. Double-click / `count` 3 selects a word or line.

Do not invent coordinates. Use the image and mapping fields that were just returned.
