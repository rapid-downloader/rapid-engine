<script setup lang="ts">
import {
Actionable,
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from '@/components/ui/dialog'
import { useVModel } from '@vueuse/core';
import { Button } from '../button';
import { DialogClose } from 'radix-vue';

const props = defineProps<{
    title?: string
    description?: string,
    open?: boolean
    action?: Actionable
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
                <DialogTitle class="font-title">{{ title }}</DialogTitle>
                <DialogDescription>
                    {{ description }}
                </DialogDescription>
            </DialogHeader>

            <slot name="content" >
                <div class="flex gap-2 justify-end">
                    <DialogClose>
                        <Button variant="ghost">Cancel</Button>
                    </DialogClose>
                    <DialogClose v-if="props.action && props.action.immediatelyClose">
                        <Button v-if="props.action" :variant="props.action.type" @click="props.action!.action">{{ props.action.label }}</Button>
                    </DialogClose>
                    <Button v-else-if="props.action && !props.action.immediatelyClose" :variant="props.action.type" @click="props.action!.action">{{ props.action.label }}</Button>
                </div>
            </slot>
        </DialogContent>
    </Dialog>
</template>