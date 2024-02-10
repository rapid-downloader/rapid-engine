<script setup lang="ts">
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from '@/components/ui/dialog'
import { ref } from 'vue';

defineProps<{
    title?: string
    description?: string,
}>()

const open = ref(false)
</script>

<template>
    <Dialog >
        <DialogTrigger>
            <button>
                <slot/>
            </button>
        </DialogTrigger>
        <DialogContent>
            <DialogHeader>
                <DialogTitle>{{ title }}</DialogTitle>
                <DialogDescription>
                    {{ description }}
                </DialogDescription>
            </DialogHeader>

            <slot name="content" />
            
            <DialogFooter>
                <slot name="footer" />
            </DialogFooter>
        </DialogContent>
    </Dialog>
</template>