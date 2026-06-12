# Frontend Guidelines

## Colors

Use the Tailwind color names from `src/app/globals.css`, not raw hex values.

| Class                          | Use for                                |
| ------------------------------ | -------------------------------------- |
| `bg-brand-bg`                  | Page background, navbar, hero sections |
| `bg-brand-surface`             | Cards, panels, sidebars                |
| `bg-brand-primary`             | Main buttons, active states            |
| `hover:bg-brand-primary-hover` | Primary button hover                   |
| `text-brand-text`              | Main text                              |
| `text-brand-text-muted`        | Paragraphs, helper text                |
| `text-brand-text-subtle`       | Labels, metadata                       |
| `border-brand-border`          | Borders and dividers                   |

## API Calls

Use `src/lib/client.ts` for frontend API calls.

```tsx
import apiClient from "@/lib/client";

const { data } = await apiClient.get("/memberships/me");
```

## Notifications

Use `sonner` for toast notifications. The themed toaster is mounted in the root layout, so components only need to import and call `toast`.

```tsx
import { toast } from "sonner";

toast.success("Membership updated");
toast.error("Unable to save changes");
```

## Components

Put shared components in `src/components`. Group similar components into the same folder and add folders as needed.
