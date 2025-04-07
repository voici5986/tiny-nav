<template>
    <div>
        <!-- 顶部操作栏 -->
        <div class="bg-white shadow">
            <div class="max-w-5xl mx-auto px-4 py-3 flex justify-between items-center">
                <h1 class="text-xl font-bold">我的导航</h1>

                <!-- 右侧菜单 -->
                <Menu as="div" class="relative">
                    <div>
                        <MenuButton
                            class="flex items-center space-x-2 px-3 py-2 rounded-md hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">
                            <svg class="w-5 h-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                    d="M19 9l-7 7-7-7" />
                            </svg>
                        </MenuButton>
                    </div>

                    <transition enter-active-class="transition ease-out duration-100"
                        enter-from-class="transform opacity-0 scale-95" enter-to-class="transform opacity-100 scale-100"
                        leave-active-class="transition ease-in duration-75"
                        leave-from-class="transform opacity-100 scale-100"
                        leave-to-class="transform opacity-0 scale-95">
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
                                <button @click="openAddDialog" :class="[
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
                                <button @click="handleLogout" :class="[
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

        <div class="mx-auto max-w-5xl mt-8 p-4">
            <div v-for="(group, category) in groupedLinks" :key="category" class="mb-8">
                <h2 class="text-xl font-bold mb-4">{{ category }}</h2>
                <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-4">
                    <LinkCard v-for="link in group" :key="link.globalIndex" :link="link" :edit-mode="editMode"
                        @update="openUpdateDialog" @delete="openDeleteDialog" />
                </div>
            </div>
        </div>

        <LinkDialog v-model:show="isAddDialogOpen" :link="newLink" mode="add" @submit="handleAdd"
            @close="closeAddDialog" />

        <LinkDialog v-model:show="isUpdateDialogOpen" :link="updatedLink" mode="update" @submit="handleUpdate"
            @close="closeUpdateDialog" />

        <Dialog v-model:open="isDeleteDialogOpen" @close="closeDeleteDialog" class="relative z-10">
            <div class="fixed inset-0 bg-black bg-opacity-25" />
            <div class="fixed inset-0 overflow-y-auto">
                <div class="flex min-h-full items-center justify-center p-4">
                    <DialogPanel class="w-full max-w-md transform overflow-hidden rounded-lg bg-white p-6">
                        <DialogTitle class="text-lg font-medium mb-4">确定要删除此链接吗？</DialogTitle>
                        <div class="flex justify-end space-x-4">
                            <button @click="handleDelete"
                                class="bg-red-500 text-white px-4 py-2 rounded-md hover:bg-red-600 focus:outline-none focus:ring focus:ring-red-300">
                                确定
                            </button>
                            <button @click="closeDeleteDialog"
                                class="bg-gray-500 text-white px-4 py-2 rounded-md hover:bg-gray-600 focus:outline-none focus:ring focus:ring-gray-300">
                                取消
                            </button>
                        </div>
                    </DialogPanel>
                </div>
            </div>
        </Dialog>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useMainStore } from '@/stores'
import { api } from '@/api'
import type { Link } from '@/api/types'
import { Dialog, DialogPanel, DialogTitle } from '@headlessui/vue'
import LinkCard from '@/components/LinkCard.vue'
import LinkDialog from '@/components/LinkDialog.vue'
import { Menu, MenuButton, MenuItems, MenuItem } from '@headlessui/vue'

const router = useRouter()
const store = useMainStore()

const editMode = ref(false)
const isAddDialogOpen = ref(false)
const isUpdateDialogOpen = ref(false)
const isDeleteDialogOpen = ref(false)
const deleteIndex = ref<number | null>(null)

const newLink = ref<Link>({
    name: '',
    url: '',
    icon: '',
    category: ''
})

const updatedLink = ref<Link>({
    name: '',
    url: '',
    icon: '',
    category: ''
})

const groupedLinks = computed(() => {
    // 添加安全检查，确保 links 是数组
    const links = store.links || []
    const groups: Record<string, (Link & { globalIndex: number })[]> = {}

    console.log(links)
    if (Array.isArray(links)) {
        links.forEach((link, index) => {
            const linkWithIndex = {
                ...link,
                globalIndex: index
            }
            if (!groups[link.category]) {
                groups[link.category] = []
            }
            groups[link.category].push(linkWithIndex)
        })
    }

    return groups
})

const fetchLinks = async () => {
    try {
        const response = await api.getNavigation()
        console.log(response)
        store.setLinks(response.data.links)
    } catch (error) {
        alert('获取链接失败')
        router.push('/')
    }
}

const handleLogout = () => {
    store.logout()
    router.push('/')
}

const openAddDialog = () => {
    newLink.value = {
        name: '',
        url: '',
        icon: '',
        category: ''
    }
    isAddDialogOpen.value = true
}

const closeAddDialog = () => {
    isAddDialogOpen.value = false
}

const handleAdd = async (link: Link) => {
    try {
        await api.addLink(link)
        await fetchLinks()
        closeAddDialog()
    } catch (error) {
        alert('添加失败')
    }
}

const openUpdateDialog = (index: number) => {
    updatedLink.value = { ...store.links[index] }
    isUpdateDialogOpen.value = true
}

const closeUpdateDialog = () => {
    isUpdateDialogOpen.value = false
}

const handleUpdate = async (link: Link) => {
    try {
        await api.updateLink(store.links.findIndex(l => l.url === link.url), link)
        await fetchLinks()
        closeUpdateDialog()
    } catch (error) {
        alert('更新失败')
    }
}

const openDeleteDialog = (index: number) => {
    deleteIndex.value = index
    isDeleteDialogOpen.value = true
}

const closeDeleteDialog = () => {
    isDeleteDialogOpen.value = false
    deleteIndex.value = null
}

const handleDelete = async () => {
    if (deleteIndex.value === null) return

    try {
        await api.deleteLink(deleteIndex.value)
        await fetchLinks()
        closeDeleteDialog()
    } catch (error) {
        alert('删除失败')
    }
}

// 切换编辑模式
const toggleEditMode = () => {
    editMode.value = !editMode.value
}

onMounted(() => {
    if (!store.token) {
        router.push('/')
        return
    }
    fetchLinks()
})
</script>
