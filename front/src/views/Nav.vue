<template>
    <AppLayout>
        <!-- 顶部操作栏 -->
        <NavHeader v-model:editMode="editMode" @add="openAddDialog" @logout="handleLogout" />

        <div class="mx-auto max-w-7xl mt-8 p-3">
            <div v-for="(_, category) in groupedLinks" :key="category" class="mb-8">
                <h2 class="text-xl font-bold mb-4">{{ category }}</h2>
                <draggable v-model="groupedLinks[category]" group="categories" handle=".drag-handle"
                    :itemKey="'globalIndex'" @end="onDragEnd"
                    class="grid grid-cols-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-2">
                    <!-- 使用 item 插槽传递数据 -->
                    <template #item="{ element }">
                        <LinkCard :link="element" :edit-mode="editMode" :key="element.globalIndex"
                            @update="openUpdateDialog" @delete="openDeleteDialog" />
                    </template>
                </draggable>
            </div>
        </div>

        <LinkDialog v-model:show="isAddDialogOpen" :link="newLink" :categories="existingCategories" mode="add"
            @submit="handleAdd" @close="closeAddDialog" />

        <LinkDialog v-model:show="isUpdateDialogOpen" :link="updatedLink" :categories="existingCategories" mode="update"
            @submit="handleUpdate" @close="closeUpdateDialog" />

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
    </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useMainStore } from '@/stores'
import { api } from '@/api'
import draggable from 'vuedraggable'
import type { Link } from '@/api/types'
import { Dialog, DialogPanel, DialogTitle } from '@headlessui/vue'
import LinkCard from '@/components/LinkCard.vue'
import LinkDialog from '@/components/LinkDialog.vue'
import AppLayout from '@/components/AppLayout.vue'
import NavHeader from '@/components/NavHeader.vue'

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

// 添加获取已存在分类的计算属性
const existingCategories = computed(() => {
    const categories = new Set<string>()
    store.links.forEach((link: Link) => {
        if (link.category) {
            categories.add(link.category)
        }
    })
    return Array.from(categories)
})

const groupedLinks = computed(() => {
    // 添加安全检查，确保 links 是数组
    const links = store.links || []
    const groups: Record<string, (Link & { globalIndex: number })[]> = {}

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
        const links = await api.getNavigation()
        console.log(links)
        store.setLinks(links)
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
        await api.updateLink(store.links.findIndex((l: Link) => l.url === link.url), link)
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

const onDragEnd = async (event: any) => {
    // 检查是否有拖动操作
    if (!event || !event.from || !event.to) {
        console.warn("无有效拖动数据")
        return;
    }

    // 获取拖动后的分类
    const oldGlobalIndex = event.oldIndex
    const newGlobalIndex = event.newIndex
    console.log(oldGlobalIndex, newGlobalIndex)

    // TODO:
};

onMounted(() => {
    if (!store.token) {
        router.push('/')
        return
    }
    fetchLinks()
})
</script>
