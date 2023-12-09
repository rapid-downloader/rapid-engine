<script setup lang="ts">
import { ComboboxInput, type ComboboxInputProps } from 'radix-vue'
import { cn } from '@/lib/utils'

const props = defineProps<ComboboxInputProps>()
</script>

<script lang="ts">
export default {
  inheritAttrs: false,
}
</script>

<template>
  <div class="flex items-center border-b border-muted-foreground px-3" cmdk-input-wrapper>
    <i-fluent-search-16-filled class="mr-2 h-4 w-4 shrink-0 opacity-50 text-foreground" />
    <ComboboxInput
      v-bind="{ ...props, ...$attrs }"
      auto-focus
      :class="cn('flex h-10 w-full rounded-md bg-transparent py-3 text-sm outline-none placeholder:text-muted-foreground disabled:cursor-not-allowed disabled:opacity-50', $attrs.class ?? '')"
    />
  </div>
</template>
