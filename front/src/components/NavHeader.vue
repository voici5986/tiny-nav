<template>
    <div class="bg-white shadow">
        <div class="max-w-7xl mx-auto px-4 py-3 flex justify-between items-center">
            <div class="flex items-center gap-2">
                <img src="/icon.svg" alt="网站图标" class="w-8 h-8" />
                <h1 class="text-xl font-bold text-gray-800">我的导航</h1>
            </div>

            <!-- 右侧菜单 -->
            <Menu as="div" class="relative z-100">
                <div>
                    <MenuButton
                        class="flex items-center gap-2 px-3 py-2 rounded-md hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors">
                        <span class="text-sm text-gray-700">菜单</span>
                        <div class="i-mdi-chevron-down text-gray-500 text-xl"></div>
                    </MenuButton>
                </div>

                <transition enter-active-class="transition ease-out duration-100"
                    enter-from-class="transform opacity-0 scale-95" enter-to-class="transform opacity-100 scale-100"
                    leave-active-class="transition ease-in duration-75"
                    leave-from-class="transform opacity-100 scale-100" leave-to-class="transform opacity-0 scale-95">
                    <MenuItems
                        class="absolute right-0 mt-2 w-56 origin-top-right divide-y divide-gray-100 rounded-md bg-white shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
                        <!-- 操作菜单 -->
                        <div class="py-1">
                            <!-- 编辑模式切换 -->
                            <MenuItem v-slot="{ active }" as="template">
                            <button @click="toggleEditMode" :class="[
                                active ? 'bg-gray-100' : '',
                                'group flex w-full items-center px-4 py-2 text-sm text-gray-700 transition-colors'
                            ]">
                                <div class="mr-3 text-xl"
                                    :class="[editMode ? 'i-mdi-eye text-blue-500' : 'i-mdi-pencil text-gray-400 group-hover:text-gray-500']">
                                </div>
                                {{ editMode ? '浏览模式' : '编辑模式' }}
                            </button>
                            </MenuItem>

                            <!-- 添加链接 -->
                            <MenuItem v-slot="{ active }" as="template">
                            <button @click="$emit('add')" :class="[
                                active ? 'bg-gray-100' : '',
                                'group flex w-full items-center px-4 py-2 text-sm text-gray-700 transition-colors'
                            ]">
                                <div class="i-mdi-plus-circle mr-3 text-xl text-gray-400 group-hover:text-gray-500">
                                </div>
                                添加网站
                            </button>
                            </MenuItem>
                        </div>

                        <!-- 退出登录 -->
                        <div class="py-1">
                            <MenuItem v-slot="{ active }" as="template">
                            <button @click="$emit('logout')" :class="[
                                active ? 'bg-gray-100' : '',
                                'group flex w-full items-center px-4 py-2 text-sm text-gray-700 transition-colors'
                            ]">
                                <div class="i-mdi-logout mr-3 text-xl text-gray-400 group-hover:text-gray-500">
                                </div>
                                退出登录
                            </button>
                            </MenuItem>
                        </div>
                    </MenuItems>
                </transition>
            </Menu>
        </div>
    </div>
</template>

<script setup lang="ts">
import { Menu, MenuButton, MenuItems, MenuItem } from '@headlessui/vue'

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
</script>
