# AI Instructions for @spilnu/backoffice-core and Related Projects

This document provides coding guidelines and preferences for AI assistants working on the Spilnu backoffice applications.

## Project Context

- **Framework**: Nuxt 4 with Vue 3 and TypeScript
- **Package Type**: Layered Nuxt application with core and region-specific layers
- **Target Platforms**: Desktop web applications (backoffice/admin tools)
- **Primary Markets**: Danish (DK) and English (UK) gaming/casino backoffice applications
- **Styling**: Tailwind CSS v4 with Nuxt UI v4 (alpha)
- **UI Framework**: @nuxt/ui for components (Dashboard, Cards, Buttons, etc.)
- **State Management**: @tanstack/vue-query via @peterbud/nuxt-query
- **Authentication**: Keycloak integration for SSO
- **Real-time**: XMPP integration for chat functionality
- **SSR**: Disabled (client-side only application)
- **Testing**: @nuxt/test-utils, Vitest for unit tests. @nuxt/test-utils, Vitest and optionally Playwright (with @nuxt/test-utils integration) for integration/e2e tests.

## AI Assistant Tools & Resources

### Nuxt MCP Server

**Always use the Nuxt MCP Server when available** for accurate, up-to-date Nuxt documentation:

- **MCP Server**: `@nuxt` (antfu/nuxt-mcp)
- **Endpoint**: `https://mcp.nuxt.com/sse`
- **Setup**: Configure in VS Code via `.vscode/settings.json` or user MCP settings
- **Capabilities**:
    - Search official Nuxt documentation (Nuxt Core, Nuxt UI, Nuxt Content, NuxtHub)
    - List and discover Nuxt modules with stats and compatibility info
    - Get accurate API references and best practices

**When to use:**
- Questions about Nuxt features, APIs, or configuration
- Looking up best practices for Nuxt-specific functionality
- Checking Nuxt module compatibility and installation
- Understanding Nuxt UI components usage
- NuxtHub database, blob, or KV store implementation

**Example queries:**
```
@nuxt How do I configure a custom Nuxt module?
@nuxt What are the best practices for SEO meta tags in Nuxt?
@nuxt Show me Nuxt UI Button component props
@nuxt How to use useAsyncData with error handling?
```

**Important**: Always prefer MCP-sourced documentation over assumptions or outdated information. If the MCP server is not available, clearly indicate you're working without live documentation access.

## Code Style & Formatting

### Linting & Code Quality

- Use `@cego/eslint-config-nuxt` for ESLint configuration
- Run `npm run lint:fix` before committing
- Follow TypeScript strict mode conventions
- Enable `noUncheckedIndexedAccess` in TypeScript config

### Indentation & Spacing

**Standard:**
- 4 spaces for TypeScript/JavaScript/Vue/SCSS/CSS files
- No trailing whitespace
- End files with a single newline

### Line Length

- No strict limit, but prioritize readability
- Break long lines sensibly at logical points
- Keep function signatures readable

## Vue Component Structure

### File Naming

- Use PascalCase for component files: `DefaultLayout.vue`, `ChatFilters.vue`, `Logo.vue`
- Group related components in directories: `Chat/`, `Emoji/`, `Base/`
- Use kebab-case for composable files: `useAuth.ts`, `useApi.ts`, `useFormat.ts`
- Organize composables by feature: `chat/`, `brand/`, `player-account/`

### Component Order

**Standard structure:**
```vue
<script setup lang="ts">
// 1. Imports (grouped by source)
// 2. Type/Interface definitions (if small and component-specific)
// 3. Props definition with `defineProps` or `withDefaults`
// 4. Emits definition with `defineEmits`
// 5. Composables
// 6. Reactive state
// 7. Computed properties
// 8. Functions
// 9. Lifecycle hooks
// 10. Provide/Inject
// 11. defineExpose
</script>

<template>
    <!-- Component template with Tailwind classes -->
</template>

<style lang="scss" scoped>
// Optional: Only use for complex custom styles that can't be achieved with Tailwind
// Prefer Tailwind utility classes in the template
</style>
```

### Script Setup Syntax

**Always use:**
- `<script setup lang="ts">` (NOT Options API)
- TypeScript for all script blocks
- Explicit type annotations for props

**Example:**
```vue
<script setup lang="ts">
export interface ComponentProps {
    variant?: 'primary' | 'secondary' | 'danger'
    disabled?: boolean
    loading?: boolean
}

const props = withDefaults(defineProps<ComponentProps>(), {
    variant: 'primary',
    disabled: false,
    loading: false
})
</script>
```

**Note:** For simple components without props or with simple logic, you can omit the props interface entirely and use inline type definitions.

### Props Definition
- Use TypeScript interface + `defineProps<T>()`
- Use `withDefaults(defineProps<ComponentProps>() { ... })` when setting defaults is needed

**Example**
```typescript
export interface ComponentProps {
    title?: string
    description?: string
    showIcon?: boolean
    variant?: 'default' | 'compact'
}

const props = defineProps<ComponentProps>()
```

## TypeScript Conventions

### Type Definitions

- Store shared types in `src/runtime/types/`
- Use type imports: `import type { ... } from '...'`
- Prefer interfaces over types for object shapes
- Define component-specific interfaces and extract to separate type file in `types/components/` if used by multiple components

### Type Naming

- Prefix component prop interfaces with component name: `ChatFiltersProps`, `LogoProps`
- Use descriptive names: `UseQueryReturn`, `ApiResponse`, `ServiceOptions`
- Suffix return types with `Return`: `UseFormatReturn`, `UseBrandsReturn`
- Suffix options with `Options`: `UseFormatAsCurrencyOptions`, `UseFormatAsDateOptions`

### Imports

**Import order:**
1. Vue core imports
2. Nuxt imports (`#imports`, `#app`, `#components`)
3. Nuxt UI types (`@nuxt/ui`)
4. Third-party packages
5. Layer imports (`#layers/core/...`)
6. Local type imports
7. Relative imports

**Use auto-imports:**
- Nuxt auto-imports components, composables and Vue.js APIs to use across the application
- Nuxt utilities are auto-imported
- Vue utilities are auto-imported (ref, computed, watch, etc.)
- Types from `core/types/` are auto-imported

**Example:**
```typescript
import { ref, computed } from '#imports'
import type { TabsItem, NavigationMenuItem } from '@nuxt/ui'
import { useChatMessages } from '#layers/core/app/composables/chat/useChatMessages'
```

### Utility Libraries

**Always prefer existing utilities from established libraries over writing custom implementations:**

**@vueuse/core** - For Vue-specific utilities and composables:
- Use `@vueuse/core` for common reactive utilities: `useToggle`, `useDebounceFn`, `useThrottleFn`, `useEventListener`, etc.
- Prefer VueUse composables for DOM interactions, localStorage, media queries, etc.
- Examples: `useStorage`, `useMediaQuery`, `useIntersectionObserver`, `useElementSize`

**es-toolkit** - For general JavaScript utilities:
- Use `es-toolkit` for array/object manipulation: `pick`, `omit`, `groupBy`, `chunk`, `debounce`, `throttle`
- Prefer es-toolkit over lodash or custom utilities for common operations
- Examples from codebase: `pick(props, allowedKeys)`, `groupBy(items, 'category')`

**Why:**
- Battle-tested, optimized implementations
- Smaller bundle sizes (especially es-toolkit)
- Type-safe TypeScript definitions
- Maintained by the community
- Reduces custom code maintenance burden

**When to write custom utilities:**
- Business logic specific to Spilnu/gaming domain
- Complex brand-specific transformations
- When no suitable library function exists

**Example - Prefer library utilities:**
```typescript
// ❌ Avoid writing custom implementations
const picked = Object.keys(props)
    .filter(key => allowedKeys.includes(key))
    .reduce((obj, key) => ({ ...obj, [key]: props[key] }), {})

// ✅ Use es-toolkit
import { pick } from 'es-toolkit'
const picked = pick(props, allowedKeys)

// ❌ Avoid custom debounce
let timeout: NodeJS.Timeout
const handleSearch = (value: string) => {
    clearTimeout(timeout)
    timeout = setTimeout(() => search(value), 300)
}

// ✅ Use @vueuse/core
import { useDebounceFn } from '@vueuse/core'
const handleSearch = useDebounceFn((value: string) => search(value), 300)
```

## Tailwind CSS & Nuxt UI Conventions

### Styling Approach

**Prefer Tailwind utility classes** for most styling needs:

```vue
<template>
    <div class="flex flex-col gap-4 p-6">
        <h2 class="text-xl font-bold text-primary">
            Title
        </h2>
        <p class="text-sm text-muted">
            Description text
        </p>
    </div>
</template>
```

### Nuxt UI Components

This project uses **@nuxt/ui v4** which provides pre-built components. Always use Nuxt UI components when available:

**Common components:**
- `UApp` - App wrapper with theme support
- `UCard` - Card container with header/footer slots
- `UButton` - Button with variants and icons
- `UInput` - Form input with validation
- `UTabs` - Tab navigation with slots
- `UAvatar` - User avatar
- `UUser` - User info with avatar and description
- `UDashboardGroup` - Dashboard layout container
- `UDashboardSidebar` - Collapsible sidebar with header/footer
- `UDashboardSearch` - Command palette search
- `UNavigationMenu` - Navigation menu with icons

**Example usage:**
```vue
<template>
    <UCard>
        <template #header>
            <div class="flex items-center justify-between">
                <span class="text-primary">Card Title</span>
                <UButton variant="outline" size="sm" icon="i-heroicons-funnel">
                    Filters
                </UButton>
            </div>
        </template>

        <!-- Card content -->
        <div class="space-y-4">
            <p>Content goes here</p>
        </div>

        <template #footer>
            <UInput 
                v-model="input" 
                placeholder="Type something..."
                @keyup.enter="handleSubmit"
            />
        </template>
    </UCard>
</template>
```

### Nuxt UI Configuration

Configure Nuxt UI theme in `app.config.ts`:

```typescript
export default defineAppConfig({
    ui: {
        colors: {
            primary: 'violet',
            secondary: 'orange',
            neutral: 'zinc',
        }
    }
})
```

### Responsive Design

**This is a desktop-first application**, but use Tailwind responsive prefixes when needed:

```vue
<template>
    <!-- Default (mobile) to larger screens -->
    <div class="p-4 md:p-6 lg:p-8">
        <h1 class="text-xl md:text-2xl lg:text-3xl">
            Responsive Heading
        </h1>
    </div>
</template>
```

**Standard breakpoints:**
- `sm`: 640px
- `md`: 768px
- `lg`: 1024px
- `xl`: 1280px
- `2xl`: 1536px

### Custom Styles

**Only use `<style>` blocks when:**
- Complex animations that can't be done with Tailwind
- Deep customization of third-party components
- Styles that need CSS features not available in Tailwind

**Example:**
```vue
<template>
    <div class="custom-component">
        <p class="text-primary">Most styling via Tailwind</p>
    </div>
</template>

<style scoped>
/* Only for complex custom needs */
.custom-component {
    /* Complex animation or layout */
    animation: slideIn 0.3s ease-out;
}

@keyframes slideIn {
    from { transform: translateX(-100%); }
    to { transform: translateX(0); }
}
</style>
```

### Theming & Design Tokens

Nuxt UI provides design tokens that work with Tailwind:

**Text colors:**
- `text-primary` - Primary theme color
- `text-secondary` - Secondary theme color
- `text-muted` - Muted/secondary text
- `text-default` - Default text color

**Background colors:**
- `bg-elevated` - Elevated surface
- `bg-default` - Default background

**Borders:**
- `border-default` - Default border color
- `ring-default` - Focus ring color

### Icons

Nuxt UI supports multiple icon libraries. Use icon prefixes in component props:

**Common icon libraries:**
- `i-heroicons-*` - Heroicons
- `i-lucide-*` - Lucide icons
- `fluent:*` - Fluent icons

**Example usage:**
```vue
<template>
    <UButton 
        icon="i-heroicons-funnel"
        @click="toggleFilters"
    >
        Filters
    </UButton>

    <UNavigationMenu
        :items="[
            { label: 'Home', icon: 'lucide:house', to: '/' },
            { label: 'Marketing', icon: 'fluent:megaphone-28-filled', to: '/campaigns' }
        ]"
    />
</template>
```

**Finding icons:**
- Heroicons: https://heroicons.com/
- Lucide: https://lucide.dev/
- Use the Nuxt UI docs for icon integration details

## Composables

### File Naming

- Prefix with `use`: `useApi.ts`, `useFormat.ts`, `useKeycloak.ts`
- Place in `core/app/composables/`
- Group by feature: `chat/`, `brand/`, `player-account/`

### Composable Structure

**Standard pattern:**
```typescript
import { createSharedComposable } from '#imports'

export const useExample = createSharedComposable(() => {
    // State
    const state = ref(null)
    
    // Computed
    const computed = computed(() => state.value)
    
    // Functions
    function doSomething() {
        // ...
    }
    
    // Return
    return {
        state,
        computed,
        doSomething
    }
})
```

- Use `createSharedComposable` for singleton composables

### Type Safety

- Always provide return type for exported composables
- Use explicit types for reactive state
- Type all function parameters

**Example:**
```typescript
export interface UseDialogReturn {
    isOpen: Ref<boolean>
    open: () => void
    close: () => void
}

export const useDialog = (): UseDialogReturn => {
    const isOpen = ref<boolean>(false)
    
    function open() {
        isOpen.value = true
    }
    
    function close() {
        isOpen.value = false
    }
    
    return {
        isOpen,
        open,
        close
    }
}
```

## API Integration

### useApi vs typed clients with useOpenapi

Some endpoints / services are not typed, which means we should use `useApi`. Otherwise in most cases we should use a typed client using `useOpenapi` and generate from openapi spec in `configs/clients/<some-client>` and generate the types using `npm run generate:clients`.

### Clients

`useQuery` and `useMutation` should be colocated based on services in client composables .e.g. `usePlayerAccountClient` for the `player-account` service.

```
// In usePlayerAccountClient.ts

export const usePlayerAccountClient = () => {
   const { site } = useCoreConfig()
   const client = useOpenapi<PlayerAccountClient.paths>({
       baseUrl: '/gateway/player-account/api',
       service: 'playerAccount'
   })
   const session_expiry = useCookie('session_expiry', {
       readonly: true
   })

   function useGetUserInfo(options: UseQueryOverrideOptions = {}) {
     const query = useQuery(async () => {
       const { data } = await client.GET("/public/user/info")
       return data
     }, {
         server: false, // For data that should not be availabile on initial render
         lazy: true, // For data that should not be availabile on initial render
         ...options,
         key: "playerAccount:userInfo",
         enabled: () => options.enabled !== undefined ? toValue(options.enabled) && !!session_expiry.value : !!session_expiry.value
     })
   }

   function useUpdateOccupation() {
     return useMutation(async (body: NonNullable<PlayerAccountClient.paths['/public/user/occupation']['post']['requestBody']>['content']['application/json']) => {
         const { data } = await client.POST("/public/user/occupation", {
             body
         })

         return data
     })
   }
}
```

Don't handle errors inside useMutation, since useMutation wraps the handler function in a try/catch and manages error-handling, request status and more.
```
### Using useLazyQuery

**For data that loads after initial render:**

```typescript
const { data, error, pending } = await useLazyQuery('/api/endpoint', {
    query: computed(() => ({
        id: route.params.id
    }))
})
```

## Layered Application Architecture

### Application Structure

This project uses Nuxt's `extends` feature to create a multi-layered application:

```
├── core/                          # Core layer (shared functionality)
│   ├── nuxt.config.ts            # Core Nuxt configuration
│   ├── app/
│   │   ├── components/           # Shared components
│   │   ├── composables/          # Shared composables
│   │   ├── layouts/              # Shared layouts
│   │   ├── middleware/           # Global middleware
│   │   ├── pages/                # Core pages
│   │   └── plugins/              # Client/server plugins
│   ├── server/                   # Server-side code
│   │   └── plugins/              # Nitro plugins (CSP, etc.)
│   └── types/                    # TypeScript types
├── regions/                       # Region-specific overrides
│   ├── dk/                       # Denmark configuration
│   │   ├── nuxt.config.ts
│   │   └── app/public/logo/      # DK-specific assets
│   └── uk/                       # UK configuration
│       ├── nuxt.config.ts
│       └── app/public/logo/      # UK-specific assets
└── nuxt.config.ts                # Root config (extends regions + core)
```

### Layer Configuration

**Root `nuxt.config.ts`:**
```typescript
export default defineNuxtConfig({
    extends: [
        './regions/dk',  // or './regions/uk'
        './core'
    ]
})
```

### Accessing Layer Resources

**Import from layers using `#layers` prefix:**
```typescript
import { useChatMessages } from '#layers/core/app/composables/chat/useChatMessages'
```

### Adding Runtime Configuration

Define runtime config in `core/nuxt.config.ts`:

```typescript
export default defineNuxtConfig({
    runtimeConfig: {
        public: {
            env: '',
            region: '',
            domain: '',
            site: '',
            keycloak: {
                url: '',
                realm: '',
                clientId: ''
            }
        }
    }
})
```

Access in components/composables:
```typescript
const config = useRuntimeConfig()
const keycloakUrl = config.public.keycloak.url
```

## Testing

### Test File Naming

- Unit tests: `*.test.ts` or `*.spec.ts`
- Place in `test/unit/` or `test/nuxt/`
- Mirror the source file structure

### Test Structure

```typescript
import { describe, it, expect } from 'vitest'

describe('ComponentName', () => {
    it('should do something', () => {
        // Arrange
        const input = 'test'
        
        // Act
        const result = functionToTest(input)
        
        // Assert
        expect(result).toBe('expected')
    })
})
```

## Commit & Documentation

### Commit Messages

- Use conventional commits format with clear, descriptive messages in imperative mood

**Examples:**
```
feat: add new UIDialog component
fix: resolve z-index issue in drawer
refactor: extract theme logic to directive
docs: update API documentation
test: add tests for useAuth composable
```

### Code Comments

- Add JSDoc comments for exported functions and composables
- Explain "why" not "what" in inline comments
- Use TODO comments for future improvements: `// TODO: Add pagination support`
- Avoid adding types in JSDoc, since typescript should manage the types

### TypeScript Documentation

```typescript
/**
 * Formats a number as currency
 * 
 * @param value - The numeric value to format
 * @param options - Formatting options
 * @returns Formatted currency string
 * 
 * @example
 * asMoney(1234.56) // "kr. 1.234,56"
 */
export const asMoney = (value: number, options?: FormatOptions): string => {
    // ...
}
```

## Common Patterns

### Provide/Inject

```typescript
// Provider component
const dialog = {
    isOpen,
    open,
    close
}
provide(UIDIALOG_KEY, dialog)

// Consumer composable
export const useDialog = () => {
    const dialog = inject(UIDIALOG_KEY, undefined)
    
    if (!dialog) {
        throw new Error('No dialog available in context!')
    }
    
    return dialog
}
```

### Computed Properties

```typescript
// Use computed for derived state
const fullName = computed(() => `${firstName.value} ${lastName.value}`)

// Use computed for reactive references
const route = useRoute()
const dialogInQuery = computed(() => !!route.query.dialog)
```

### Watchers

- Use `watch` for side effects based on reactive changes

```typescript
// Watch specific sources
watch([source1, source2], ([newVal1, newVal2]) => {
    // React to changes
})

// Watch with options
watch(source, (newVal, oldVal) => {
    // ...
}, {
    immediate: true,
    deep: true
})
```

## Security & Performance

### Security Headers

CSP (Content Security Policy) is configured in `core/server/plugins/csp.ts`:

```typescript
export default defineNitroPlugin((nitroApp) => {
    nitroApp.hooks.hook('render:html', (response, { event }) => {
        const content = `frame-ancestors 'self' *.cego.dk;`
        setResponseHeader(event, 'content-security-policy', content)
    })
})
```

**Important:**
- Never expose secrets in runtime config public section
- All sensitive values should be server-only
- Use environment variables for deployment-specific values

### Performance

- **SSR disabled**: This is a client-side only application (`ssr: false`)
- Use `lazy: true` with `useQuery` for non-critical data
- Use `server: false` with `useQuery` for client-only data
- Leverage @tanstack/vue-query caching
- Use `v-once` for static content
- Use `v-memo` for expensive renders with dependencies

**Example:**
```typescript
const { data } = useQuery(async () => {
    // Fetch data
}, {
    key: 'unique-key',
    lazy: true,      // Don't fetch on mount
    server: false,   // Client-only
    staleTime: 60000 // Cache for 1 minute
})
```

## Multi-Region & Brand Support

### Region Configuration

The application supports multiple regions (DK, UK) through Nuxt layers:

**Region structure:**
```
regions/
├── dk/
│   ├── nuxt.config.ts
│   └── app/public/logo/
└── uk/
    ├── nuxt.config.ts
    └── app/public/logo/
```

**Region-specific config:**
```typescript
// regions/dk/nuxt.config.ts
export default defineNuxtConfig({
    $meta: {
        name: 'dk'
    }
})
```

### Brand Support

Brands are managed through composables and runtime config:

**Site switcher in `app.config.ts`:**
```typescript
export default defineAppConfig({
    siteSwitcher: [
        {
            name: 'spilnu',
            siteShort: 'sn',
            locale: 'da-DK',
            currency: 'DKK',
        },
        {
            name: 'happytiger',
            siteShort: 'ht',
            locale: 'en-GB',
            currency: 'GBP',
        },
    ]
})
```

**Using brand composables:**
```typescript
import { useBrands } from '#layers/core/app/composables/brand'

const { currentBrand } = useBrands()
const brandValue = currentBrand.value?.value
```

### Internationalization

Format dates, times, and currency based on runtime config:

```typescript
const { asCurrency, asDate, asDateTime } = useFormat()

// Currency formatting (uses locale and currency from runtimeConfig)
const formatted = asCurrency(1234.56) // "kr. 1.234,56" or "£1,234.56"

// Date formatting with timezone support
const date = asDate(new Date()) // "2025-10-28"
const dateTime = asDateTime(new Date()) // "2025-10-28 14:30 CET"
```

### Theming

Colors are configured in `app.config.ts` using Nuxt UI theme system:

```typescript
export default defineAppConfig({
    ui: {
        colors: {
            primary: 'violet',
            secondary: 'orange',
            neutral: 'zinc',
        }
    }
})
```

## Authentication & Real-time Features

### Keycloak Authentication

The application uses Keycloak for Single Sign-On (SSO):

**Configuration:**
```typescript
// In runtimeConfig
public: {
    keycloak: {
        url: '',
        realm: '',
        clientId: '',
        clientPrefix: ''
    },
    featureFlags: {
        keycloak: false  // Enable/disable Keycloak
    }
}
```

**Using the Keycloak composable:**
```typescript
const keycloak = useKeycloak()

// Check if enabled
if (keycloak.isEnabled) {
    // Access the Keycloak client
    const token = keycloak.client?.token
}
```

**Auth middleware:**
```typescript
// core/app/middleware/auth.global.ts
// Global middleware handles authentication checks
```

### XMPP Chat Integration

Real-time chat uses XMPP protocol via `@cego/xmpp`:

**Chat composables:**
```typescript
import { useChatMessages } from '#layers/core/app/composables/chat/useChatMessages'
import { useChatFilters } from '#layers/core/app/composables/chat/useChatFilters'

const { sendChatMessage } = useChatMessages()
const { showFilters, toggleFilters } = useChatFilters()
```

**Sending messages:**
```typescript
sendChatMessage({
    message: 'Hello',
    brand: currentBrand.value,
    displayName: 'username'
})
```

**Chat components:**
- `Chat.vue` - Main chat container with tabs
- `ChatLive.vue` - Live chat view
- `ChatSearchResults.vue` - Search results view
- `ChatFilters.vue` - Filter controls
- `ChatBadWords.vue` - Forbidden words management

## Questions & Edge Cases

When encountering ambiguous situations:

1. **Check existing code patterns** in similar components
2. **Refer to this document** for guidance
3. **Ask for clarification** with specific options:
    - "Should I use OPTION A or OPTION B for this case?"
    - Provide context and your reasoning
4. **Prioritize consistency** with existing codebase over personal preference

## File Structure Reference

```
├── core/                              # Core layer
│   ├── nuxt.config.ts                # Core Nuxt configuration
│   ├── app/
│   │   ├── app.config.ts             # App configuration (UI theme, site switcher)
│   │   ├── app.vue                   # Root app component
│   │   ├── assets/
│   │   │   └── css/
│   │   │       └── tailwind.css      # Tailwind imports
│   │   ├── components/               # Vue components (auto-imported)
│   │   │   ├── Base/                # Base components (DefaultLayout, etc.)
│   │   │   ├── Chat/                # Chat feature components
│   │   │   ├── Emoji/               # Emoji picker components
│   │   │   ├── Logo.vue             # Logo component
│   │   │   ├── Can.vue              # Permission components
│   │   │   └── Cannot.vue
│   │   ├── composables/             # Composables (auto-imported)
│   │   │   ├── brand/               # Brand management composables
│   │   │   ├── chat/                # Chat composables
│   │   │   ├── player-account/      # Player account composables
│   │   │   ├── useApi.ts            # API client
│   │   │   ├── useFormat.ts         # Formatting utilities
│   │   │   └── useKeycloak.ts       # Keycloak auth
│   │   ├── layouts/
│   │   │   └── default.vue          # Default layout
│   │   ├── middleware/
│   │   │   └── auth.global.ts       # Global auth middleware
│   │   ├── pages/
│   │   │   └── index.vue            # Home page
│   │   └── plugins/
│   │       ├── keycloak.client.ts   # Keycloak plugin
│   │       └── xmpp.client.ts       # XMPP chat plugin
│   ├── public/
│   │   └── silent-check-sso.html    # Keycloak SSO check
│   ├── server/
│   │   ├── plugins/
│   │   │   └── csp.ts               # Content Security Policy
│   │   └── tsconfig.json
│   └── types/                        # TypeScript types (auto-imported)
│       ├── chat/                     # Chat-related types
│       ├── promotion/                # Promotion types
│       ├── domain.ts
│       ├── global.ts
│       ├── keycloak.ts
│       └── route.ts
├── regions/                           # Region-specific layers
│   ├── dk/                           # Denmark
│   │   ├── nuxt.config.ts
│   │   └── app/public/logo/
│   └── uk/                           # United Kingdom
│       ├── nuxt.config.ts
│       └── app/public/logo/
├── nuxt.config.ts                    # Root config (extends layers)
├── package.json
├── tsconfig.json
└── eslint.config.mjs
```

## Version Information

- **Nuxt**: 4.1.x
- **Vue**: 3.5.x
- **TypeScript**: 5.x
- **Nuxt UI**: 4.0.0-alpha.1
- **Tailwind CSS**: 4.x
- **Node**: >= 22.x
- **Package Manager**: npm

## Key Dependencies

- **@tanstack/vue-query**: 5.87.x - Data fetching and caching
- **@peterbud/nuxt-query**: 1.x - Nuxt integration for TanStack Query
- **@vueuse/core**: 13.9.x - Vue composition utilities
- **es-toolkit**: 1.39.x - Modern utility library
- **keycloak-js**: 26.2.x - Authentication
- **dayjs**: 1.11.x - Date/time manipulation
- **@cego/xmpp**: 2.1.x - Real-time chat

---

**Note**: This document is specific to the @spilnu/backoffice-core project. Keep it updated as patterns and dependencies evolve.