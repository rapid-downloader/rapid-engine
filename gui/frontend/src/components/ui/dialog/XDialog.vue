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
import { useVModel } from '@vueuse/core';
import { ref } from 'vue';

const props = defineProps<{
    title?: string
    description?: string,
    open?: boolean
}>()

const open = useVModel(props, 'open')

</script>

<template>
    <Dialog v-model:open="open">
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