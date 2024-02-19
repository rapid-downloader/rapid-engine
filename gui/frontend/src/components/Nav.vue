<script setup lang="ts">
import { useRoute } from 'vue-router';
import { Button } from '@/components/ui/button';
import { Tooltip } from './ui/tooltip';

const route = useRoute()

const props = defineProps<{
    to: string
    tooltip: string
}>()
</script>

<template>
    <router-link :to="props.to">
        <Tooltip :text="props.tooltip" location="right">
            <Button :class="`group shadow-none bg-background py-5 px-2 m-0.5 ${route.path === props.to ? 'bg-background brightness-150 hover:bg-background hover:brightness-150' : ''}`">
                <slot />
            </Button>
        </Tooltip>
    </router-link>
</template>