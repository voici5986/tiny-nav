<template>
    <div v-if="show" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4">
        <div class="bg-white rounded-lg p-6 max-w-lg w-full">
            <h3 class="text-lg font-medium mb-4">
                {{ mode === 'add' ? '添加链接' : '编辑链接' }}
            </h3>

            <form @submit.prevent="handleSubmit" class="space-y-4">
                <!-- 名称输入 -->
                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">
                        名称
                    </label>
                    <input v-model="formData.name" type="text" required
                        class="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500" />
                </div>

                <!-- URL输入 -->
                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">
                        URL
                    </label>
                    <input v-model="formData.url" type="url" required
                        class="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500" />
                </div>

                <!-- 图标相关输入 -->
                <div class="w-full max-w-full overflow-hidden">
                    <label class="block text-sm font-medium text-gray-700 mb-1">
                        图标
                    </label>
                    <div class="flex space-x-2">
                        <input v-model="formData.icon" type="text" placeholder="图标URL"
                            class="flex-1 max-w-[calc(100%-3.5rem)] border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500" />
                        <button type="button" @click="fetchIcon" :disabled="isLoading"
                            class="shrink-0 px-4 py-2 bg-gray-100 text-gray-700 rounded-md hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed">
                            <div class="i-mdi-image-sync-outline"></div>
                        </button>
                    </div>
                </div>

                <!-- 分类输入 -->
                <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">
                        分类
                    </label>
                    <input v-model="formData.category" type="text" required
                        class="w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500" />
                </div>

                <!-- 按钮组 -->
                <div class="flex justify-end space-x-4 mt-6">
                    <button type="button" @click="$emit('update:show', false)"
                        class="px-4 py-2 text-gray-600 hover:text-gray-800 focus:outline-none">
                        取消
                    </button>
                    <button type="submit"
                        class="px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500">
                        确定
                    </button>
                </div>
            </form>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import type { Link } from '@/api/types'
import { api } from '@/api'

const props = defineProps<{
    show: boolean
    mode: 'add' | 'update'
    link?: Link
}>()

const emit = defineEmits<{
    (e: 'update:show', value: boolean): void
    (e: 'submit', value: Link): void
}>()

const formData = ref({
    name: '',
    url: '',
    icon: '',
    category: ''
})

const isLoading = ref(false)

// 当 link 属性改变时更新表单数据
watch(() => props.link, (newLink) => {
    if (newLink) {
        formData.value = { ...newLink }
    } else {
        formData.value = {
            name: '',
            url: '',
            icon: '',
            category: ''
        }
    }
}, { immediate: true })

// 获取图标
const fetchIcon = async () => {
    if (!formData.value.url || isLoading.value) return

    isLoading.value = true
    try {
        console.log("url", formData.value.url)
        const icon = await api.getWebsiteIcon(formData.value.url)
        console.log("icon", icon)
        formData.value.icon = icon
    } catch (error) {
        console.error('Error fetching icon:', error)
        alert('获取图标失败')
    } finally {
        isLoading.value = false
    }
}

// 处理图标加载错误
const handleIconError = () => {
    formData.value.icon = '' // 清空无效的图标URL
}

// 提交表单
const handleSubmit = () => {
    emit('submit', { ...formData.value })
    emit('update:show', false)
}
</script>
