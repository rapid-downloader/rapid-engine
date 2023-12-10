<script setup lang="ts">
import { Input } from '@/components/ui/input'
import { useRouteQuery } from '@vueuse/router'

const search = useRouteQuery<string>('search')
defineProps<{
    noSearch?: boolean

}>()

</script>

<template>
    <div class="flex gap-5 justify-between">
        <div class="flex gap-5 items-center">
            <span class="flex gap-5">
                <button variant="link">
                    <i-fluent-arrow-left-16-regular @click="$router.back()" class="text-lg text-foreground"/>
                </button>
                <button variant="link">
                    <i-fluent-arrow-right-16-regular @click="$router.forward()" class="text-lg text-foreground"/>
                </button>
            </span>
            
            <slot />
        </div>
        
        <Input v-if="!noSearch" placeholder="Search" v-model="search" class="max-w-[20rem] rounded bg-secondary">
            <template v-slot:append-icon>
                <i-fluent-search-16-filled class="text-muted-foreground"/>
            </template>
        </Input>
    </div>
</template>