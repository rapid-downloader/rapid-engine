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
import { computed, onMounted, onUnmounted, ref } from 'vue';
import { useRouteQuery } from '@vueuse/router';
import { Button } from '@/components/ui/button';
import Dialog from '@/components/ui/dialog/XDialog.vue';
import { Actionable } from '@/components/ui/dialog';
import Confirmation from '@/components/Confirmation.vue';
import { Checkbox } from '@/components/ui/checkbox';

const props = defineProps<{
    items: Record<string, Download>,
}>()

const asc = ref(false)
const selected = ref<Sort>('date')

function sort(row: Sort) {
    selected.value = row
    asc.value = !asc.value
}

const emit = defineEmits<{
    (e: 'paginate'): void,
    (e: 'delete', id: string): void
}>()

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

const pagination = ref<HTMLElement>()
function isVisible(element: HTMLElement) {
    const rect = element.getBoundingClientRect()
    
    return (
        rect.top >= 0 &&
        rect.left >= 0 &&
        rect.bottom <= 1000 &&
        rect.right <= (window.innerWidth || document.documentElement.clientWidth)
    );
}

function onScroll() {
    console.log(isVisible(pagination.value!));
    
    if (pagination.value && isVisible(pagination.value)) {
        emit('paginate')
    }
}

function isScrollable() {
    const documentElement = document.documentElement;
    const body = document.body;

    const documentHeight = Math.max(
        documentElement.scrollHeight,
        body.scrollHeight
    );

    return documentHeight > documentElement.clientHeight;
}

onMounted(() => {
    if (isScrollable()) {
        document.addEventListener('scroll', onScroll)
        document.addEventListener('resize', onScroll)
    }
    
})

onUnmounted(() => {
    document.removeEventListener('scroll', onScroll)
    document.removeEventListener('resize', onScroll)
})

const wantRemoveFromDisk = ref(false)
function remove(id: string) {

}

</script>

<template>
    <div :class="`${!items || Object.keys(items).length === 0 ? '' : 'bg-secondary border border-muted mb-3 rounded-md px-2'}`">
        <div v-if="!items || Object.keys(items).length === 0" class="w-fit mx-auto">
            <img :src="Cato" alt="empty" class="mx-auto my-auto w-[20rem] h-[80vh]">
        </div>

        <div ref="container" v-else class="">
            <Table  class="min-w-max">
                <table-caption></table-caption>
                <table-header>
                    <table-row class="hover:bg-secondary border-muted-foreground">
                        <table-head class="w-[2rem]">Type</table-head>
                        <table-head @click="sort('name')" class="cursor-pointer">
                            <div class="flex justify-between items-center">
                                <p>Name</p>
                                <i-radix-icons-caret-sort />
                            </div>
                        </table-head>
                        <table-head @click="sort('size')" class="cursor-pointer min-w-[7rem]">
                            <div class="flex justify-between items-center">
                                <p>Size</p>
                                <i-radix-icons-caret-sort />
                            </div>
                        </table-head>
                        <table-head class="w-[8%]">Progress</table-head>
                        <table-head class="min-w-[5rem]">Time Left</table-head>
                        <table-head class="min-w-[7rem] lg:w-[10%]">Speed</table-head>
                        <table-head class="w-[10%]">Status</table-head>
                        <table-head @click="sort('date')" class="cursor-pointer min-w-[12rem]">
                            <div class="flex justify-between items-center">
                                <p>Date</p>
                                <i-radix-icons-caret-sort />
                            </div>
                        </table-head>
                    </table-row>
                </table-header>
                <table-body>
                    <table-row v-for="[id, item] in items" :key="id" class="group relative cursor-pointer">
                        <table-cell class="font-medium">
                            <file-type :type="item.type" /> 
                        </table-cell>
                        <table-cell class="min-w-[30rem] relative group/cell">
                            <p class="w-[95%] group-hover/cell:w-[75%] truncate">{{ item.name }}</p>
                            <div class="hidden group-hover/cell:flex absolute top-1/2 -translate-y-1/2 right-5 gap-2 py-2">
                                <Button size="icon" variant="ghost" class="bg-muted group/action rounded-full hover:bg-secondary hover:text-foreground w-7 h-7">
                                    <i-fluent:pause-16-regular class="group-hover/action:hidden" />
                                    <i-fluent:pause-16-filled class="hidden group-hover/action:block" />
                                </Button>
                                <Button size="icon" variant="ghost" class="group/action rounded-full text-warning hover:bg-warning hover:text-warning-foreground w-7 h-7">
                                    <i-fluent:stop-16-regular class="group-hover/action:hidden"/>
                                    <i-fluent:stop-16-filled class="hidden group-hover/action:block"/>
                                </Button>
                                <Confirmation 
                                    title="Delete download" 
                                    :description="`This will delete ${item.name} from the list forever. If this is an on going download, it will be stopped and then deleted. Are you sure?`" 
                                    :actionable="{ label: 'Delete', type: 'destructive', action: () => remove(item.id) }" >
                                    <Button size="icon" variant="ghost" class="rounded-full text-destructive hover:bg-destructive hover:text-destructive-foreground w-7 h-7">
                                        <i-fluent:delete-16-regular class="group-hover/action:hidden"/>
                                        <i-fluent:delete-16-filled class="hidden group-hover/action:block"/>
                                    </Button>

                                    <template #additional>
                                        <div class="w-full flex gap-2 items-center">
                                            <Checkbox v-model="wantRemoveFromDisk" />
                                            <p class="text-xs opacity-80">Remove from disk</p>
                                        </div>
                                    </template>
                                </Confirmation>
                            </div>
                        </table-cell>
                        <table-cell>{{ parseSize(item.size) }}</table-cell>
                        <table-cell>{{ `${item.progress.toFixed(2)}%` }}</table-cell>
                        <table-cell>{{ parseTimeleft(item.timeLeft) }}</table-cell>
                        <table-cell>{{ `${parseSize(item.speed)}/s` }}</table-cell>
                        <table-cell :class="`font-medium flex gap-1 items-center ${statusColor(item.status)}`">
                            <status-icon :status="item.status" /> {{ item.status }}
                        </table-cell>
                        <table-cell>{{ parseDate(item.date) }}</table-cell>
                    </table-row>
                </table-body>
            </Table>
        </div>

        <div ref="pagination" class="h-1 w-1"/>
    </div>
</template>