<script setup lang="ts">
import { useVModel } from '@vueuse/core'
import { cn } from '@/lib/utils'

const props = defineProps<{
  defaultValue?: string | number
  modelValue?: string | number
  placeholder?: string
}>()

const emits = defineEmits<{
  (e: 'update:modelValue', payload: string | number): void
}>()

const modelValue = useVModel(props, 'modelValue', emits, {
  passive: true,
  defaultValue: props.defaultValue,
})
</script>

<template>
  <div class="relative w-full">
    <div class="absolute top-1/2 left-1 transform -translate-y-1/2 w-[1.75rem]">
      <slot name="prepend-icon"/>
    </div>

    <input v-model="modelValue" type="text" :placeholder="placeholder" :class="cn('flex h-9 w-full rounded border border-muted bg-transparent px-3 py-1 text-sm shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50', $attrs.class ?? '')">
    <div class="absolute top-1/2 right-1 transform -translate-y-1/2 w-[1.75rem]">
      <slot name="append-icon"/>
    </div>
  </div>
</template>
