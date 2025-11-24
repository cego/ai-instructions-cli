# AI Instructions for @spilnu/core and Related Projects

This document provides coding guidelines and preferences for AI assistants working on projects that use `@spilnu/core` or similar codebases.

## Project Context

- **Framework**: Nuxt 4 with Vue 3 and TypeScript
- **Package Type**: Nuxt module providing core functionality for multiple brands
- **Target Platforms**: Desktop and mobile web applications (responsive design)
- **Primary Markets**: Danish and English gaming/casino applications
- **Styling**: SCSS with mobile-first approach
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

- Use PascalCase for component files: `AccountBalanceOverview.vue`, `UIButton.vue`,
- Prefix UI components with `UI`: `UIDialog.vue`, `UIInput.vue`
- Use kebab-case for composable files: `usePlayerAccountClient.ts`, `useApi.ts`

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
    <!-- Component template -->
</template>

<style lang="scss" scoped>
// Component styles
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
export interface UIButtonProps {
    variant?: 'primary' | 'secondary' | 'danger'
    disabled?: boolean
    loading?: boolean
}

const props = withDefaults(defineProps<UIButtonProps>(), {
    variant: 'primary',
    disabled: false,
    loading: false
})
</script>
```

### Props Definition

- Use TypeScript interface + `defineProps<T>()`
- Use `withDefaults(defineProps<UIButtonProps>() { ... })` when setting defaults is needed

**Example**

```typescript
export interface UILoaderProps {
    bg?: boolean
    color?: string
    inline?: boolean
    width?: string
    height?: string
    text?: string | false
}

const props = defineProps<UILoaderProps>()
```

## TypeScript Conventions

### Type Definitions

- Store shared types in `src/runtime/types/`
- Use type imports: `import type { ... } from '...'`
- Prefer interfaces over types for object shapes
- Define component-specific interfaces and extract to separate type file in `types/components/` if used by multiple components

### Type Naming

- Prefix component prop interfaces with component name: `UIButtonProps`, `UIDialogProps`
- Use descriptive names: `UseQueryReturn`, `ApiResponse`, `ServiceOptions`
- Suffix return types with `Return`: `UseDisplayReturn`, `UseApiReturn`
- Suffix options with `Options`: `UseQueryOptions`, `UseFormatOptions`

### Imports

**Import order:**

1. Vue core imports
2. Nuxt imports (`#imports`, `#app`)
3. Third-party packages
4. Local type imports (`@core/types`)
5. Relative imports

**Use auto-imports:**

- Nuxt auto-imports components, composables and Vue.js APIs to use across the application
- Nuxt utilities are auto-imported
- Vue utilities are auto-imported (ref, computed, watch, etc.)

**Example:**

```typescript
import type { TransitionProps, HTMLAttributes } from 'vue'
import type { UIDialogInstance } from '@core/types'
import { onUnmounted, onMounted, watch, computed, ref } from '#imports'
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

## SCSS/CSS Conventions

### Mobile-First Approach

**Always design mobile-first**, then add larger breakpoint styles:

```scss
.container {
    // Mobile styles (default)
    padding: 15px;

    // Tablet
    @media (min-width: $screen-md) {
        padding: 30px;
    }

    // Desktop
    @media (min-width: $screen-lg) {
        padding: 40px;
    }
}
```

### Media Query Mixins

**Use provided mixins instead of raw media queries:**

```scss
.container {
    padding: 15px;

    @include tablet() {
        padding: 30px;
    }

    @include desktop() {
        padding: 40px;
    }

    @include desktop-xl() {
        max-width: 1024px;
    }
}
```

### Variables & Functions

- Use SCSS variables from `variables.scss` (auto-imported)
- Use the `z()` function for z-index values: `z-index: z(dialog)`
- Avoid using deprecated Sass functions (e.g., `darken`, `lighten`, `adjust-color`). Use modern alternatives from the `sass:color` module like `color.adjust()`, `color.mix()`.
- Use math functions: `math.floor()`, `math.ceil()`

**Required imports (already available globally):**

```scss
@use "sass:math";
@use "sass:color";
```

### BEM-Like Naming

- We prefer BEM-like naming with `&` for nested elements, but allow simple class names when component is small and straightforward

**Example:**

```scss
.dialog {
    // Block

    &__header {
        // Element
    }

    &__footer {
        // Element
    }

    &--wide {
        // Modifier
    }

    &--bounce-down {
        // Modifier
    }
}
```

### Scoped Styles

- Always use `<style lang="scss" scoped>` in components
- Use `:deep()` for styling elements rendered by components or slots within the component
- Avoid `::v-deep` (deprecated)

**Example:**

```vue
<template>
    <div v-show="modelValue" class="'alert'">
        <UIButton
            v-if="closeable"
            as="close"
            class="alert__btn--close"
            @click.prevent="onCloseClick"
        >
            <!-- ... -->
        </UIButton>

        <div class="alert__message">
            <!-- ... -->
            <div>
                <slot>
                    <p>{{ message }}</p>
                </slot>
            </div>
        </div>
    </div>
</template>
<style lang="scss" scoped>
.alert {
    padding: 12px;
    position: relative;
    font-size: $font-size-sm;
    z-index: 3;

    // Target child p elements in slot
    :slotted(p) {
        color: inherit;
        margin: 0 0 5px;
        padding: 0;
        font-size: inherit;
    }

    // Target child element in UIButton
    :deep(.btn__content) {
        font-size: inherit;
    }
}
</style>
```

## Composables

### File Naming

- Prefix with `use`: `usePlayerAccountClient.ts`, `useApi.ts`, `useDayjs.ts`
- Place in `src/runtime/composables/`
- Data-fetching composables go in `src/runtime/composables/data/`

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
         enabled: () => options.eanbled !== undefined ? toValue(options.enabled) && !!session_expiry.value : !!session_expiry.value
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

## Nuxt Module Development

### Module Structure

- Main module file: `src/module.ts`
- Runtime code: `src/runtime/`
- Type definitions: `src/runtime/types/`
- Build config: `build.config.ts`

### Module Options

- Define options interface in `src/runtime/types/global.ts`
- Provide sensible defaults in `src/options/index.ts`
- Make options available via `useCoreConfig()`

### Adding Components

```typescript
addComponentsDir({
    path: resolve(runtimeDir, 'components'),
    pathPrefix: true // Enables UI/Button.vue -> UIButton
})
```

### Adding Composables

```typescript
addImportsDir([
    resolve(runtimeDir, 'composables'),
    resolve(runtimeDir, 'composables', 'data'),
])
```

### Adding Plugins

```typescript
addPlugin({
    src: resolve(runtimeDir, 'plugins/dev'),
    mode: 'client' // or 'server' or omit for both
})
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
test: add tests for usePlayerAccountClient composable
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

- All CSP rules are configured in `src/module.ts`
- Use `getTrustedOrigins()` helper for environment-based domains
- Never expose secrets in public module options (they become public)

### Performance

- Use `lazy: true` for non-critical data (below the fold)
- Use `server: false` when data is client-only and for non-critical data (below the fold)
- Leverage query caching with TTL
- Use `v-once` for static content
- Use `v-memo` for expensive renders

## Brand-Specific Considerations

### Multi-Brand Support

- Define brand-specific values in SCSS variables
- Use runtime config for brand-specific API endpoints
- Support multiple jurisdictions (DK and UK)
- Handle brand-specific translations

### Theming

- Make colors configurable via SCSS variables
- Use CSS custom properties for runtime theming

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
├── src/
│   ├── module.ts                      # Main module entry
│   ├── options/                       # Module configuration
│   ├── runtime/
│   │   ├── components/                # Vue components (auto-imported)
│   │   │   ├── UI/                   # UI components (prefixed with UI)
│   │   │   ├── Account/              # Feature-specific components
│   │   │   └── Base/                 # Base components
│   │   ├── composables/              # Composables (auto-imported)
│   │   │   └── data/                 # Data-fetching composables
│   │   ├── constants/                # Constants and enums
│   │   ├── directives/               # Vue directives
│   │   ├── middleware/               # Nuxt middleware
│   │   ├── plugins/                  # Nuxt plugins
│   │   ├── scss/                     # Global styles
│   │   │   ├── base/                # Base styles, mixins, typography
│   │   │   ├── components/          # Component-specific styles
│   │   │   ├── variables.scss       # SCSS variables
│   │   │   └── main.scss            # Main style entry
│   │   ├── server/                   # Server utilities
│   │   │   ├── api/                 # API routes
│   │   │   └── plugins/             # Nitro plugins
│   │   ├── types/                    # TypeScript types
│   │   └── utils/                    # Utility functions
│   └── utils/                        # Build-time utilities
├── test/                             # Tests
│   ├── unit/                        # Unit tests
│   └── nuxt/                        # Nuxt-specific tests
└── playground/                       # Development playground
```

## Version Information

- **Nuxt**: 4.x
- **Vue**: 3.x
- **TypeScript**: 5.x
- **Node**: >= 22.x
- **Package Manager**: npm

---

**Note**: This document should be copied to any project using `@spilnu/core` as a distribution copi. Keep this file (the original) updated as patterns evolve.
