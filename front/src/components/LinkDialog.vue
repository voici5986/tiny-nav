<template>
    <Dialog :open="show" @close="closeModal" class="relative z-10">
        <div class="fixed inset-0 bg-black bg-opacity-25" />
        <div class="fixed inset-0 overflow-y-auto">
            <div class="flex min-h-full items-center justify-center p-4">
                <DialogPanel
                    class="w-full max-w-md transform overflow-hidden rounded-lg bg-white p-6 shadow-xl transition-all">
                    <DialogTitle class="text-lg font-medium mb-4">
                        {{ mode === 'add' ? '添加链接' : '修改链接' }}
                    </DialogTitle>
                    <form @submit.prevent="handleSubmit" class="space-y-4">
                        <div>
                            <label class="block text-sm font-medium text-gray-700">名称</label>
                            <input v-model="formData.name" type="text"
                                class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-300 focus:ring focus:ring-blue-200"
                                required />
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-gray-700">URL</label>
                            <input v-model="formData.url" type="url"
                                class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-300 focus:ring focus:ring-blue-200"
                                required />
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-gray-700">分类</label>
                            <input v-model="formData.category" type="text"
                                class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-300 focus:ring focus:ring-blue-200"
                                required />
                        </div>
                        <div class="flex justify-end space-x-4 mt-6">
                            <button type="button" @click="closeModal"
                                class="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500">
                                取消
                            </button>
                            <button type="submit"
                                class="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">
                                {{ mode === 'add' ? '添加' : '保存' }}
                            </button>
                        </div>
                    </form>
                </DialogPanel>
            </div>
        </div>
    </Dialog>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import {
    Dialog,
    DialogPanel,
    DialogTitle,
} from '@headlessui/vue'
import type { Link } from '@/api/types'
import { api } from '@/api'

interface Props {
    show: boolean
    mode: 'add' | 'update'
    link?: Link
}

const props = defineProps<Props>()
const emit = defineEmits<{
    (e: 'close'): void
    (e: 'update:show', value: boolean): void
    (e: 'submit', link: Link): void
}>()

const formData = ref<Link>({
    name: '',
    url: '',
    icon: '',
    category: ''
})

watch(() => props.link, (newLink) => {
    if (newLink) {
        formData.value = { ...newLink }
    }
}, { immediate: true })

const closeModal = () => {
    emit('update:show', false)
    emit('close')
}

const handleSubmit = async () => {
    try {
        // 获取网站图标
        const iconRes = await api.getWebsiteIcon(formData.value.url)
        formData.value.icon = iconRes.iconUrl

        emit('submit', formData.value)
        closeModal()
    } catch (error) {
        console.error('Failed to submit:', error)
        alert('操作失败')
    }
}

// 当模态框打开时重置表单
watch(() => props.show, (newShow) => {
    if (newShow && props.mode === 'add') {
        formData.value = {
            name: '',
            url: '',
            icon: '',
            category: ''
        }
    }
})
</script>
