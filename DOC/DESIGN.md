---
version: "alpha"
name: "[Project Name]"
description: "[One-line visual identity summary — e.g.: Bold dark-mode dashboard with electric-blue accents and sharp geometric rhythm]"

colors:
  primary:          "[#XXXXXX]"   # Main interactive color — buttons, links, primary actions
  primary-dark:     "[#XXXXXX]"   # Hover / pressed state of primary
  secondary:        "[#XXXXXX]"   # Supporting accents, badges, tags
  accent:           "[#XXXXXX]"   # Highlight, focus ring, notification indicators
  background:       "[#XXXXXX]"   # Page/app background
  surface:          "[#XXXXXX]"   # Cards, panels, raised surfaces
  surface-alt:      "[#XXXXXX]"   # Alternate surface (sidebars, modals)
  on-primary:       "[#XXXXXX]"   # Text/icon color over primary
  on-surface:       "[#XXXXXX]"   # Default body text color
  on-surface-muted: "[#XXXXXX]"   # Captions, placeholders, metadata
  border:           "[#XXXXXX]"   # Default border/divider color
  error:            "[#XXXXXX]"   # Validation errors, destructive actions
  success:          "[#XXXXXX]"   # Confirmations, success badges
  warning:          "[#XXXXXX]"   # Warnings, alerts

typography:
  h1:
    fontFamily: "[Font Name]"
    fontSize:   "2.5rem"
    fontWeight: "700"
    lineHeight: "1.15"
  h2:
    fontFamily: "[Font Name]"
    fontSize:   "2rem"
    fontWeight: "600"
    lineHeight: "1.2"
  h3:
    fontFamily: "[Font Name]"
    fontSize:   "1.5rem"
    fontWeight: "600"
    lineHeight: "1.25"
  body-md:
    fontFamily: "[Font Name]"
    fontSize:   "1rem"
    fontWeight: "400"
    lineHeight: "1.6"
  body-sm:
    fontFamily: "[Font Name]"
    fontSize:   "0.875rem"
    fontWeight: "400"
    lineHeight: "1.5"
  label:
    fontFamily: "[Font Name]"
    fontSize:   "0.875rem"
    fontWeight: "500"
    letterSpacing: "0.02em"
  label-caps:
    fontFamily: "[Font Name]"
    fontSize:   "0.75rem"
    fontWeight: "600"
    letterSpacing: "0.08em"
  code:
    fontFamily: "[Monospace Font]"
    fontSize:   "0.875rem"
    fontWeight: "400"

rounded:
  none: "0"
  sm:   "4px"
  md:   "8px"
  lg:   "16px"
  xl:   "24px"
  full: "9999px"

spacing:
  xs:  "4px"
  sm:  "8px"
  md:  "16px"
  lg:  "24px"
  xl:  "48px"
  2xl: "80px"

components:
  button-primary:
    backgroundColor: "{colors.primary}"
    textColor:       "{colors.on-primary}"
    rounded:         "{rounded.md}"
    padding:         "10px 20px"
    typography:      "{typography.label}"
  button-primary-hover:
    backgroundColor: "{colors.primary-dark}"
  button-secondary:
    backgroundColor: "transparent"
    textColor:       "{colors.primary}"
    rounded:         "{rounded.md}"
    padding:         "10px 20px"
  card:
    backgroundColor: "{colors.surface}"
    rounded:         "{rounded.lg}"
    padding:         "{spacing.lg}"
  input:
    backgroundColor: "{colors.surface}"
    textColor:       "{colors.on-surface}"
    rounded:         "{rounded.md}"
    padding:         "10px 14px"
  badge:
    backgroundColor: "{colors.secondary}"
    textColor:       "{colors.on-primary}"
    rounded:         "{rounded.full}"
    padding:         "2px 10px"
---

# DESIGN.md — Visual Identity

> **GSD Rule — no inspiration references here.**
> This file contains only the **final, decided design values** the agent must implement.
> References to products, apps, or visual sources that inspired the design belong in `PROJECT.md`.
> If the user says *"make it look like Product X"*, extract the relevant values (colors, fonts, spacing)
> and record them here as concrete tokens. Never mention Product X in this file.

> **Immutability Rule — append only.**
> Do not overwrite existing values. When a decision changes, add a new commented entry with a
> timestamp below the old one. Format: `# [YYYY-MM-DD HH:MM] — superseded by: [new value]`

> **Agent protocol — how to collect design preferences (run once per project):**
> Before filling this file, ask the developer:
> 1. "What is the overall visual mood? (e.g.: dark/light, minimal/rich, serious/playful)"
> 2. "Are there any colors you want as a starting point? (brand colors, hex codes)"
> 3. "Any font preferences? (type 'system default' if no preference)"
> 4. "Any apps, sites, or UI patterns you like the feel of? (I will extract values — not copy them)"
> 5. "Any animations or motion preferences? (subtle/none, energetic, smooth/bouncy)"
> 6. "Any layout preferences? (card-heavy, list-based, sidebar navigation, etc.)"
> Record answers in `PROJECT.md` if inspiration sources are mentioned, then translate to tokens here.

---

## Overview

> *High-level description of the visual identity — mood, metaphor, and design philosophy.*
> *Written once; updated only when the identity significantly changes (with timestamp).*

[Describe the visual personality of this project. E.g.: "Clean, professional dashboard aesthetic
with a dark base and electric-blue primary accents. Typography is sharp and utilitarian.
Interactions are smooth but not showy — micro-animations under 200ms only."]

**Last updated:** [YYYY-MM-DD HH:MM] — [Responsible]

---

## Colors

> *Rationale for each color decision. Reference the token name from YAML front matter.*

- **`primary`** — [Explain the choice and its emotional/functional role]
- **`secondary`** — [Explain]
- **`accent`** — [Explain]
- **`background` / `surface`** — [Explain the depth layering logic]
- **`error` / `success` / `warning`** — [Explain semantic color choices]

**Dark mode:** [Yes / No / Planned — describe any token overrides for dark mode if applicable]

**Contrast compliance:** All text-on-background combinations must meet WCAG AA (4.5:1 for body, 3:1 for large text).
Validate with: `npx @google/design.md lint DESIGN.md`

---

## Typography

> *Font choices, their loading method, and usage rules.*

- **Primary font (`[Font Name]`):** [Why this font — personality, legibility, licensing]
  - Loading: [Google Fonts / self-hosted / system stack] — `@import url('[URL]')`
- **Monospace font (`[Font Name]`):** [For code blocks, terminal-style labels]
- **Fallback stack:** `[Primary], [Category fallback], sans-serif`

**Scale rationale:** [Describe the type scale — modular scale? fixed scale? visual hierarchy goal?]

---

## Motion & Animations

> *Rules for micro-animations, transitions, and motion behavior.*

| Element | Duration | Easing | Notes |
| :--- | :--- | :--- | :--- |
| Page transitions | [e.g.: 200ms] | [e.g.: ease-out] | [e.g.: fade + slight translate-Y] |
| Button hover | [e.g.: 120ms] | [e.g.: ease-in-out] | [e.g.: scale(1.02) + shadow] |
| Modal open/close | [e.g.: 250ms] | [e.g.: cubic-bezier(0.16,1,0.3,1)] | [e.g.: slide-up + fade] |
| Sidebar expand | [e.g.: 300ms] | [e.g.: ease-in-out] | |
| Toast / notification | [e.g.: 350ms] | [e.g.: spring] | [e.g.: slide-in from right] |
| Loading skeleton | [e.g.: 1.5s loop] | [e.g.: linear] | [e.g.: shimmer effect] |

**Motion philosophy:** [Subtle and purposeful / Energetic and expressive / Minimal — content first]
**Reduced motion:** All animations must be disabled or minimized when `prefers-reduced-motion: reduce`.

---

## Layout & Spacing

> *Grid system, spacing philosophy, and responsive behavior.*

- **Grid:** [e.g.: 12-column, 24px gutter / CSS Grid with named areas]
- **Max content width:** [e.g.: 1280px]
- **Sidebar width (if applicable):** [e.g.: 240px collapsed → 280px expanded]
- **Breakpoints:**

| Name | Min-width | Description |
| :--- | :--- | :--- |
| mobile | 0px | Single column, bottom navigation |
| tablet | 768px | Two-column, collapsible sidebar |
| desktop | 1024px | Full layout, persistent sidebar |
| wide | 1440px | Max-width container centered |

**Spacing philosophy:** [e.g.: 8px base unit — all spacing is a multiple of 8]

---

## Component Patterns

> *Specific UI component rules beyond what's in the token schema.*

### Navigation
- [e.g.: Top bar on mobile, left sidebar on desktop]
- [Active state: primary color left border + background tint]

### Forms
- [e.g.: Labels above inputs, not inline]
- [Error messages below inputs in `error` color]
- [Focus ring: 2px solid `{colors.accent}` with 2px offset]

### Cards
- [e.g.: Flat cards (no shadow by default), hover adds shadow-md]
- [Clickable cards: add `cursor: pointer` + scale(1.01) on hover]

### Tables
- [e.g.: Zebra striping / dividers only / full borders]
- [Sticky header on scroll]

### Feedback States
- [e.g.: Empty states with illustration + CTA]
- [Loading: skeleton screens, not spinners]
- [Errors: inline (inputs), toast (server), full-page (critical)]

---

## Aesthetic Notes

> *Qualitative guidelines that don't reduce to tokens but guide implementation decisions.*

[E.g.: "Prefer iconography over text labels in dense UIs. Use line icons, not filled.
Avoid drop shadows on flat-design surfaces — rely on surface color layering instead.
Photography (if used) must be high-contrast, desaturated to fit the palette."]

---

## Version History (Immutable)

| Timestamp | Author | Change |
| :--- | :--- | :--- |
| [YYYY-MM-DD HH:MM] | [Responsible] | Initial DESIGN.md created. All values are placeholders — to be filled with project-specific decisions. |

---

*Validate tokens: `npx @google/design.md lint DESIGN.md`*
*Compare versions: `npx @google/design.md diff DESIGN.md DESIGN-v2.md`*
