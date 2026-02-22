# AGENTS.md - Website Development Guidelines

Documentation website for `gh wt` GitHub CLI extension. Built with Astro and MDX, hosted on GitHub Pages.

## Tech Stack

- **Framework**: Astro 5.x
- **Content**: MDX for CLI documentation
- **Hosting**: GitHub Pages
- **Styling**: Custom CSS with CSS variables

## Development Commands

The project uses [Task](https://taskfile.dev) for running development commands. Make sure you have Task installed.

### Running the Website

```bash
task website:dev      # Start development server (interactive)
task website:build    # Build for production
task website:preview  # Preview built website locally
```

### Documentation Generation

```bash
task website:docs:generate  # Generate CLI docs from command help
```

### Testing Changes

Use the `agent-browser` skill to test website changes in a real browser:

```bash
use_skill with skill_name: "agent-browser"
```

After activating the skill, you can navigate to the development server, verify pages load correctly, test interactions, and take screenshots to confirm visual changes work as expected.

### Astro Commands

If needed, you can also run Astro commands directly:

```bash
cd website
npm run astro -- --help  # View all Astro commands
```

## Project Structure

```
website/
├── src/
│   ├── components/     # Astro components
│   │   ├── CommandCard.astro
│   │   ├── NeonText.astro
│   │   ├── Scanlines.astro
│   │   ├── Sidebar.astro
│   │   ├── TerminalDemo.astro
│   │   └── TerminalHero.astro
│   ├── layouts/        # Page layouts
│   │   ├── CLILayout.astro
│   │   └── DocsLayout.astro
│   ├── pages/          # Route pages
│   │   ├── index.astro
│   │   └── docs/
│   │       ├── *.md
│   │       └── cli/*.mdx
│   └── styles/
│       └── global.css
├── public/             # Static assets
├── scripts/            # Build scripts
│   └── generate-cli-docs.js
├── astro.config.mjs
└── package.json
```

## Design Guidelines

### Aesthetic Direction

Retro-futuristic/cyberpunk terminal aesthetic - CRT monitors, neon glows, scanlines, typing animations.

### Color Palette

| Role | Color | Hex |
|------|-------|-----|
| Background | Deep charcoal | #0a0a0f |
| Primary accent | Electric cyan | #00fff9 |
| Secondary accent | Hot magenta | #ff00ff |
| Tertiary | Amber | #ffb800 |
| Success | Bright green | #00ff88 |
| Error | Coral red | #ff3366 |

### Visual Effects

- CRT scanline overlay (subtle, CSS pseudo-element)
- Neon glow on hover (box-shadow, text-shadow)
- Typing cursor animations
- Grid background with flicker effect
- Staggered reveal animations on page load

### Landing Page Elements

- Animated terminal demos showing CLI usage
- Install command prominently displayed

## CLI Documentation

CLI documentation is auto-generated from the Go source code using the `docs:generate` script. The generated MDX files are placed in `src/pages/docs/cli/`.

### Updating CLI Docs

After modifying CLI commands in Go code:

```bash
task website:docs:generate
```

This will regenerate all CLI documentation pages from command help output.

## Adding New Documentation

### Markdown Pages (.md)

Add markdown files to `src/pages/docs/`:

```markdown
---
layout: ../../layouts/DocsLayout.astro
title: "Page Title"
---

# Content
```

### CLI Command Pages (.mdx)

CLI command pages are auto-generated. To add a new command:

1. Add the command in `cmd/` package
2. Run `task website:docs:generate` to regenerate docs

## CSS Guidelines

- Use CSS custom properties (variables) for colors defined in the palette
- Use semantic class names
- Keep global styles in `src/styles/global.css`
- Component-specific styles can be in the component file

## Commit Messages

Use Conventional Commits format:

```
<type>(<scope>): <description>

[optional body]
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

Example: `docs(website): add installation guide page`
