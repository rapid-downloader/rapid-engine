<script setup lang="ts">
import {
    ContextMenu,
    ContextMenuContent,
    ContextMenuItem,
    ContextMenuSeparator,
    ContextMenuShortcut,
    ContextMenuSub,
    ContextMenuSubContent,
    ContextMenuSubTrigger,
    ContextMenuTrigger,
} from '@/components/ui/context-menu'

interface Menu {
    text: string
    shortcut?: string
    seperate?: boolean
    action?: Function
    disabled?: boolean
    danger?: boolean
    submenus?: Menu[]
    props?: any
}

defineProps<{
    menus: Menu[]
}>()
</script>

<template>
    <ContextMenu>
        <ContextMenuTrigger>
            <slot />
        </ContextMenuTrigger>
        <ContextMenuContent class="w-64">
            <div v-for="menu in menus">
                <ContextMenuSub v-if="menu.submenus">
                    <ContextMenuSubTrigger inset>
                        {{ menu.text }}
                    </ContextMenuSubTrigger>
                    <ContextMenuSubContent class="w-48">
                        <ContextMenuItem v-for="sub in menu.submenus" :disabled="sub.disabled" inset :class="`group ${menu.danger ? 'focus:bg-destructive' : ''}`"
                            @select="event => sub.action?.call(event, sub.props)">
                            <p :class="`${menu.danger ? 'text-destructive group-hover:text-destructive-foreground' : ''}`">{{ sub.text }}</p>
                            <ContextMenuShortcut v-if="sub.shortcut">{{ sub.shortcut }}</ContextMenuShortcut>
                            <ContextMenuSeparator v-if="menu.seperate" />
                        </ContextMenuItem>
                    </ContextMenuSubContent>
                </ContextMenuSub>

                <ContextMenuItem v-else :disabled="menu.disabled" inset :class="`group ${menu.danger ? 'focus:bg-destructive' : ''}`" @select="event => menu.action?.call(event, menu.props)">
                    <p :class="`${menu.danger ? 'text-destructive group-hover:text-destructive-foreground' : ''}`">{{ menu.text }}</p>
                    <ContextMenuShortcut v-if="menu.shortcut">{{ menu.shortcut }}</ContextMenuShortcut>
                </ContextMenuItem>
                <ContextMenuSeparator v-if="menu.seperate" />
            </div>
        </ContextMenuContent>
    </ContextMenu>
</template>