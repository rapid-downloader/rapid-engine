<script setup lang="ts">
import { Label } from '@/components/ui/label';
import { Checkbox } from '@/components/ui/checkbox'
import {
    Popover,
    PopoverContent,
    PopoverTrigger,
} from '@/components/ui/popover'
import {
    Command,
    CommandEmpty,
    CommandGroup,
    CommandInput,
    CommandItem,
    CommandList,
} from '@/components/ui/command'
import { Separator } from '@/components/ui/separator'
import Pill from './Pill.vue'
import { ref } from 'vue';

interface Props {
    modelValue?: string[],
    title?: string,
    options?: {
        label: string,
        value: string,
    }[]
}

defineProps<Props>()
const emit = defineEmits<{
    (e: 'update:modelValue', val?: string[]): void
}>()

const selected = ref<string[]>([])

function select(val: string) {
    const i = selected.value.indexOf(val)
    if (i !== -1) selected.value.splice(i, 1)
    else selected.value.push(val)

    emit('update:modelValue', selected.value)
}

function remove(val: string) {
    const i = selected.value.indexOf(val)
    if (i !== -1) selected.value.splice(i, 1)
}

</script>

<template>
    <div class="flex gap-1 border-opacity-25 rounded ">
        <Popover class="w-fit">
            <PopoverTrigger class="relative">
                <slot />
            </PopoverTrigger>
            <PopoverContent class="absolute -left-[2.5rem] !pt-2">
                <Command class="w-[10rem] -mt-2">
                    <CommandInput placeholder="Search" />
                    <CommandList>
                        <CommandEmpty>No results found.</CommandEmpty>
                        <CommandGroup :heading="title">
                            <CommandItem v-for="(option, i) in options" :value="option.value" @select="select(option.value)">
                                <div class="flex gap-3 items-center">
                                    <Checkbox :checked="selected.includes(option.value)" v-model="selected" :value="option.value" />
                                    <Label class="cursor-pointer">{{ option.label }}</Label>
                                </div>
                            </CommandItem>
                        </CommandGroup>
                    </CommandList>
                </Command>
            </PopoverContent>
        </Popover>
        <div class="flex gap-1">
            <Separator v-if="selected.length > 0" orientation="vertical" class=""/>
            <Pill v-for="value in selected" @click="remove(value)" class="flex items-center gap-1">
                <p>{{ value }}</p>
                <i-fluent-dismiss-circle-16-regular />
            </Pill>
        </div>
    </div>
</template>