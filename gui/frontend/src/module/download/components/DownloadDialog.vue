<script setup lang="ts">
import { computed, reactive, ref } from 'vue';
import { Label } from '@/components/ui/label';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { required, url, helpers } from '@vuelidate/validators'
import { useVuelidate } from '@vuelidate/core'
import { Downloader } from '../api'
import { Download, Request } from '@/module/home/types';
import { Progress } from '@/components/ui/progress';

const emits = defineEmits<{
    (e: 'fetched', entry: Download): void
    (e: 'close'): void
}>()

const fetchForm = reactive({
    url: ''
})

const fetchRules = {
    url: {
        required: helpers.withMessage("URL can't be empty", required),
        url: helpers.withMessage('Please provide valid URL', url),
    }
}

const downloader = Downloader()

enum State {
    Init,
    Fetching,
    Fetched,
}

const fetchValidation = useVuelidate(fetchRules, fetchForm)
const state = ref<State>(State.Init)
const result = ref<Download>()
const ext = ref('')

const cached = ref('')
async function fetch(e: Event) {
    e.preventDefault()
    state.value = State.Fetching

    if (!await fetchValidation.value.$validate()) {
        state.value = State.Init
        return
    }

    if (isCached.value) {
        state.value = State.Fetched
        return
    }
    
    result.value = await downloader.fetch({
        provider: "default",
        url: fetchForm.url
    } as Request)

    state.value = State.Fetched
    cached.value = fetchForm.url
    
    if (result.value) {
        const splitname = result.value.name.split('.')
        ext.value = `.${splitname[splitname.length-1]}`

        downloadForm.name = result.value.name.replace(ext.value, '')
        downloadForm.location = result.value.location.split('/').toSpliced(-1).join('/')
    }
}

const isCached = computed(() => {
    return cached.value && cached.value === fetchForm.url
})

const downloadForm = reactive({
    name: '',
    location: ''
})

const downloadRules = {
    name: {
        required: helpers.withMessage("File name can't be empty", required),
    },
    location: {
        required: helpers.withMessage("File location can't be empty", required),
    }
}

const downloadValidation = useVuelidate(downloadRules, downloadForm)

async function download(e: Event) {
    e.preventDefault()

    if (!await downloadValidation.value.$validate()) {
        state.value = State.Fetched
        return
    }
    
    if (result.value) {
        result.value.name = downloadForm.name + ext.value
        result.value.location = `${downloadForm.location}/${result.value.name}`

        await downloader.download(result.value.id, result.value.name, result.value.location)

        emits('fetched', result.value)
    }
}

function back(e: Event) {
    e.preventDefault()
    state.value = State.Init
}

function close(e: Event) {
    e.preventDefault()
    emits('close')
}

</script>

<template>
    <Transition name="slide-fade" mode="out-in">
        <form v-if="state === State.Fetched" @submit="download" class="flex flex-col gap-1">
            <div class="flex">
                <div class="basis-9/12">
                    <div class="flex gap-1">
                        <Input v-model="downloadForm.name" :class="`bg-secondary basis-10/12 ${downloadValidation.name.$error && 'border-destructive focus-visible:ring-destructive'}`">
                            <template v-slot:append-icon>
                                <i-fluent:edit-48-regular class="bg-secondary w-[1.75rem]" />
                            </template>
                        </Input>
                        <Input v-model="ext" class="basis-2/12 disabled:opacity-100 bg-secondary disabled:cursor-default" disabled />
                    </div>
                    <Input disabled v-model="downloadForm.location" :class="`bg-secondary mt-2 ${downloadValidation.location.$error && 'border-destructive focus-visible:ring-destructive'}`">
                        <template v-slot:append-icon>
                            <i-fluent:folder-28-regular class="bg-secondary w-[1.75rem]" />
                        </template>
                    </Input>
                </div>
                <div class="basis-3/12 my-auto pl-3">
                    <i-fluent-document-text-16-regular class="text-4xl mx-auto text-info" v-if="result!.type === 'Document'"/>
                    <i-fluent-video-clip-16-regular class="text-4xl mx-auto text-success" v-if="result!.type === 'Video'"/>
                    <i-fluent-speaker-16-regular class="text-4xl mx-auto text-accent" v-if="result!.type === 'Audio'"/>
                    <i-fluent-folder-zip-16-regular class="text-4xl mx-auto text-warning" v-if="result!.type === 'Compressed'"/>
                    <i-fluent-image-16-regular class="text-4xl mx-auto text-destructive" v-if="result!.type === 'Image'"/>
                    <i-fluent-document-16-regular class="text-4xl mx-auto" v-if="result!.type === 'Other'"/>

                    <p class="text-center text-xs mt-2 opacity-70">{{ result?.type }}</p>
                </div>
            </div>

            <div class="flex justify-between mt-5">
                <Button :disabled="downloadValidation.$error" variant="ghost" @click="back">
                    <i-radix-icons-arrow-left/>
                </Button>
                
                <div class="flex gap-2">
                    <Button variant="ghost" @click="close">
                        Cancel
                    </Button>
                    
                    <Button :disabled="downloadValidation.$error" @click="download">
                        Download
                    </Button>
                </div>
            </div>
        </form>
        <form v-else @submit="fetch" class="flex flex-col gap-1">
            <Label>URL</Label>
            <div class="flex gap-3">
                <div class="w-full relative">
                    <Input v-model="fetchForm.url" :class="`bg-secondary ${fetchValidation.$error && 'border-destructive focus-visible:ring-destructive'}`">
                        <template v-slot:append-icon>
                            <i-fluent-link-16-regular class="bg-secondary w-[1.75rem]" />
                        </template>
                    </Input>
                    <div class="absolute h-1 w-full">
                        <Progress v-if="!isCached && state === State.Fetching" indetermined class="h-0.5 rounded-full"/>
                    </div>
                    <p v-for="error in fetchValidation.url.$errors" :key="error.$uid" class="pt-1 text-xs text-destructive">{{ error.$message }}</p>
                </div>
                <Button type="submit" class="flex justify-center w-[5rem] ml-auto">
                    <i-radix-icons-arrow-right v-if="isCached"/>
                    <p v-else>Fetch</p>
                </Button>
            </div>
        </form>
    </Transition>
</template>

<style>
.slide-fade-enter-active {
  transition: all 0.2s ease-out;
}

.slide-fade-leave-active {
  transition: all 0.2s cubic-bezier(1, 0.5, 0.8, 1);
}

.slide-fade-enter-from,
.slide-fade-leave-to {
  transform: translateX(20px);
  opacity: 0;
}
</style>