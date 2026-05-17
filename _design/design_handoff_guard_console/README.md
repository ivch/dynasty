# Handoff: Guard Console Redesign

## Overview

This is a redesign of the **ЖК Династія** КПП (checkpoint) guard console — the web interface used by security staff to view, filter, and action resident requests (guests, taxis, deliveries, cargo). The redesign applies the **Dynasty Design System** consistently throughout.

---

## About the Design Files

The files in this bundle are **HTML design references** — high-fidelity prototypes showing intended look, layout, and interactive behaviour. They are **not** production code to copy directly.

The task is to **recreate these designs inside the existing `dynasty/_ui/guard/` Go template + Bootstrap 4 environment**, replacing the inline styles and scattered overrides in the original `index.html` with a clean, token-based CSS layer that follows the Dynasty Design System.

---

## Fidelity

**High-fidelity.** The prototypes use final colors, typography, spacing, and interactions drawn directly from the Dynasty Design System tokens. The developer should match these pixel-for-pixel using the codebase's existing Bootstrap 4 setup plus a new `guard-ds.css` stylesheet that layers the design tokens on top.

---

## Screens / Views

### 1 — Guard Console (single screen)

The entire UI lives on one page. It auto-refreshes every 60 seconds (`<meta http-equiv="refresh" content="60">`).

**Layout (top → bottom):**
```
┌─────────────────────────────────────────────────────────┐
│ Navbar (sticky, full-width, gradient)                   │
├─────────────────────────────────────────────────────────┤
│ .guard-container  (max-width: 1140px, centered)         │
│  ┌─────────────────────────────────────────────────┐    │
│  │ Summary strip (open / closed / total)           │    │
│  └─────────────────────────────────────────────────┘    │
│  ┌─────────────────────────────────────────────────┐    │
│  │ Filter bar (pills left, apt search right)       │    │
│  └─────────────────────────────────────────────────┘    │
│  ┌─────────────────────────────────────────────────┐    │
│  │ Request table (shadow card)                     │    │
│  │  row · row · row · …                            │    │
│  └─────────────────────────────────────────────────┘    │
│  Pager (centered below table)                           │
└─────────────────────────────────────────────────────────┘
```

---

## Components

### Navbar

- **Background:** `linear-gradient(0deg, #22C1C3 0%, #1F9991 100%)`
- **Position:** `sticky; top: 0; z-index: 100`
- **Box-shadow:** `0 2px 8px rgba(31,153,145,.25)`
- **Padding:** `12px 20px`
- **Layout:** `display: flex; align-items: center; gap: 16px`
- **Logo:** `assets/img/logo.png` at `width: 30px; height: 30px`
- **Brand text:** `"ЖК Династія"` — `font-size: 18px; font-weight: 600; color: #F5F5F5`
- **Sub-label:** `"Система заявок · КПП"` — `font-size: 12px; color: rgba(245,245,245,.75); text-transform: uppercase; letter-spacing: 0.06em`
- **Right side:** `"Оновлення кожні 60 с"` — `margin-left: auto; font-size: 12px; color: rgba(245,245,245,.6)`

---

### Summary Strip

- **Background:** `#fff`
- **Border-radius:** `0.25rem`
- **Box-shadow:** `0 .125rem .25rem rgba(0,0,0,.075)`
- **Padding:** `12px 16px`
- **Margin-bottom:** `16px`
- **Layout:** `display: flex; gap: 24px; align-items: center; flex-wrap: wrap`

Three stats, separated by `1px solid #dee2e6` vertical dividers (`height: 32px`):

| Stat | Value source | Number color |
|---|---|---|
| Відкритих | count of `status === "new"` | `#1F9991` (brand-teal-600) |
| Закритих | count of `status === "closed"` | `#434652` (fg-primary) |
| Усього сьогодні | total requests | `#434652` |

- **Number:** `font-size: 28px; font-weight: 700; line-height: 1.2; display: block`
- **Label:** `font-size: 12px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.06em; color: #6c757d`

---

### Filter Bar

- **Layout:** `display: flex; align-items: center; gap: 12px; flex-wrap: wrap; margin-bottom: 12px`

**Pills (left side):**
- **Default state:** `padding: 7px 14px; border-radius: 0.25rem; background: transparent; border: 1px solid #dee2e6; color: #495057; font-size: 14px; font-weight: 500`
- **Hover:** `background: #F5F5F5; border-color: #ced4da`
- **Active state:** `background: #1A92B0; border-color: #1A92B0; color: #fff`
- **Active hover:** `background: #146f88; border-color: #146f88`
- Two pills: `"Тільки для КПП"` and `"Тільки відкриті"`
- Active state persisted in `localStorage` (`reqType`, `reqStatus`)

**Apartment search (right side, `margin-left: auto`):**
- Grouped input with left icon (search/apartment icon, 14×14), text field, and conditional clear button (✕)
- Container: `display: flex; align-items: center; border: 1px solid #dee2e6; border-radius: 0.25rem; background: #fff; overflow: hidden`
- Focus-within: `border-color: #22C1C3; box-shadow: 0 0 0 3px rgba(34,193,195,.15)`
- Input: `height: 34px; padding: 0 6px; border: none; background: transparent; font-size: 14px; width: 120px; max-length: 4`
- Clear button (shown only when input has value): `height: 34px; padding: 0 10px; border-left: 1px solid #dee2e6; background: transparent; color: #6c757d; font-size: 12px`
- Clears both the input and `localStorage("apt")`

---

### Request Table

**Wrapper:**
- `background: #fff; border-radius: 0.25rem; box-shadow: 0 .5rem 1rem rgba(0,0,0,.15); overflow: hidden`

**Table:**
- `width: 100%; border-collapse: collapse`
- `td { padding: 16px; border-top: 1px solid #e9ecef; vertical-align: middle }`
- First row: `border-top: none`
- Row hover: `td { background: rgba(0,0,0,.025) }`
- Closed row: `td { background: #fafafa; opacity: .75 }` — restored to full opacity on hover

**Each row has two cells:**

**Cell 1 — main content:**

1. **Head row** (`display: flex; align-items: center; gap: 8px; flex-wrap: wrap; margin-bottom: 8px`):
   - **Type badge** — see Badge spec below
   - **Date/time** — `dd.mm.yyyy, HH:MM` format (`uk-UA` locale), `font-size: 12px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.06em; color: #6c757d`
   - **"Закрито" chip** (closed requests only) — `font-size: 12px; font-weight: 600; text-transform: uppercase; background: #e9ecef; padding: 2px 7px; border-radius: 0.2rem; color: #6c757d`

2. **Address row** (`display: flex; align-items: baseline; gap: 8px; flex-wrap: wrap`):
   - Building + apartment number: `font-weight: 700; font-size: 16px`
   - Apartment `#NNN` portion: `color: #1F9991` (brand-teal-600)
   - Resident name + phone: `font-size: 14px; color: #6c757d` — separated by ` · `

3. **Description** (if present): `font-size: 14px; color: #495057; max-width: 72ch; margin-top: 4px; line-height: 1.5`

4. **Photo thumbnails** (if `images.length > 0`): `display: flex; gap: 8px; margin-top: 12px`
   - Each thumb: `52×52px; border-radius: 0.25rem; background: #F5F5F5; border: 1px solid #dee2e6; display: flex; align-items: center; justify-content: center`
   - Icon: camera/image SVG, `20×20px`, `color: #6c757d`
   - Hover: `background: #e8f7f7; border-color: #3fcacb; box-shadow: 0 1px 2px rgba(0,0,0,.06); color: #197872`
   - Click → opens PhotoModal with the full image

**Cell 2 — action button** (`width: 1%; white-space: nowrap; vertical-align: middle; padding-left: 12px`):
- **Закрити** (status = `"new"`): `background: #28a745; color: #fff; padding: 8px 18px; border-radius: 0.25rem; font-size: 14px; font-weight: 600; border: none`
  - Hover: `background: #218838`
  - API call: `PUT /requests/v1/guard/request/{id}` with body `{"status":"closed"}`
- **Відкрити** (status = `"closed"`): same shape, `background: #6c757d`
  - Hover: `background: #5a6268`
  - API call: body `{"status":"new"}`
- After action: `window.location.reload()` (original) or optimistic update + toast

---

### Type Badges

| Type | Ukrainian label | Background color |
|---|---|---|
| `guest` | Гість | `#17A2B8` |
| `taxi` | Таксі | `#F0A500` |
| `delivery` | Доставка | `#4A90A4` |
| `cargo` | 37-Б Розвантаження | `#C0392B` |

Style: `display: inline-block; padding: 3px 9px; border-radius: 0.2rem; font-size: 12px; font-weight: 700; color: #fff; letter-spacing: 0.02em; white-space: nowrap`

Request type ID is numeric in the API: `1=guest, 2=taxi, 3=delivery, 4=cargo`. Type 4 (`cargo`) is КПП-only and filtered out when `reqTypeFilter === "kpp"` is OFF (confusingly, "for КПП" means **include** cargo).

---

### Pager

Shown only when `count > limit`. Bootstrap-style linked page numbers.

- Wrapper: `display: flex; justify-content: center; margin-top: 16px`
- List: `display: flex; border-radius: 0.25rem; overflow: hidden; box-shadow: 0 .125rem .25rem rgba(0,0,0,.075)`
- Each item: `display: block`
- Link: `padding: 9px 16px; background: #fff; border: 1px solid #dee2e6; border-right: none; color: #434652; font-size: 14px; font-weight: 500`
  - Last item: `border-right: 1px solid #dee2e6`
  - Hover: `background: #F5F5F5`
- Active page: `background: #1A92B0; border-color: #1A92B0; color: #fff`

---

### Photo Modal

Triggered when user clicks an image thumbnail.

- Overlay: `position: fixed; inset: 0; background: rgba(0,0,0,.5); display: flex; align-items: center; justify-content: center; z-index: 1000`
- Dialog: `background: #fff; border-radius: 0.25rem; padding: 16px; max-width: 460px; width: 92%; box-shadow: 0 1rem 3rem rgba(0,0,0,.175); position: relative`
- Close button: top-right, `32×32px`, `background: #F5F5F5; border: 1px solid #dee2e6; border-radius: 0.25rem` — contains an × SVG icon
  - Hover: `background: #e9ecef`
- Photo area: `width: 100%; aspect-ratio: 4/3; background: #F5F5F5; border-radius: 0.25rem; margin-top: 8px; overflow: hidden`
  - Real image: `object-fit: contain; width: 100%; height: 100%`
  - Placeholder (no image): centered camera SVG with `"Фото відсутнє"` label in `#6c757d`
- Close on overlay click or close button click

---

### Toast Notification

Shown briefly after Закрити/Відкрити actions.

- `position: fixed; bottom: 24px; left: 50%; transform: translateX(-50%)`
- `background: #212529; color: #fff; padding: 10px 20px; border-radius: 0.25rem; font-size: 14px; font-weight: 500; white-space: nowrap`
- `box-shadow: 0 .5rem 1rem rgba(0,0,0,.15); z-index: 2000`
- Animation: fade + slide up from `translateY(10px)` over `200ms ease`
- Auto-dismiss after **1800ms**
- Messages: `"Заявку закрито"` / `"Заявку відкрито"`

---

### Empty State

When no requests match active filters:

- Single `<td>` spanning full table width
- `padding: 64px 20px; text-align: center; color: #6c757d; font-size: 14px`
- Text: `"Немає заявок за обраними фільтрами."`

---

## Interactions & Behaviour

| Trigger | Behaviour |
|---|---|
| Page load | Reads `localStorage` for `reqStatus`, `reqType`, `apt`; applies saved filters |
| "Тільки для КПП" pill | Toggles `reqTypeFilter` between `"kpp"` and `"all"`; saves to `localStorage("reqType")`; reloads to page 1 |
| "Тільки відкриті" pill | Toggles `reqStatusFilter` between `"new"` and `"all"`; saves to `localStorage("reqStatus")`; reloads to page 1 |
| Apartment input | Filters list as user types (max 4 chars); saves to `localStorage("apt")`; resets to page 1 |
| Clear (✕) button | Clears apartment input and `localStorage("apt")`; resets to page 1 |
| Thumbnail click | Opens PhotoModal with the image src |
| PhotoModal overlay click | Closes modal |
| "Закрити" button | `PUT /requests/v1/guard/request/{id}` `{"status":"closed"}` → `window.location.reload()` |
| "Відкрити" button | `PUT /requests/v1/guard/request/{id}` `{"status":"new"}` → `window.location.reload()` |
| Pager link | Changes active page; re-renders visible rows |
| `<meta http-equiv="refresh" content="60">` | Full page reload every 60 s |

---

## State Management (original jQuery approach)

The original codebase uses jQuery + inline JS. Keep that pattern; no framework needed.

Key `localStorage` keys:
- `reqStatus` — `"new"` | `"all"` (default `"new"`)
- `reqType` — `"kpp"` | `"all"` (default `"kpp"`)
- `apt` — string, max 4 chars (default `""`)

API endpoints (unchanged):
- `GET  /requests/v1/guard/list?offset=&limit=&status=&place=&apartment=`
- `PUT  /requests/v1/guard/request/{id}` — body `{"status":"closed"|"new"}`
- `GET  /dictionary/v1/request-types` — returns `{ data: { 1: {ua,en}, … } }`

---

## Design Tokens

All values are defined as CSS custom properties in `colors_and_type.css` (included in this bundle). Reference them via `var(--*)`.

| Token | Value | Usage |
|---|---|---|
| `--brand-gradient` | `linear-gradient(0deg, #22C1C3 0%, #1F9991 100%)` | Navbar background |
| `--brand-accent` | `#1A92B0` | Active filter pills, active pager |
| `--brand-teal-600` | `#1F9991` | Open-count number, apartment number tint |
| `--brand-teal-100` | `#e8f7f7` | Thumbnail hover background |
| `--fg-primary` | `#434652` | Default body text |
| `--fg-secondary` | `#495057` | Description text |
| `--fg-muted` | `#6c757d` | Meta labels, placeholders |
| `--fg-on-brand` | `#F5F5F5` | Navbar text |
| `--bg-card` | `#ffffff` | Table background, summary strip |
| `--bg-page` | `#F5F5F5` | Page background |
| `--neutral-50` | `#fafafa` | Closed row background |
| `--border-default` | `#dee2e6` | Input borders, table dividers |
| `--shadow-md` | `0 .125rem .25rem rgba(0,0,0,.075)` | Summary strip, pager |
| `--shadow-lg` | `0 .5rem 1rem rgba(0,0,0,.15)` | Request table |
| `--shadow-xl` | `0 1rem 3rem rgba(0,0,0,.175)` | Modal |
| `--shadow-brand` | `0 8px 24px rgba(31,153,145,.25)` | Navbar bottom |
| `--success` | `#28a745` | Закрити button |
| `--secondary` | `#6c757d` | Відкрити button |
| `--rtype-guest` | `#17A2B8` | Guest badge |
| `--rtype-taxi` | `#F0A500` | Taxi badge |
| `--rtype-delivery` | `#4A90A4` | Delivery badge |
| `--rtype-cargo` | `#C0392B` | Cargo badge |
| `--radius-md` | `0.25rem` | All UI elements |
| `--radius-sm` | `0.2rem` | Badges |
| `--font-system` | `"Helvetica Neue", Helvetica, Arial, sans-serif` | All text |
| `--transition-fast` | `120ms ease` | Hover transitions |

---

## Assets

| File | Usage | Notes |
|---|---|---|
| `dynasty/_ui/guard/img/logo.png` | Navbar brand mark, 30×30 | Already in place |
| `dynasty/_ui/guard/img/favicon.ico` | Browser tab | Already in place |

No new image assets are introduced. All icons in the redesign are inline SVGs (no icon font required).

---

## Files in This Bundle

| File | Purpose |
|---|---|
| `index.html` | Main guard console — **design reference** (React prototype) |
| `guard.css` | Design system CSS — extract the relevant rules into the real `dynasty/_ui/guard/` CSS |
| `colors_and_type.css` | All Dynasty Design System CSS custom properties — link or inline into production |
| `NavBar.jsx` | Navbar component — reference for markup and class names |
| `FilterBar.jsx` | Filter pill + search bar component |
| `RequestRow.jsx` | Single table row component |
| `RequestTypeBadge.jsx` | Coloured request-type badge |
| `Pager.jsx` | Pagination component |
| `PhotoModal.jsx` | Image preview modal |
| `data.js` | Mock API data — **not for production** |

---

## Implementation Notes for Claude Code

1. **Target file:** `dynasty/_ui/guard/index.html` — preserve the Go template variables (`{{.PageURI}}`, `{{.APIHost}}`, `{{.PagerLimit}}`), jQuery-based JS, and Bootstrap 4. Do **not** introduce React.

2. **Add a stylesheet:** Create `dynasty/_ui/guard/css/guard-ds.css` that imports or duplicates `colors_and_type.css` tokens, then add the component-level rules from `guard.css` in this bundle.

3. **Replace the inline `<style>` block** in `index.html` with `<link rel="stylesheet" href="/ui/assets/css/guard-ds.css">`.

4. **Rename Bootstrap classes** to their DS equivalents where possible, e.g.:
   - `nav nav-pills nav-fill` → `.guard-filterbar + .guard-pills`
   - `table table-hover shadow` → `.guard-table-wrap > .guard-table`
   - `badge badge-info` → `.guard-badge` with appropriate `--rtype-*` background

5. **Preserve server-side rendering:** The request rows are rendered via `renderItems()` — update the HTML string templates inside that function to output the new class names and markup structure documented above.

6. **Toast:** Add a `<div id="guard-toast" class="guard-toast" hidden>` to the DOM; show/hide it after action button clicks before calling `window.location.reload()`.
