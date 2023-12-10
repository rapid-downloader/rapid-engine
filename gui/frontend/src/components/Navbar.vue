<script setup lang="ts">
import { useRoute } from 'vue-router';
import { Button } from '@/components/ui/button';

import Nav from './Nav.vue';

const route = useRoute()
</script>

<template>
    <div class="flex my-1 ml-1">
        <div class="bg-background flex flex-col">
            <Nav to="/" tooltip="Downloads">
                <i-fluent-document-multiple-16-regular
                    :class="`text-xl p-0 text-primary-foreground ${route.path === '/' ? 'brightness-100' : 'brightness-75 group-hover:brightness-150'}`" />
            </Nav>
            <Nav to="/log" tooltip="Logs">
                <i-fluent-bug-16-regular
                    :class="`text-xl p-0 text-primary-foreground ${route.path === '/log' ? 'brightness-100' : 'brightness-75 group-hover:brightness-150'}`" />
            </Nav>
            <Nav to="/setting" tooltip="Settings">
                <i-fluent-settings-16-regular
                    :class="`text-xl p-0 text-primary-foreground ${route.path === '/setting' ? 'brightness-100' : 'brightness-75 group-hover:brightness-150'}`" />
            </Nav>
            <Nav to="/about" tooltip="About">
                <i-fluent-question-circle-16-regular
                    :class="`text-xl p-0 text-primary-foreground ${route.path === '/about' ? 'brightness-100' : 'brightness-75 group-hover:brightness-150'}`" />
            </Nav>
        </div>
    </div>
</template>