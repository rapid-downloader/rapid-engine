<script setup lang="ts">
import { ToastProps, toastVariants } from '.';
import { Button } from '../button';
import { Tooltip } from '../tooltip';
import { Progress } from '../progress';

const props = defineProps<ToastProps>()
const emit = defineEmits<{
    (e: 'closeToast'): void
}>()

function dismiss() {
    if (props.onDismiss) props.onDismiss({ id: props.id! })
    emit('closeToast')
}

function click() {
    if (props.action) props.action.onClick()
    emit('closeToast')
}

</script>

<template>
    <div :class="`w-[23rem] flex items-center gap-4 p-3 bg-background rounded-md border relative group ${toastVariants({ border: type })}`">
        <i-material-symbols-check-circle v-if="type === 'success'" class="text-success"/>
        <i-material-symbols-info v-if="type === 'info'" class="text-info"/>
        <i-clarity-error-standard-solid v-if="type === 'warning'" class="text-warning"/>
        <i-bx-bxs-error v-if="type === 'destructive'" class="text-destructive"/>

        <div>
            <p :class="`text-sm font-medium ${toastVariants({ text: type })}`">{{ title }}</p>
            <p class="text-xs text-foreground/50</p>">{{ description }}</p>
        </div>

        <Tooltip v-if="action && action.icon" :text="action.label" location="left">
            <Button size="sm" @click="click" :class="`ml-auto ${action.icon ? 'aspect-square' : ''}`">
                <component v-if="action.icon" :is="action.icon" />
                <p v-else-if="action.label" class="text-xs">{{ action.label }}</p>
            </Button>
        </Tooltip>

        <Button v-else-if="action" size="sm" @click="click" :class="`ml-auto ${action.icon ? 'aspect-square' : ''}`">
            <component v-if="action.icon" :is="action.icon" />
            <p v-else-if="action.label" class="text-xs">{{ action.label }}</p>
        </Button>

        <Button v-if="dismissible" @click="dismiss" size="icon" class="hidden group-hover:block absolute -right-2 -top-3 rounded-full w-6 h-6 aspect-square">
            <i-material-symbols-close-rounded class="text-xs mx-auto" />
        </Button>
    </div>
</template>