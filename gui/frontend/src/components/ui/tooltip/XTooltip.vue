<script setup lang="ts">
import {
    Tooltip,
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
        <Tooltip>
            <TooltipTrigger as-child>
                <slot />
            </TooltipTrigger>
            <TooltipContent class="rounded border border-muted" avoid-collisions :side="location">
                <p class="brightness-75">{{ props.text }}</p>
            </TooltipContent>
        </Tooltip>
    </TooltipProvider>
</template>