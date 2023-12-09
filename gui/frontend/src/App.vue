<script setup lang="ts">
import Navbar from './components/Navbar.vue';

</script>

<template>
    <div class="flex bg-background min-w-[50rem]">
        <Navbar />

        <div class="flex flex-col w-full mx-3 my-2">
            <router-view/>
        </div>
    </div>
</template>
