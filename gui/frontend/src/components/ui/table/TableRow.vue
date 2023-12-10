<script setup lang="ts">
import { cn } from '@/lib/utils'
import { useCssVar } from '@vueuse/core';
import { ref } from 'vue';

const props = defineProps<{ class?: string }>()

const el = ref()
const mutedFg = useCssVar('--muted-foreground', el)

</script>

<template>
  <tr :class="cn('border-b border-muted transition-colors hover:bg-muted data-[state=selected]:bg-muted', props.class)">
    <slot />
  </tr>
</template>
