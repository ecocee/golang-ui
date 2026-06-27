# UI Context

## Theme (shadcn/ui style)

The visual language is a **modern, clean design system** inspired by shadcn/ui. 
Because we use the Glyra Architecture, we can achieve perfect 1:1 parity with shadcn/ui using raw CSS.

Design characteristics:
- **Near-white / near-black base** — clean backgrounds, no muddy grays (Zinc palette).
- **Subtle borders** — Borders define structure without shouting.
- **CSS Variables** — All colors are mapped to CSS variables in `:root`.

## CSS Design Tokens (`frontend/css/theme.css`)

We use the Zinc palette. The tokens are defined in CSS variables and consumed by component CSS.

### Light Mode Variables
```css
:root {
  --background: 0 0% 100%;
  --foreground: 240 10% 3.9%;
  --card: 0 0% 100%;
  --card-foreground: 240 10% 3.9%;
  --popover: 0 0% 100%;
  --popover-foreground: 240 10% 3.9%;
  --primary: 240 5.9% 10%;
  --primary-foreground: 0 0% 98%;
  --secondary: 240 4.8% 95.9%;
  --secondary-foreground: 240 5.9% 10%;
  --muted: 240 4.8% 95.9%;
  --muted-foreground: 240 3.8% 46.1%;
  --accent: 240 4.8% 95.9%;
  --accent-foreground: 240 5.9% 10%;
  --destructive: 0 84.2% 60.2%;
  --destructive-foreground: 0 0% 98%;
  --border: 240 5.9% 90%;
  --input: 240 5.9% 90%;
  --ring: 240 5.9% 10%;
  --radius: 0.5rem;
}
```

## Component Architecture

Components are built using Semantic HTML and Component-wise CSS.

### Example: Button (`frontend/css/components/button.css`)
```css
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: calc(var(--radius) - 2px);
  font-size: 0.875rem;
  font-weight: 500;
  height: 2.25rem;
  padding: 0 1rem;
  transition: colors 0.15s;
  cursor: pointer;
  border: none;
}
.btn-primary {
  background-color: hsl(var(--primary));
  color: hsl(var(--primary-foreground));
}
.btn-primary:hover { opacity: 0.9; }
```

## Layout Patterns

- We use native **CSS Flexbox** and **CSS Grid**. 
- No custom Go layout engine is needed anymore, the browser engine handles it natively.
- Use `gap`, `justify-content`, and `align-items` for precise spacing.
