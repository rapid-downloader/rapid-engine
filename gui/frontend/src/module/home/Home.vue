<script setup lang="ts">
import { Button } from '@/components/ui/button'
import Header from '@/components/Header.vue';
import DownloadList from './components/DownloadList.vue';
import DownloadListSkeleton from './components/DownloadListSkeleton.vue';
import { computed, watch, onUnmounted, ref, onMounted } from 'vue';
import XTooltip from '@/components/ui/tooltip/XTooltip.vue';
import Filter from './components/Filter.vue';
import Entries from './api'
import XDialog from '@/components/ui/dialog/XDialog.vue';
import DownloadDialog from '../download/components/DownloadDialog.vue';
import { useNow, useTimeout } from '@vueuse/core';
//@ts-ignore
import { EventsOn } from '@/../wailsjs/runtime'
import { Download } from './types';

const types = [
    { value: 'Document', label: 'Document' },
    { value: 'Video', label: 'Video' },
    { value: 'Audio', label: 'Audio' },
    { value: 'Compressed', label: 'Compressed' },
    { value: 'Other', label: 'Other' },
]
const statuses = [
    { value: 'Downloading', label: 'Downloading' },
    { value: 'Queued', label: 'Queued' },
    { value: 'Paused', label: 'Paused' },
    { value: 'Failed', label: 'Failed' },
    { value: 'Stoped', label: 'Stoped' },
    { value: 'Completed', label: 'Completed' },
]

const loading = ref(true)
const dlentries = ref<Record<string, Download>>({})

const entries = Entries()

onMounted(async () => {
    dlentries.value = await entries.all()
    loading.value = false
})

onUnmounted(async () => {
    await entries.updateAll(dlentries.value)
})

async function fetched(result: Download) {
    dlentries.value[result.id] = result
    await entries.update(result)
}

const filteredtype = ref<string[]>([])
const filteredstatus = ref<string[]>([])

const items = computed(() => {
    let filtered = dlentries.value
    if (filteredtype.value.length > 0) {
        filtered = Object.fromEntries(
            Object.entries(dlentries.value).filter(([, entry]) => filteredtype.value.some(val => val.includes(entry.type)))
        )
    }

    if (filteredstatus.value.length > 0) {
        filtered = Object.fromEntries(
            Object.entries(dlentries.value).filter(([, entry]) => filteredstatus.value.some(val => val.includes(entry.status)))
        )
    }
    
    return filtered
})

interface Progress {
    id: string
    index: number
    downloaded: number
    size: number
    progress: number
    done: boolean
}

const now = useNow()
const { ready, start } = useTimeout(1000, { controls: true })

async function update(progress: Progress) {
    if (dlentries.value[progress.id].downloadedChunks === undefined) {
        dlentries.value[progress.id].downloadedChunks = new Array<number>()
    }

    dlentries.value[progress.id].downloadedChunks[progress.index] = progress.downloaded
    dlentries.value[progress.id].status = 'Downloading'

    const downloadedTotal = dlentries.value[progress.id]
        .downloadedChunks.reduce((total, chunk) => total + chunk)
    
    // calculate the total downloaded percentage
    const size = dlentries.value[progress.id].size
    dlentries.value[progress.id].progress = (downloadedTotal / size) * 100

    // refresh calculation every second
    if (ready.value) {
        // calculcate the download speed
        const elapsedSecond = new Date(
                now.value.getTime() - new Date(dlentries.value[progress.id].date).getTime()
            ).getTime() / 1000

        const speed = downloadedTotal / elapsedSecond
        dlentries.value[progress.id].speed = speed

        // calculate time left
        const remainingSize = size - downloadedTotal
        dlentries.value[progress.id].timeLeft = remainingSize / speed
        
        start()
    }

    if (progress.done) {
        dlentries.value[progress.id].status = 'Completed'
        dlentries.value[progress.id].timeLeft = 0
        dlentries.value[progress.id].progress = 100

        await entries.update(dlentries.value[progress.id])
    }
}

EventsOn('progress', async (...event: any) => {
    update(event[0] as Progress)
})

</script>

<template>
    <Header>
        <div class="flex gap-3">
            <x-tooltip text="New Download" location="bottom">
                <x-dialog title="New Download" description="Provide a link to start a new download">
                    <template v-slot:trigger>
                        <Button class="flex gap-2 bg-accent justify-center hover:bg-accent/90">
                            <i-fluent-add-16-filled class="text-accent-foreground" />
                        </Button>
                    </template>

                    <template v-slot:content>
                        <download-dialog @fetched="fetched" /> 
                    </template>
                </x-dialog>
            </x-tooltip>

            <x-tooltip text="New batch download" location="bottom">
                <Button class="flex gap-2 border border-accent group hover:bg-accent" variant="outline">
                    <i-fluent-add-square-multiple-16-filled class="text-accent text-lg group-hover:text-accent-foreground" />
                </Button>
            </x-tooltip>

            <x-tooltip text="Resume all" location="bottom">
                <Button class="flex gap-2 border border-accent group hover:bg-accent" variant="outline">
                    <i-fluent-play-16-filled class="text-accent group-hover:text-accent-foreground" />
                </Button>
            </x-tooltip>

            <x-tooltip text="Pause all" location="bottom">
                <Button class="flex gap-2 border border-accent group hover:bg-accent" variant="outline">
                    <i-fluent-pause-16-filled class="text-accent group-hover:text-accent-foreground" />
                </Button>
            </x-tooltip>
            
            <x-tooltip text="Stop all" location="bottom">
                <Button class="flex gap-2 border border-accent group hover:bg-accent" variant="outline">
                    <i-fluent-stop-16-filled class="text-accent group-hover:text-accent-foreground" />
                </Button>    
            </x-tooltip>
        </div>
    </Header>

    <div class="flex flex-col gap-2 mt-5 items-start">
        <div class="flex gap-3">
            <Filter v-if="Object.keys(dlentries).length > 0" :options="types" v-model="filteredtype">
                <Button variant="outline" size="sm" class="flex gap-1 border-muted rounded-md">
                    <i-radix-icons-mixer-horizontal />
                    <p>Type</p>
                </Button>
            </Filter>

            <Filter v-if="Object.keys(dlentries).length > 0" :options="statuses" v-model="filteredstatus">
                <Button variant="outline" size="sm" class="flex gap-1 border-muted rounded-md">
                    <i-radix-icons-mixer-horizontal />
                    <p>Status</p>
                </Button>
            </Filter>
        </div>

        <download-list-skeleton v-if="loading" />
        <download-list v-else class="w-full" :items="items"/>
    </div>
</template>
