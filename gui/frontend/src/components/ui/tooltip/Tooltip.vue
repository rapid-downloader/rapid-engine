<script setup lang="ts">
import {
    TooltipRoot,
    TooltipContent,
    TooltipProvider,
    TooltipTrigger
} from '@/components/ui/tooltip'

const props = defineProps<{
    text: string
    location?: "top" | "right" | "bottom" | "left"
}>()
</script>

<template>
    <TooltipProvider>
        <TooltipRoot>
            <TooltipTrigger as-child>
                <slot />
            </TooltipTrigger>
            <TooltipContent class="rounded border border-muted" avoid-collisions :side="location">
                <p class="brightness-75">{{ props.text }}</p>
            </TooltipContent>
        </TooltipRoot>
    </TooltipProvider>
</template>