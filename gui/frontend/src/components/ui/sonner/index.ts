import { cva } from 'class-variance-authority'
import { Component, defineComponent, h, markRaw } from 'vue'
import { ExternalToast } from 'vue-sonner/lib/types'
export { default as Toaster } from './Sonner.vue'
import { toast as callToast } from 'vue-sonner'
import Toast from './Toast.vue'

export const toastVariants = cva('', {
    variants: {
        border: {
            default: 'border-foreground/60',
            warning: 'border-warning',
            destructive: 'border-destructive',
            info: 'border-info',
            success: 'border-success',
        },
        text: {
            default: 'text-foreground',
            warning: 'text-warning',
            destructive: 'text-destructive',
            info: 'text-info',
            success: 'text-success',
        },
    },
    defaultVariants: {
      border: 'default',
      text: 'default',
    },
})

export type ToastOption = Omit<ExternalToast, 'invert' | 'icon' | 'important' | 'style' | 'unstyled' | 'descriptionClassName' | 'className' | 'promise' | 'action' > & {
    type?: NonNullable<Parameters<typeof toastVariants>[0]>['text']
    action?: {
        label?: string,
        icon?: Component,
        onClick: () => void
    }
}

export type ToastProps = ToastOption & {
    title: string
}

export function toast(title: string, data?: ToastOption) {    
    callToast(markRaw(
        defineComponent({ render: () => h(Toast, { title, ...data }) })), 
        { 
            description: undefined,
            onAutoClose: data?.onAutoClose,
            onDismiss: data?.onDismiss,
            duration: data?.duration,
            cancel: data?.cancel,
            id: data?.id || new Date().getTime().toString()
        }
    )
}