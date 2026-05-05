---
version: "alpha"
name: "PF_ai"
description: "Clean technical workbench with a light neutral base, sharp hierarchy, calm green accents, and dense operational layouts."

colors:
  primary:          "#1F7A5C"
  primary-dark:     "#155C45"
  secondary:        "#2F5D9F"
  accent:           "#D97706"
  background:       "#F7F8F6"
  surface:          "#FFFFFF"
  surface-alt:      "#EEF1ED"
  on-primary:       "#FFFFFF"
  on-surface:       "#1D2320"
  on-surface-muted: "#5E6B64"
  border:           "#D6DDD8"
  error:            "#B42318"
  success:          "#1F7A5C"
  warning:          "#B54708"

typography:
  h1:
    fontFamily: "Inter"
    fontSize:   "2.25rem"
    fontWeight: "700"
    lineHeight: "1.15"
  h2:
    fontFamily: "Inter"
    fontSize:   "1.75rem"
    fontWeight: "650"
    lineHeight: "1.2"
  h3:
    fontFamily: "Inter"
    fontSize:   "1.25rem"
    fontWeight: "650"
    lineHeight: "1.25"
  body-md:
    fontFamily: "Inter"
    fontSize:   "1rem"
    fontWeight: "400"
    lineHeight: "1.55"
  body-sm:
    fontFamily: "Inter"
    fontSize:   "0.875rem"
    fontWeight: "400"
    lineHeight: "1.45"
  label:
    fontFamily: "Inter"
    fontSize:   "0.875rem"
    fontWeight: "600"
    letterSpacing: "0"
  code:
    fontFamily: "JetBrains Mono"
    fontSize:   "0.875rem"
    fontWeight: "400"

rounded:
  none: "0"
  sm:   "4px"
  md:   "8px"
  lg:   "8px"
  full: "9999px"

spacing:
  xs:  "4px"
  sm:  "8px"
  md:  "16px"
  lg:  "24px"
  xl:  "40px"
  2xl: "64px"

components:
  button-primary:
    backgroundColor: "{colors.primary}"
    textColor:       "{colors.on-primary}"
    rounded:         "{rounded.md}"
    padding:         "9px 14px"
    typography:      "{typography.label}"
  button-secondary:
    backgroundColor: "{colors.surface}"
    textColor:       "{colors.on-surface}"
    rounded:         "{rounded.md}"
    padding:         "9px 14px"
  panel:
    backgroundColor: "{colors.surface}"
    rounded:         "{rounded.md}"
    padding:         "{spacing.md}"
  input:
    backgroundColor: "{colors.surface}"
    textColor:       "{colors.on-surface}"
    rounded:         "{rounded.md}"
    padding:         "9px 12px"
  badge:
    backgroundColor: "{colors.surface-alt}"
    textColor:       "{colors.on-surface}"
    rounded:         "{rounded.full}"
    padding:         "2px 8px"
---

# DESIGN.md - Visual Identity

## Overview
PF_ai should feel like a precise technical workbench: calm, readable, structured, and built for repeated operational use. The UI prioritizes scanning, comparison, and control over marketing-style presentation.

**Last updated:** 2026-05-05 15:06 - [CEO]

## Colors
- **`primary`** - Green communicates active operation, successful orchestration, and system readiness.
- **`secondary`** - Blue supports informational states without turning the palette into a single-hue UI.
- **`accent`** - Amber marks attention, pending decisions, and user intervention.
- **`background` / `surface`** - Light neutral layers keep the interface calm and document-like.
- **`error` / `success` / `warning`** - Semantic colors are reserved for system states and validation.

**Dark mode:** Planned after MVP.
**Contrast compliance:** All text-on-background combinations must meet WCAG AA.

## Typography
- **Primary font (`Inter`):** Neutral, highly legible, suited for dashboards and tools.
- **Monospace font (`JetBrains Mono`):** Used for handoff JSON, paths, model names, and logs.
- **Fallback stack:** `Inter, system-ui, sans-serif`

## Motion & Animations

| Element | Duration | Easing | Notes |
| :--- | :--- | :--- | :--- |
| Page transitions | 160ms | ease-out | Fade only |
| Button hover | 120ms | ease-in-out | Color and border change |
| Modal open/close | 180ms | ease-out | Fade plus small vertical movement |
| Sidebar expand | 180ms | ease-in-out | Width transition |
| Toast | 220ms | ease-out | Slide from top-right |
| Loading skeleton | 1.2s loop | linear | Subtle shimmer |

**Motion philosophy:** Subtle and purposeful.
**Reduced motion:** Disable non-essential animation when `prefers-reduced-motion: reduce`.

## Layout & Spacing
- **Grid:** CSS Grid and flex layouts with 8px spacing base.
- **Max content width:** 1440px for operational screens.
- **Sidebar width:** 248px expanded; 64px collapsed.
- **Breakpoints:**

| Name | Min-width | Description |
| :--- | :--- | :--- |
| mobile | 0px | Single column, compact navigation |
| tablet | 768px | Two-column layouts |
| desktop | 1024px | Persistent sidebar and dense panels |
| wide | 1440px | Full workbench layout |

## Component Patterns

### Navigation
- Persistent left sidebar on desktop.
- Compact top or bottom navigation on mobile.
- Active state uses primary left border and subtle surface-alt background.

### Forms
- Labels above inputs.
- Inline validation below fields.
- Focus ring uses `accent` with 2px offset.

### Panels
- Use panels for tools, registries, logs, and repeated items.
- Do not nest panels.
- Keep radius at 8px or less.

### Tables
- Dividers and sticky headers for agent/model registries.
- Prefer dense rows with clear status badges.

### Feedback States
- Loading uses skeletons.
- Empty states should expose the next action directly.
- Errors must be specific and actionable.

## Aesthetic Notes
- Prefer icons plus concise labels for commands.
- Avoid decorative gradients, oversized heroes, and marketing composition.
- Use real operational data as the main visual signal.
- Keep text compact and scannable.

## Version History (Immutable)

| Timestamp | Author | Change |
| :--- | :--- | :--- |
| [2026-05-05 15:06] | [CEO] | Initial PF_ai design tokens filled during onboarding. |
