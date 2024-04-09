<script lang="ts" setup>
import Header from '@/components/Header.vue';
import H3 from '@/components/ui/H3.vue';
import SettingApi from './api'
import { onMounted, ref } from 'vue';
import { Setting } from './types';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';

const api = SettingApi()
const setting = ref<Setting>({
    DataLocation: '',
    DownloadLocation: '',
    MaxChunkCount: 8,
    MaxRetry: 3,
    MinChunkSize: 1024 * 1024 * 5
})

onMounted(async () => {
    const s = await api.get()
    if (s) setting.value = s
})

</script>

<template>
    <Header />
    <div class="mx-auto my-5 lg:max-w-screen-lg w-full p-2">
        <H3 class="font-title mb-5">Setting</H3>

        <div class="flex justify-between items-center border border-muted rounded-md p-2 bg-secondary mb-3">
            <div>
                <p class="font-medium">Download Location</p>
                <p class="text-xs opacity-80">Location where the download will be saved</p>
            </div>
            <Input disabled v-model="setting.DownloadLocation" class="max-w-xs bg-background" />
        </div>

        <div class="flex justify-between items-center border border-muted rounded-md p-2 bg-secondary mb-3">
            <div>
                <p class="font-medium">Number of Chunk</p>
                <p class="text-xs opacity-80">The number of chunks that will be downloaded concurrently at a time</p>
            </div>
            <Input v-model="setting.MaxChunkCount" type="number" class="max-w-xs bg-background" />
        </div>

        <div class="flex justify-between items-center border border-muted rounded-md p-2 bg-secondary mb-3">
            <div>
                <p class="font-medium">Chunk Size</p>
                <p class="text-xs opacity-80">Minimum size of chunk that will be split into (Bytes)</p>
            </div>
            <Input v-model="setting.MinChunkSize" type="number" class="max-w-xs bg-background" />
        </div>

        <div class="flex justify-between items-center border border-muted rounded-md p-2 bg-secondary mb-3">
            <div>
                <p class="font-medium">Retries</p>
                <p class="text-xs opacity-80">Number of retries that will be done if download process is failing</p>
            </div>
            <Input v-model="setting.MaxRetry" type="number" class="max-w-xs bg-background" />
        </div>
        
        <div class="w-full flex justify-end">
            <Button class="mt-5" size="lg">Save</Button>
        </div>
    </div>
</template>