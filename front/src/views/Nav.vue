<template>
    <div>
        <div class="top-bar bg-gray-800 text-white p-4 flex justify-between items-center">
            <div class="flex items-center">
                <h1 class="text-2xl font-bold">导航</h1>
                <nav class="ml-4">
                    <Switch v-model="editMode"
                        class="relative inline-flex items-center h-6 rounded-full w-11 bg-gray-200"
                        :class="[editMode ? 'bg-blue-600' : 'bg-gray-200']">
                        <span class="sr-only">编辑模式</span>
                        <span class="inline-block w-4 h-4 transform bg-white rounded-full transition-transform"
                            :class="[editMode ? 'translate-x-6' : 'translate-x-1']" />
                    </Switch>
                    <span class="ml-2">{{ editMode ? '编辑' : '浏览' }}</span>
                    <button v-if="editMode" @click="openAddDialog"
                        class="ml-4 bg-red-500 text-white px-4 py-2 rounded-md hover:bg-red-600 focus:outline-none focus:ring focus:ring-red-300">
                        新增链接
                    </button>
                </nav>
            </div>
            <button @click="handleLogout"
                class="bg-red-500 text-white px-4 py-2 rounded-md hover:bg-red-600 focus:outline-none focus:ring focus:ring-red-300">
                退出
            </button>
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
import { Dialog, DialogPanel, DialogTitle, Switch } from '@headlessui/vue'
import LinkCard from '@/components/LinkCard.vue'
import LinkDialog from '@/components/LinkDialog.vue'

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

onMounted(() => {
    if (!store.token) {
        router.push('/')
        return
    }
    fetchLinks()
})
</script>
