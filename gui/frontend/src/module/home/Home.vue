<script setup lang="ts">
import { Button } from '@/components/ui/button'
import Header from '@/components/Header.vue';
import DownloadList from './components/DownloadList.vue';
import { computed, onMounted, onUnmounted, ref } from 'vue';
import { Download } from './types';
import XTooltip from '@/components/ui/tooltip/XTooltip.vue';
import Filter from './components/Filter.vue';
import Entries from './api'
import XDialog from '@/components/ui/dialog/XDialog.vue';
import DownloadDialog from '../download/components/DownloadDialog.vue';
import { Websocket } from '@/composable'
import { Progress } from '@/components/ui/progress';
import { useNow, useTimeout } from '@vueuse/core';

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

const dlentries = ref<Record<string, Download>>({})

const entries = Entries()

interface Progress {
    id: string
    index: number
    downloaded: number
    size: number
    progress: number
    done: boolean
}

const socket = Websocket()
const now = useNow()
const { ready, start } = useTimeout(1000, { controls: true })

socket.connect()   

socket.onmessage(async (data: Progress) => {
    dlentries.value[data.id].downloadedChunks[data.index] = data.downloaded
    dlentries.value[data.id].status = 'Downloading'

    const downloadedTotal = dlentries.value[data.id]
        .downloadedChunks.reduce((total, chunk) => total + chunk)
    
    // calculate the total downloaded percentage
    const size = dlentries.value[data.id].size
    dlentries.value[data.id].progress = (downloadedTotal / size) * 100

    // refresh calculation every second
    if (ready.value) {
        // calculcate the download speed
        const elapsedSecond = new Date(
                now.value.getTime() - new Date(dlentries.value[data.id].date).getTime()
            ).getTime() / 1000

        const speed = downloadedTotal / elapsedSecond
        dlentries.value[data.id].speed = speed

        // calculate time left
        const remainingSize = size - downloadedTotal
        dlentries.value[data.id].timeLeft = remainingSize / speed
        
        start()
    }

    if (data.done) {
        dlentries.value[data.id].status = 'Completed'
        dlentries.value[data.id].timeLeft = 0
        dlentries.value[data.id].progress = 100

        await entries.update(dlentries.value[data.id])
    }
})

onMounted(async () => {
    dlentries.value = await entries.all()
})

onUnmounted(async () => {
    await entries.updateAll(dlentries.value)
})


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
                        <download-dialog @fetched="res => dlentries[res.id] = res" /> 
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

        <download-list class="w-full" :items="items"/>
    </div>
</template>
