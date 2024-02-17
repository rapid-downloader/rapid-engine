<script setup lang="ts">
import Navbar from './components/Navbar.vue';
import { Toaster } from '@/components/ui/sonner'

</script>

<template>
    <Toaster />
    <div class="flex bg-background min-w-[50rem]">
        <Navbar class="sticky top-0 left-0 p-1 bg-background self-start h-screen z-[49]"/>

        <div class="flex flex-col w-full mx-3 my-2">
            <router-view/>
        </div>
    </div>
</template>
