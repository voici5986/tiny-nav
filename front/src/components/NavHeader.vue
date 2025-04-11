<template>
    <div class="bg-white shadow">
        <div class="max-w-7xl mx-auto px-4 py-3 flex justify-between items-center">
            <!-- 左侧菜单按钮 -->
            <div class="flex items-center gap-2">
                <button @click="toggleMobileMenu" class="md:hidden p-2 rounded-md hover:bg-gray-100">
                    <div class="i-mdi-menu text-2xl text-gray-500"></div>
                </button>
                <img src="/icon.svg" alt="网站图标" class="w-8 h-8" />
                <h1 class="text-xl font-bold text-gray-800">我的导航</h1>
            </div>

            <!-- PC端按钮组 -->
            <div class="hidden md:flex items-center gap-4">
                <button @click="toggleEditMode" class="flex items-center gap-2 px-3 py-2 rounded-md hover:bg-gray-100"
                    :title="editMode ? '浏览模式' : '编辑模式'">
                    <div :class="[editMode ? 'i-mdi-eye text-blue-500' : 'i-mdi-pencil text-gray-400']"></div>
                </button>
                <button @click="$emit('add')" class="flex items-center gap-2 px-3 py-2 rounded-md hover:bg-gray-100"
                    title="添加网站">
                    <div class="i-mdi-plus-circle text-gray-400"></div>
                </button>
                <button @click="toggleTheme" class="flex items-center gap-2 px-3 py-2 rounded-md hover:bg-gray-100"
                    :title="isDarkTheme ? '浅色模式' : '深色模式'">
                    <div class="i-mdi-theme-light-dark text-gray-400"></div>
                </button>
                <button @click="$emit('logout')" class="flex items-center gap-2 px-3 py-2 rounded-md hover:bg-gray-100"
                    title="退出登录">
                    <div class="i-mdi-logout text-gray-400"></div>
                </button>
            </div>
        </div>

        <!-- 移动端菜单 -->
        <div v-if="isMobileMenuOpen" class="md:hidden bg-white shadow-lg">
            <div class="flex flex-col gap-2 p-4">
                <button @click="toggleEditMode" class="flex items-center px-3 py-2 rounded-md hover:bg-gray-100">
                    <div :class="[editMode ? 'i-mdi-eye text-blue-500' : 'i-mdi-pencil text-gray-400']"></div>
                    {{ editMode ? '浏览模式' : '编辑模式' }}
                </button>
                <button @click="$emit('add')" class="flex items-center px-3 py-2 rounded-md hover:bg-gray-100">
                    <div class="i-mdi-plus-circle text-gray-400"></div>
                    添加网站
                </button>
                <button @click="toggleTheme" class="flex items-center px-3 py-2 rounded-md hover:bg-gray-100">
                    <div class="i-mdi-theme-light-dark text-gray-400"></div>
                    切换到{{ isDarkTheme ? '浅色模式' : '深色模式' }}
                </button>
                <button @click="$emit('logout')" class="flex items-center px-3 py-2 rounded-md hover:bg-gray-100">
                    <div class="i-mdi-logout text-gray-400"></div>
                    退出登录
                </button>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const props = defineProps<{
    editMode: boolean
}>()

const emit = defineEmits<{
    (e: 'update:editMode', value: boolean): void
    (e: 'add'): void
    (e: 'logout'): void
}>()

const toggleEditMode = () => {
    emit('update:editMode', !props.editMode)
}

// 主题切换
const isDarkTheme = ref(false)
const toggleTheme = () => {
    isDarkTheme.value = !isDarkTheme.value
    document.documentElement.classList.toggle('dark', isDarkTheme.value)
}

// 移动端菜单状态
const isMobileMenuOpen = ref(false)
const toggleMobileMenu = () => {
    isMobileMenuOpen.value = !isMobileMenuOpen.value
}
</script>
