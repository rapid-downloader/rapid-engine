<script setup lang="ts">
import { Status } from '../types'
import { statusColor } from '@/lib/parse'

defineProps<{
    status: Status
}>()
</script>

<template>
    <div v-if="status === 'Completed' ">
        <i-fluent-checkmark-circle-16-regular :class="`group-hover:hidden mx-auto ${statusColor(status)}`" />
        <i-fluent-checkmark-circle-16-filled :class="`hidden group-hover:block mx-auto ${statusColor(status)}`" />
    </div>
    <div v-if="status === 'Downloading' ">
        <i-fluent-circle-hint-half-vertical-16-regular :class="`group-hover:hidden mx-auto ${statusColor(status)}`" />
        <i-fluent-circle-hint-half-vertical-16-filled :class="`hidden group-hover:block mx-auto ${statusColor(status)}`" />
    </div>
    <div v-if="status === 'Failed' ">
        <i-fluent-dismiss-circle-16-regular :class="`group-hover:hidden mx-auto ${statusColor(status)}`" />
        <i-fluent-dismiss-circle-16-filled :class="`hidden group-hover:block mx-auto ${statusColor(status)}`" />
    </div>
    <div v-if="status === 'Paused' ">
        <i-fluent-pause-circle-20-regular :class="`group-hover:hidden mx-auto ${statusColor(status)}`" />
        <i-fluent-pause-circle-20-filled :class="`hidden group-hover:block mx-auto ${statusColor(status)}`" />
    </div>
    <div v-if="status === 'Queued' ">
        <i-fluent-people-queue-20-regular :class="`group-hover:hidden mx-auto ${statusColor(status)}`" />
        <i-fluent-people-queue-20-filled :class="`hidden group-hover:block mx-auto ${statusColor(status)}`" />
    </div>
    <div v-if="status === 'Stoped' ">
        <i-fluent-subtract-circle-16-regular :class="`group-hover:hidden mx-auto ${statusColor(status)}`" />
        <i-fluent-subtract-circle-16-filled :class="`hidden group-hover:block mx-auto ${statusColor(status)}`" />
    </div>
</template>