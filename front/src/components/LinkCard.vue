<template>
    <div class="relative group">
        <!-- 整个卡片作为可点击区域 -->
        <a :href="link.url" target="_blank" rel="noopener noreferrer"
            class="block p-4 bg-white rounded-lg shadow-md hover:shadow-lg transition-all">
            <div class="flex flex-col items-center">
                <div class="mb-2 w-12 h-12">
                    <!-- SVG 代码 -->
                    <div v-if="isSvgContent(link.icon)" v-html="link.icon"
                        class="w-full h-full flex items-center justify-center"></div>
                    <!-- 图片 URL -->
                    <img v-else :src="getIconUrl(link.icon)" :alt="link.name" class="w-full h-full object-contain"
                        @error="onImageError" />
                </div>

                <div class="w-full px-2"> <!-- 添加容器控制宽度 -->
                    <span class="text-blue-500 text-center mt-2 text-sm block truncate" :title="link.name">
                        {{ link.name }}
                    </span>
                </div>
            </div>
        </a>

        <!-- 悬浮按钮容器 -->
        <div v-if="editMode"
            class="absolute inset-x-0 top-1/2 -translate-y-1/2 flex justify-between px-2 opacity-50 duration-200">
            <!-- 修改按钮 -->
            <button @click.prevent="$emit('update', link.globalIndex)"
                class="w-8 h-8 flex items-center justify-center text-gray-400 hover:text-yellow-500 bg-white rounded-full shadow-md transition-all hover:scale-110"
                title="修改">
                <div class="i-mdi-pencil text-xl"></div>
            </button>

            <!-- 拖动图标 -->
            <button
                class="w-8 h-8 flex items-center justify-center text-gray-400 hover:text-blue-500 rounded-full shadow-md transition-all hover:scale-110 drag-handle"
                title="拖动">
                <div class="i-mdi-drag text-xl"></div>
            </button>


            <!-- 删除按钮 -->
            <button @click.prevent="$emit('delete', link.globalIndex)"
                class="w-8 h-8 flex items-center justify-center text-gray-400 hover:text-red-500 bg-white rounded-full shadow-md transition-all hover:scale-110"
                title="删除">
                <div class="i-mdi-delete text-xl"></div>
            </button>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { Link } from '@/api/types'

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

const isSvgContent = (icon: string) => {
    if (!icon) return false
    const iconLower = icon.toLowerCase()
    return iconLower.startsWith('<svg') || iconLower.includes('xmlns="http://www.w3.org/2000/svg"')
}

const getIconUrl = (icon: string) => {
    if (!icon) {
        const svg = defaultIcon.value
        return svg ? `data:image/svg+xml;charset=utf-8,${encodeURIComponent(svg.outerHTML)}` : ''
    }

    if (icon.startsWith('http')) {
        return icon
    }

    return icon
}

const onImageError = (e: Event) => {
    const img = e.target as HTMLImageElement
    const svg = defaultIcon.value
    if (svg) {
        img.src = `data:image/svg+xml;charset=utf-8,${encodeURIComponent(svg.outerHTML)}`
        img.onerror = null
    }
}
</script>

<style scoped>
/* 确保 SVG 图标居中且大小合适 */
:deep(svg) {
    width: 100%;
    height: 100%;
}

/* 按钮容器始终保持在卡片上方 */
.group:hover .absolute {
    pointer-events: auto;
}

/* 卡片悬浮效果 */
.group:hover a {
    transform: translateY(-2px);
}

/* 文本截断样式 */
.truncate {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 100%;
    margin: 0 auto;
}
</style>
