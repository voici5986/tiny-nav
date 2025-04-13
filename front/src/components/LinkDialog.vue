<template>
    <div v-if="show" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4">
        <div class="bg-white dark:bg-gray-800 text-gray-800 dark:text-gray-100 rounded-lg p-6 max-w-lg w-full">
            <h3 class="text-lg font-medium mb-4">
                {{ mode === 'add' ? '添加链接' : '编辑链接' }}
            </h3>

            <form @submit.prevent="handleSubmit" class="space-y-4">
                <!-- 名称输入 -->
                <div>
                    <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        名称
                    </label>
                    <input v-model="formData.name" type="text" required
                        class="w-full border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-800 dark:text-gray-100 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400" />
                </div>

                <!-- URL输入 -->
                <div>
                    <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        URL
                    </label>
                    <input v-model="formData.url" type="url" required
                        class="w-full border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-800 dark:text-gray-100 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400" />
                </div>

                <!-- 图标相关输入 -->
                <div class="w-full max-w-full overflow-hidden">
                    <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        图标
                    </label>
                    <div class="flex space-x-2">
                        <input v-model="formData.icon" type="text" placeholder="图标URL"
                            class="flex-1 max-w-[calc(100%-3.5rem)] border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-800 dark:text-gray-100 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400" />
                        <button type="button" @click="fetchIcon" :disabled="isLoading"
                            class="shrink-0 px-4 py-2 bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-300 rounded-md hover:bg-gray-200 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400 disabled:opacity-50 disabled:cursor-not-allowed">
                            <div class="i-mdi-image-sync-outline"></div>
                        </button>
                    </div>
                </div>

                <!-- 分类输入（使用 Combobox） -->
                <div>
                    <Combobox v-model="formData.category">
                        <div class="relative">
                            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                                分类
                            </label>
                            <div
                                class="relative w-full cursor-default overflow-hidden rounded-md bg-white dark:bg-gray-700 text-left border border-gray-300 dark:border-gray-600 focus:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 dark:focus-visible:ring-blue-400">
                                <ComboboxInput
                                    class="w-full border-none px-3 py-2 text-gray-900 dark:text-gray-100 bg-white dark:bg-gray-700 focus:outline-none"
                                    :displayValue="(category: unknown) => category as string"
                                    @change="handleCategoryInput" />
                                <ComboboxButton class="absolute inset-y-0 right-0 flex items-center pr-2">
                                    <div class="i-mdi-chevron-up-down h-5 w-5 text-gray-400 dark:text-gray-300"
                                        aria-hidden="true" />
                                </ComboboxButton>
                            </div>
                            <TransitionRoot leave="transition ease-in duration-100" leaveFrom="opacity-100"
                                leaveTo="opacity-0" @after-leave="query = ''">
                                <ComboboxOptions
                                    class="absolute mt-1 max-h-60 w-full overflow-auto rounded-md bg-white dark:bg-gray-700 py-1 text-base shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none sm:text-sm">
                                    <ComboboxOption v-if="filteredCategories.length === 0 && query !== ''"
                                        :value="query" v-slot="{ active }">
                                        <li class="relative cursor-default select-none py-2 px-4" :class="{
                                            'bg-blue-500 text-white': active,
                                            'text-gray-900 dark:text-gray-100': !active
                                        }">
                                            创建新分类 "{{ query }}"
                                        </li>
                                    </ComboboxOption>

                                    <ComboboxOption v-for="category in filteredCategories" :key="category"
                                        :value="category" as="template" v-slot="{ selected, active }">
                                        <li class="relative cursor-default select-none py-2 px-4" :class="{
                                            'bg-blue-500 text-white': active,
                                            'text-gray-900 dark:text-gray-100': !active
                                        }">
                                            <span class="block truncate"
                                                :class="{ 'font-medium': selected, 'font-normal': !selected }">
                                                {{ category }}
                                            </span>
                                        </li>
                                    </ComboboxOption>
                                </ComboboxOptions>
                            </TransitionRoot>
                        </div>
                    </Combobox>
                </div>

                <!-- 按钮组 -->
                <div class="flex justify-end space-x-4 mt-6">
                    <button type="button" @click="$emit('update:show', false)"
                        class="px-4 py-2 text-gray-600 dark:text-gray-300 rounded-md hover:text-gray-800 dark:hover:text-gray-100 focus:outline-none">
                        取消
                    </button>
                    <button type="submit"
                        class="px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 dark:hover:bg-blue-400 focus:outline-none focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400">
                        确定
                    </button>
                </div>
            </form>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import type { Link } from '@/api/types'
import { api } from '@/api'
import {
    Combobox,
    ComboboxInput,
    ComboboxButton,
    ComboboxOptions,
    ComboboxOption,
    TransitionRoot,
} from '@headlessui/vue'

const props = defineProps<{
    show: boolean
    mode: 'add' | 'update'
    link?: Link
    categories: string[] // 新增 categories prop
}>()

const emit = defineEmits<{
    (e: 'update:show', value: boolean): void
    (e: 'submit', value: Link): void
}>()

const query = ref('')
const isLoading = ref(false)

const formData = ref({
    name: '',
    url: '',
    icon: '',
    category: '',
    sortIndex: 0,
})

// 处理分类输入
const handleCategoryInput = (event: Event) => {
    const value = (event.target as HTMLInputElement).value
    query.value = value
    // 直接更新 formData 的 category
    formData.value.category = value
}

// 根据输入过滤分类
const filteredCategories = computed(() => {
    if (query.value === '') return props.categories

    return props.categories.filter((category) =>
        category.toLowerCase().includes(query.value.toLowerCase())
    )
})

// 当 link 属性改变时更新表单数据
watch(() => props.link, (newLink) => {
    if (newLink) {
        formData.value = { ...newLink }
    } else {
        formData.value = {
            name: '',
            url: '',
            icon: '',
            category: '',
            sortIndex: 0,
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

// 提交表单
const handleSubmit = () => {
    emit('submit', { ...formData.value })
    emit('update:show', false)
}
</script>
