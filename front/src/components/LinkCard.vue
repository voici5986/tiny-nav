<template>
    <div class="flex flex-col items-center p-4 bg-white rounded-lg shadow-md">
        <div class="mb-2 w-12 h-12">
            <!-- SVG 代码 -->
            <div v-if="isSvgContent(link.icon)" v-html="link.icon"
                class="w-full h-full flex items-center justify-center"></div>
            <!-- 图片 URL -->
            <img v-else :src="getIconUrl(link.icon)" :alt="link.name" class="w-full h-full object-contain"
                @error="onImageError" />
            <!-- 默认图标（隐藏，仅作为参考） -->
            <svg v-show="false" ref="defaultIcon" class="w-full h-full text-gray-400" fill="none" stroke="currentColor"
                viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
            </svg>
        </div>
        <a :href="link.url" target="_blank" rel="noopener noreferrer"
            class="text-blue-500 hover:underline text-center mt-2 text-sm" :title="link.name">{{ link.name }}</a>
        <div v-if="editMode" class="mt-2 space-x-2">
            <button @click="$emit('update', link.globalIndex)"
                class="bg-yellow-400 text-white px-2 py-1 rounded-md hover:bg-yellow-500 focus:outline-none focus:ring focus:ring-yellow-300 text-sm">
                修改
            </button>
            <button @click="$emit('delete', link.globalIndex)"
                class="bg-red-500 text-white px-2 py-1 rounded-md hover:bg-red-600 focus:outline-none focus:ring focus:ring-red-300 text-sm">
                删除
            </button>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { Link } from '@/api/types'

const apiBase = import.meta.env.VITE_API_BASE
const defaultIcon = ref<SVGElement>()

interface Props {
    link: Link & { globalIndex: number }
    editMode: boolean
}

defineProps<Props>()
defineEmits<{
    (e: 'update', index: number): void
    (e: 'delete', index: number): void
}>()

// 检查是否是 SVG 内容
const isSvgContent = (icon: string) => {
    if (!icon) return false
    const iconLower = icon.toLowerCase()
    return iconLower.startsWith('<svg') || iconLower.includes('xmlns="http://www.w3.org/2000/svg"')
}

// 处理图标 URL
const getIconUrl = (icon: string) => {
    console.log(icon)
    if (!icon) {
        const svg = defaultIcon.value
        return svg ? `data:image/svg+xml;charset=utf-8,${encodeURIComponent(svg.outerHTML)}` : ''
    }

    // 如果是完整的 URL（以 http 或 https 开头）
    if (icon.startsWith('http')) {
        return icon
    }

    // 如果是相对路径（以 /data/ 开头）
    if (icon.startsWith('/data/')) {
        // 移除开头的斜杠以避免双斜杠问题
        const url = `${apiBase}${icon}`
        console.log(url)
        return url
    }

    return icon
}

// 图片加载失败时显示默认图标
const onImageError = (e: Event) => {
    const img = e.target as HTMLImageElement
    const svg = defaultIcon.value
    if (svg) {
        img.src = `data:image/svg+xml;charset=utf-8,${encodeURIComponent(svg.outerHTML)}`
        img.onerror = null // 防止死循环
    }
}
</script>

<style scoped>
/* 确保 SVG 图标居中且大小合适 */
:deep(svg) {
    width: 100%;
    height: 100%;
}

/* 链接文本超出时显示省略号 */
a {
    max-width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    display: block;
}
</style>
