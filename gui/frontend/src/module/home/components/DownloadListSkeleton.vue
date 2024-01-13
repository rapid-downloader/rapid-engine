<script setup lang="ts">
import { Skeleton } from '@/components/ui/skeleton'
import { Separator } from '@/components/ui/separator'

</script>

<template>
    <div class="flex gap-4 -mt-2">
        <Skeleton class="bg-muted w-[4.2rem] h-8"/>
        <Skeleton class="bg-muted w-20 h-8"/>
    </div>
    
    <div class="w-full bg-secondary flex flex-col gap-3 border border-muted mb-3 rounded-md p-2">
        <div v-for="_ in Array(20)" class="flex gap-2">
            <Skeleton class="w-[47%] h-6 bg-muted"/>
            <Skeleton class="w-[7%] h-6 bg-muted"/>
            <Skeleton class="w-[7%] h-6 bg-muted"/>
            <Skeleton class="w-[7%] h-6 bg-muted"/>
            <Skeleton class="w-[7%] h-6 bg-muted"/>
            <Skeleton class="w-[7%] h-6 bg-muted"/>
            <Skeleton class="flex-grow h-6 bg-muted"/>
        </div>
    </div>
</template>