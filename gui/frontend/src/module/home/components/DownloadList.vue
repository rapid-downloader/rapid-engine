<script setup lang="ts">
import Cato from '@/assets/images/cato.svg'
import { Download, Sort } from '../types'
import { parseSize, parseDate, statusColor, parseTimeleft } from '@/lib/parse';
import FileType from './FileType.vue'
import StatusIcon from './StatusIcon.vue';

import {
    Table,
    TableBody,
    TableCaption,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from '@/components/ui/table'
import { computed, ref } from 'vue';
import { useRouteQuery } from '@vueuse/router';
import { useNow } from '@vueuse/core';

const props = defineProps<{
    items: Record<string, Download>,
}>()

const asc = ref(false)
const selected = ref<Sort>('date')

function sort(row: Sort) {
    selected.value = row
    asc.value = !asc.value
}

const search = useRouteQuery('search', '')

const items = computed(() => {
    return Object.entries(props.items)
        .filter(([_, item]) => item.name.toLowerCase().includes(search.value.toLowerCase()))
        .sort(([, v1], [, v2]) => {
            if (selected.value === 'date') {
                return asc.value
                    ? new Date(v1.date).getTime() - new Date(v2.date).getTime()
                    : new Date(v2.date).getTime() - new Date(v1.date).getTime()
            }

            if (selected.value === 'size') {
                return asc.value
                    ? v1.size - v2.size
                    : v2.size - v1.size
            }

            return asc.value
                ? v1.name.localeCompare(v2.name)
                : v2.name.localeCompare(v1.name)
        })
})

</script>

<template>
    <div
        :class="`${!items || Object.keys(items).length === 0 ? '' : 'bg-secondary border border-muted mb-3 rounded-md px-2'}`">
        <div v-if="!items || Object.keys(items).length === 0" class="w-fit mx-auto">
            <img :src="Cato" alt="empty" class="mx-auto my-auto w-[20rem] h-[80vh]">
        </div>

        <Table v-else class="min-w-max">
            <table-caption></table-caption>
            <table-header>
                <table-row class="hover:bg-secondary border-muted-foreground">
                    <table-head class="w-[2rem]">
                        Type
                    </table-head>
                    <table-head @click="sort('name')" class="cursor-pointer">
                        <div class="flex justify-between items-center">
                            <p>Name</p>
                            <i-radix-icons-caret-sort />
                        </div>
                    </table-head>
                    <table-head @click="sort('size')" class="cursor-pointer w-[10%]">
                        <div class="flex justify-between items-center">
                            <p>Size</p>
                            <i-radix-icons-caret-sort />
                        </div>
                    </table-head>
                    <table-head class="w-[8%]">
                        Progress
                    </table-head>
                    <table-head class="w-[8%]">
                        Time Left
                    </table-head>
                    <table-head class="w-[15%] lg:w-[10%]">
                        Speed
                    </table-head>
                    <table-head class="w-[10%]">
                        Status
                    </table-head>
                    <table-head @click="sort('date')" class="cursor-pointer w-[20%]">
                        <div class="flex justify-between items-center">
                            <p>Date</p>
                            <i-radix-icons-caret-sort />
                        </div>
                    </table-head>
                </table-row>
            </table-header>
            <table-body>
                <table-row v-for="[id, item] in items" :key="id" class="group relative cursor-pointer" @click="$router.push(`/download/${item.id}`)">
                    <table-cell class="font-medium">
                        <file-type :type="item.type" /> 
                    </table-cell>
                    <table-cell>
                        <p class="w-[95%] truncate">{{ item.name }}</p>
                    </table-cell>
                    <table-cell>
                        {{ parseSize(item.size) }}
                    </table-cell>
                    <table-cell>
                        {{ `${item.progress.toFixed(2)}%` }}
                    </table-cell>
                    <table-cell>
                        {{ parseTimeleft(item.timeLeft) }}
                    </table-cell>
                    <table-cell>
                        {{ `${parseSize(item.speed)}/s` }}
                    </table-cell>
                    <table-cell :class="`font-medium flex gap-1 items-center ${statusColor(item.status)}`">
                        <status-icon :status="item.status" /> {{ item.status }}
                    </table-cell>
                    <table-cell>
                        {{ parseDate(item.date) }}
                    </table-cell>
                </table-row>
            </table-body>
        </Table>
    </div>
</template>