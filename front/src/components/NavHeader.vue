<template>
    <div class="bg-white shadow">
        <div class="max-w-5xl mx-auto px-4 py-3 flex justify-between items-center">
            <h1 class="text-xl font-bold">我的导航</h1>

            <!-- 右侧菜单 -->
            <Menu as="div" class="relative z-100">
                <div>
                    <MenuButton
                        class="flex items-center space-x-2 px-3 py-2 rounded-md hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">
                        <span class="text-sm text-gray-700">菜单</span>
                        <svg class="w-5 h-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                        </svg>
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
                                'group flex w-full items-center px-4 py-2 text-sm text-gray-700'
                            ]">
                                <svg class="mr-3 h-5 w-5 text-gray-400 group-hover:text-gray-500"
                                    :class="{ 'text-blue-500': editMode }" fill="none" stroke="currentColor"
                                    viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                        d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
                                </svg>
                                {{ editMode ? '浏览模式' : '编辑模式' }}
                            </button>
                            </MenuItem>

                            <!-- 添加链接 -->
                            <MenuItem v-slot="{ active }" as="template">
                            <button @click="$emit('add')" :class="[
                                active ? 'bg-gray-100' : '',
                                'group flex w-full items-center px-4 py-2 text-sm text-gray-700'
                            ]">
                                <svg class="mr-3 h-5 w-5 text-gray-400 group-hover:text-gray-500" fill="none"
                                    stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                        d="M12 4v16m8-8H4" />
                                </svg>
                                添加网站
                            </button>
                            </MenuItem>
                        </div>

                        <!-- 退出登录 -->
                        <div class="py-1">
                            <MenuItem v-slot="{ active }" as="template">
                            <button @click="$emit('logout')" :class="[
                                active ? 'bg-gray-100' : '',
                                'group flex w-full items-center px-4 py-2 text-sm text-gray-700'
                            ]">
                                <svg class="mr-3 h-5 w-5 text-gray-400 group-hover:text-gray-500" fill="none"
                                    stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                        d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                                </svg>
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
