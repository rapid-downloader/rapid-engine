<script setup lang="ts">
import { Input } from '@/components/ui/input';
import { Date } from '@/module/log/types';
import { computed } from 'vue';

const props = defineProps<{
    modelValue: Date
}>()

const emit = defineEmits<{
    (e: 'update:modelValue', val: Date): void
}>()

const date = computed({
    get: () => props.modelValue,
    set: val => emit('update:modelValue', val)
})

</script>

<template>
    <Input v-model="date.day" type="text" placeholder="Day" max="31" class="w-[3rem] text-center" />
    <Input v-model="date.month" type="text" placeholder="Month" class="w-[3rem] text-center" />
    <Input v-model="date.year" type="text" placeholder="Year" class="w-[4rem] text-center" />
</template>