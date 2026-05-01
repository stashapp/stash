# Develop a Frontend Component

Guide for building and modifying React components in the Stash UI.

## When to use
- Creating new UI components
- Modifying existing components
- Adding new pages or features to the frontend

## Project structure

```
ui/v2.5/src/
  components/       # React components organized by feature
    Scenes/         # Scene-related components
    Performers/     # Performer-related components
    Movies/         # Movie/group components
    Studios/        # Studio components
    Tags/           # Tag components
    Galleries/      # Gallery components
    Images/         # Image components
    Settings/       # Settings pages
    Shared/         # Shared/common components
    List/           # List/filter components
  hooks/            # Custom React hooks
  core/             # Generated GraphQL hooks (from codegen)
  models/           # TypeScript type definitions
  utils/            # Utility functions
  locales/          # i18n translation files
  styles/           # Shared SCSS styles
```

## Steps

### 1. If adding new GraphQL data, update the schema first

See the "Add a new GraphQL field" skill, then run `make generate`.

### 2. Create or modify the component

Components typically follow this structure:
```tsx
// ui/v2.5/src/components/Scenes/MyNewComponent.tsx
import React from "react";
import { useScenesQuery } from "src/core/StashService";

export const MyNewComponent: React.FC = () => {
  const { data, loading } = useScenesQuery();
  // ...
};
```

### 3. Use shared components and patterns

- Use components from `Shared/` for consistent UI
- Follow existing patterns in similar feature components
- Use Bootstrap 4 classes for layout (project uses Bootstrap)
- Use FontAwesome icons via `@fortawesome/react-fontawesome`

### 4. Add i18n strings

If adding user-visible text:
```tsx
import { useIntl } from "react-intl";
const intl = useIntl();
const label = intl.formatMessage({ id: "component.my_label" });
```

Add the key to locale files in `src/locales/`.

### 5. Format and lint

```bash
cd ui/v2.5
pnpm run format       # Format with Prettier
pnpm run lint         # Run ESLint + Stylelint
pnpm run check        # TypeScript type check
```

Or from project root:
```bash
make fmt-ui           # Format UI code
make validate-ui      # Full UI validation
```

### 6. Test in the browser

```bash
# Terminal 1: Start backend
make server-start

# Terminal 2: Start frontend dev server
make ui-start
```

Open `http://localhost:3000/` and verify your changes.

## Tech stack
- React with TypeScript
- Apollo Client for GraphQL
- Bootstrap 4 for layout
- FontAwesome for icons
- react-intl for i18n
- Vite for dev server and bundling
- pnpm for package management (NOT npm or yarn)
