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
import { Button } from './ui/button';

interface Actionable {
    label: string
    type?: 'warning' | 'destructive' | 'default'
    action: () => void
}
const props = defineProps<{
    title: string
    description?: string
    actionable: Actionable
}>()

const open = defineModel('open', { default: false })

function accept() {
    if (props.actionable) {
        props.actionable.action()
    }

    open.value = false
}
</script>

<template>
<Dialog v-model:open="open">
    <DialogTrigger>
        <slot/>
    </DialogTrigger>
    <DialogContent>
        <DialogHeader>
            <DialogTitle>{{ title }}</DialogTitle>
            <DialogDescription v-if="description">{{ description }}</DialogDescription>
        </DialogHeader>

        <DialogFooter class="flex justify-between items-center">
            <slot name="additional" />
            
            <div class="flex gap-2">
                <Button variant="ghost">Cancel</Button>
                <Button :variant="actionable.type" @click="accept">{{ actionable.label }}</Button>
            </div>
        </DialogFooter>
    </DialogContent>
</Dialog>
</template>