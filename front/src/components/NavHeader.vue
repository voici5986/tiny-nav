<template>
    <div class="bg-white dark:bg-gray-800 shadow dark:shadow-gray-900">
        <div class="max-w-7xl mx-auto px-4 py-3 flex justify-between items-center">
            <!-- 左侧图标和标题 -->
            <div class="flex items-center gap-2">
                <img src="/icon.svg" alt="网站图标" class="w-8 h-8" />
                <h1 class="text-xl font-bold text-gray-800 dark:text-gray-100">我的导航</h1>
            </div>

            <!-- 右侧菜单按钮 -->
            <div class="flex items-center gap-2">
                <button @click="toggleMobileMenu"
                    class="md:hidden p-2 w-10 h-10 rounded-md bg-gray-100 dark:bg-gray-800 hover:bg-gray-200 dark:hover:bg-gray-700">
                    <div :class="[isMobileMenuOpen ? 'i-mdi-close text-red-500' : 'i-mdi-menu text-gray-500 dark:text-gray-300']"
                        class="text-2xl"></div>
                </button>
            </div>

            <!-- PC端按钮组 -->
            <div class="hidden md:flex items-center gap-4">
                <button @click="toggleEditMode"
                    class="flex items-center gap-2 px-3 py-2 rounded-md bg-gray-100 dark:bg-gray-800 hover:bg-gray-200 dark:hover:bg-gray-700"
                    :title="editMode ? '浏览模式' : '编辑模式'">
                    <div
                        :class="[editMode ? 'i-mdi-eye text-blue-500 dark:text-blue-300' : 'i-mdi-pencil text-gray-400 dark:text-gray-300']">
                    </div>
                </button>
                <button @click="$emit('add')"
                    class="flex items-center gap-2 px-3 py-2 rounded-md bg-gray-100 dark:bg-gray-800 hover:bg-gray-200 dark:hover:bg-gray-700"
                    title="添加网站">
                    <div class="i-mdi-plus-circle text-gray-400 dark:text-gray-300">
                    </div>
                </button>
                <button @click="themeStore.toggleTheme"
                    class="flex items-center gap-2 px-3 py-2 rounded-md bg-gray-100 dark:bg-gray-800 hover:bg-gray-200 dark:hover:bg-gray-700"
                    :title="themeStore.isDarkTheme ? '浅色模式' : '深色模式'">
                    <div
                        :class="themeStore.isDarkTheme ? 'i-mdi-white-balance-sunny text-blue-500 dark:text-blue-300' : 'i-mdi-moon-waxing-crescent text-gray-400 dark:text-gray-300'">
                    </div>
                </button>

                <button v-if="!isNoAuthMode" @click="$emit('logout')"
                    class="flex items-center gap-2 px-3 py-2 rounded-md bg-gray-100 dark:bg-gray-800 hover:bg-gray-200 dark:hover:bg-gray-700"
                    title="退出登录">
                    <div class="i-mdi-logout text-gray-400 dark:text-gray-300">
                    </div>
                </button>
            </div>
        </div>

        <!-- 移动端菜单 -->
        <div v-if="isMobileMenuOpen" class="md:hidden bg-white dark:bg-gray-800 shadow-lg dark:shadow-gray-900">
            <div class="flex flex-col gap-2 p-4">
                <button @click="toggleEditMode"
                    class="flex items-center gap-3 px-2 py-2 rounded-md bg-gray-100 dark:bg-gray-800 hover:bg-gray-200 dark:hover:bg-gray-700">
                    <div
                        :class="[editMode ? 'i-mdi-eye text-blue-500 dark:text-blue-300' : 'i-mdi-pencil text-gray-400 dark:text-gray-300']">
                    </div>
                    {{ editMode ? '浏览模式' : '编辑模式' }}
                </button>
                <button @click="$emit('add')"
                    class="flex items-center gap-3 px-3 py-2 rounded-md bg-gray-100 dark:bg-gray-800 hover:bg-gray-200 dark:hover:bg-gray-700">
                    <div class="i-mdi-plus-circle text-gray-400 dark:text-gray-300">
                    </div>
                    添加网站
                </button>
                <button @click="themeStore.toggleTheme"
                    class="flex items-center gap-3 px-3 py-2 rounded-md bg-gray-100 dark:bg-gray-800 hover:bg-gray-200 dark:hover:bg-gray-700">
                    <div
                        :class="themeStore.isDarkTheme ? 'i-mdi-white-balance-sunny text-blue-500 dark:text-blue-300' : 'i-mdi-moon-waxing-crescent text-gray-400 dark:text-gray-300'">
                    </div>
                    {{ themeStore.isDarkTheme ? '浅色模式' : '深色模式' }}
                </button>
                <button v-if="!isNoAuthMode" @click="$emit('logout')"
                    class="flex items-center gap-3 px-3 py-2 rounded-md bg-gray-100 dark:bg-gray-800 hover:bg-gray-200 dark:hover:bg-gray-700">
                    <div class="i-mdi-logout text-gray-400 dark:text-gray-300">
                    </div>
                    退出登录
                </button>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useThemeStore } from '@/stores/themeStore'

const themeStore = useThemeStore()

const props = defineProps<{
    editMode: boolean
    isNoAuthMode: boolean
}>()

const emit = defineEmits<{
    (e: 'update:editMode', value: boolean): void
    (e: 'add'): void
    (e: 'logout'): void
}>()

const toggleEditMode = () => {
    emit('update:editMode', !props.editMode)
}

// 移动端菜单状态
const isMobileMenuOpen = ref(false)
const toggleMobileMenu = () => {
    isMobileMenuOpen.value = !isMobileMenuOpen.value
}

</script>
