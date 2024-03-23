<script setup lang="ts">
import { Button } from '@/components/ui/button'
import Header from '@/components/Header.vue';
import DownloadList from './components/DownloadList.vue';
import DownloadListSkeleton from './components/DownloadListSkeleton.vue';
import { computed, watch, onUnmounted, ref, onMounted } from 'vue';
import { Tooltip } from '@/components/ui/tooltip';
import Filter from './components/Filter.vue';
import Entries from './api'
import Dialog from '@/components/ui/dialog/XDialog.vue';
import DownloadDialog from '../download/components/DownloadDialog.vue';
import { useNow, useTimeout } from '@vueuse/core';
//@ts-ignore
import { EventsOn } from '@/../wailsjs/runtime'
import { Download } from './types';
import { Downloader } from '../download/api';

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
const downloader = Downloader()

const page = ref(1)
onMounted(async () => {
    dlentries.value = await entries.all(page.value)
    loading.value = false
})

watch(page, async p => {
    const nextEntries = await entries.all(p)

    dlentries.value = {
        ...dlentries.value,
        ...nextEntries,
    }
})

// onUnmounted(async () => {
//     await entries.updateAll(dlentries.value)
// })

const dialogOpen = ref(false)
async function fetched(result: Download) {
    dlentries.value[result.id] = result
    await entries.update(result)
    dialogOpen.value = false
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
    downloaded?: number
    size?: number
    progress?: number
    length: number
    done: boolean
}

const now = useNow()
const { ready, start } = useTimeout(1000, { controls: true })

const singo = ref()

async function update(progress: Progress) {
    singo.value = progress
    if (!dlentries.value[progress.id].downloadedChunks) {
        dlentries.value[progress.id].downloadedChunks = new Array<number>(progress.length)
    }

    if (progress.downloaded) {
        dlentries.value[progress.id].downloadedChunks[progress.index] = progress.downloaded
        dlentries.value[progress.id].status = 'Downloading'
    }

    const downloadedTotal = dlentries.value[progress.id]
        .downloadedChunks.reduce((total, chunk) => total + chunk)
    
    // // calculate the total downloaded percentage
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
        dlentries.value[progress.id].progress = (downloadedTotal / size) * 100
        if (dlentries.value[progress.id].progress == 100) {
            dlentries.value[progress.id].status =  'Completed'
            dlentries.value[progress.id].timeLeft = 0
            
            await entries.update(dlentries.value[progress.id])
        }
    }
}

EventsOn('progress', async (...event: any) => {
    update(event[0] as Progress)
})

async function removeEntry(id: string, fromDisk: boolean) {
    delete dlentries.value[id]
    await entries.deleteEntry(id, fromDisk)
}

async function stopDownload(id: string) {
    await downloader.stop(id)

    const entry = dlentries.value[id]
    entry.status = 'Stoped'
    await entries.update(dlentries.value[id])
}

async function pauseDownload(id: string) {
    await downloader.pause(id)

    dlentries.value[id].status = 'Paused'
    await entries.update(dlentries.value[id])
}

async function resumeDownload(id: string) {
    await downloader.resume(id)

    const entry = dlentries.value[id]
    entry.status = 'Downloading'
    await entries.update(dlentries.value[id])
}

async function restartDownload(id: string) {
    await downloader.restart(id)

    const entry = dlentries.value[id]
    entry.status = 'Downloading'
    await entries.update(dlentries.value[id])
}
</script>

<template>
    <Header>
        <div class="flex gap-3">
            <Tooltip text="New Download" location="bottom">
                <Dialog v-model:open="dialogOpen" title="New Download" description="Provide a link to start a new download">
                    <Button class="flex gap-2 bg-accent justify-center hover:bg-accent/90">
                        <i-fluent-add-16-filled class="text-accent-foreground" />
                    </Button>

                    <template v-slot:content>
                        <download-dialog @fetched="fetched" @close="dialogOpen = false"/> 
                    </template>
                </Dialog>
            </Tooltip>

            <Tooltip text="New batch download" location="bottom">
                <Button class="flex gap-2 border border-accent group hover:bg-accent" variant="outline">
                    <i-fluent-add-square-multiple-16-filled class="text-accent text-lg group-hover:text-accent-foreground" />
                </Button>
            </Tooltip>

            <Tooltip text="Resume all" location="bottom">
                <Button class="flex gap-2 border border-accent group hover:bg-accent" variant="outline">
                    <i-fluent-play-16-filled class="text-accent group-hover:text-accent-foreground" />
                </Button>
            </Tooltip>

            <Tooltip text="Pause all" location="bottom">
                <Button class="flex gap-2 border border-accent group hover:bg-accent" variant="outline">
                    <i-fluent-pause-16-filled class="text-accent group-hover:text-accent-foreground" />
                </Button>
            </Tooltip>
            
            <Tooltip text="Stop all" location="bottom">
                <Button class="flex gap-2 border border-accent group hover:bg-accent" variant="outline">
                    <i-fluent-stop-16-filled class="text-accent group-hover:text-accent-foreground" />
                </Button>    
            </Tooltip>
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
        <download-list v-else 
            @pause="pauseDownload"
            @resume="resumeDownload"
            @restart="restartDownload"
            @delete="removeEntry" 
            @cancel="stopDownload"
            class="w-full" 
            :items="items" 
            @paginate="page++"
            
        />
    </div>
</template>