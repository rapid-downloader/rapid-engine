<script setup lang="ts">
import {
  ProgressIndicator,
  ProgressRoot,
  type ProgressRootProps,
} from 'radix-vue'
import { cn } from '@/lib/utils'

const props = withDefaults(
  defineProps<ProgressRootProps & { class?: string, indetermined?: boolean }>(),
  {
    modelValue: 0,
    indetermined: false,
  },
)
</script>

<template>
  <ProgressRoot
    v-if="!indetermined"
    :class="
      cn(
        'relative h-2 w-full overflow-hidden rounded-full bg-primary/20',
        props.class,
      )
    "
    v-bind="props"
  >
    <ProgressIndicator
      :class="
        cn(
          'h-full w-full flex-1 bg-primary transition-all',
          props.class,
        )
      "
      :style="`transform: translateX(-${100 - (props.modelValue ?? 0)}%);`"
    />
  </ProgressRoot>
  <div v-else :class="cn('w-full h-2 overflow-hidden', props.class)" >
      <div class="w-full h-full bg-accent progress-bar-value" />
  </div>
</template>

<style>
.progress-bar-value {
  animation: indetermined 1s infinite linear;
  transform-origin: 0% 50%;
}

@keyframes indetermined {
  0% {
    transform:  translateX(0) scaleX(0);
  }
  40% {
    transform:  translateX(0) scaleX(0.7);
  }
  80% {
    transform:  translateX(100%) scaleX(0.2);
  }
  100% {
    transform:  translateX(100%) scaleX(0.1);
  }
}
</style>