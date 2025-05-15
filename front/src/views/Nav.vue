<template>
    <AppLayout>
        <!-- 顶部操作栏 -->
        <NavHeader v-model:editMode="editMode" @add="openAddDialog" @logout="handleLogout" @login="handleLogin" />

        <div class="mx-auto max-w-7xl mt-8 p-3">
            <draggable v-model="categories" @end="onCategoryDragEnd" item-key="category" handle=".category-drag-handle">
                <template #item="{ element: category }">
                    <div :key="category" class="mb-8">
                        <div class="flex items-center gap-2 text-xl font-bold mb-4">
                            <button v-if="editMode"
                                class="flex items-center justify-center text-gray-400 hover:text-blue-500 rounded-full shadow-md transition-all hover:scale-110 category-drag-handle"
                                title="拖动分类">
                                <div class="i-mdi-drag text-xl"></div>
                            </button>
                            <div v-else class="i-mdi-chevron-double-right text-xl text-gray-400"></div>
                            <div class="flex items-center">
                                <h2>{{ category }}</h2>
                            </div>
                        </div>

                        <!-- 链接卡片拖拽区域 -->
                        <draggable v-model="groupedLinksData[category]" group="categories" handle=".drag-handle"
                            :data-category="category" :itemKey="'globalIndex'" @end="onDragEnd($event, category)"
                            class="grid grid-cols-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-2">
                            <!-- 使用 item 插槽传递数据 -->
                            <template #item="{ element }">
                                <LinkCard :link="element" :edit-mode="editMode" :key="element.globalIndex"
                                    @update="openUpdateDialog" @delete="openDeleteDialog" />
                            </template>
                        </draggable>
                    </div>
                </template>
            </draggable>
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
import { ref, watch, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useMainStore } from '@/stores'
import { api } from '@/api'
import draggable from 'vuedraggable'
import type { Link, SortIndexUpdate } from '@/api/types'
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
    category: '',
    sortIndex: 0,
})

const updatedLink = ref<Link>({
    name: '',
    url: '',
    icon: '',
    category: '',
    sortIndex: 0,
})

// 添加获取已存在分类的计算属性
const existingCategories = computed(() => {
    const categories = new Set<string>()
    const links = store.links || []
    links.forEach((link: Link) => {
        if (link.category) {
            categories.add(link.category)
        }
    })
    return Array.from(categories)
})

const categories = ref<string[]>([])
const updateCategories = () => {
    const links = store.links || []
    // 创建当前链接中存在的分类集合
    const currentCategories = new Set<string>()
    for (const link of links) {
        if (link.category) {
            currentCategories.add(link.category)
        }
    }

    // 删除不存在的分类（保持原有顺序）
    const newCategories: string[] = []
    const storeCategories = store.categories || []

    // 保持原有顺序，只保留仍在使用的分类
    for (const category of storeCategories) {
        if (currentCategories.has(category)) {
            newCategories.push(category)
            currentCategories.delete(category) // 从当前分类集合中删除已处理的分类
        }
    }

    // 添加新的分类（将剩余的分类追加到列表末尾）
    newCategories.push(...Array.from(currentCategories))
    categories.value = newCategories
}

const groupedLinksData = ref<Record<string, (Link & { globalIndex: number })[]>>({})
// 计算分组数据
const updateGroupedLinksData = () => {
    const links = store.links || []
    const groups: Record<string, (Link & { globalIndex: number })[]> = {}

    if (Array.isArray(links)) {
        links.forEach((link, index) => {
            const linkWithIndex = {
                ...link,
                globalIndex: index,
                sortIndex: link.sortIndex || 0
            }
            if (!groups[link.category]) {
                groups[link.category] = []
            }
            groups[link.category].push(linkWithIndex)
        })

        // 对每个分组进行排序
        for (const category in groups) {
            groups[category].sort((a, b) => {
                if (a.sortIndex === 0 && b.sortIndex === 0) {
                    return a.globalIndex - b.globalIndex
                }
                if (a.sortIndex === 0) return 1
                if (b.sortIndex === 0) return -1
                return a.sortIndex - b.sortIndex
            })
        }
    }

    groupedLinksData.value = groups
}

// 监听 store.links 的变化，更新分组数据
watch(() => store.links, () => {
    updateGroupedLinksData()
}, { immediate: true })

watch(() => store.categories, () => {
    updateCategories()
}, { immediate: true })

const fetchLinks = async () => {
    try {
        const { links, categories } = await store.getNavigation()
        console.log(links, categories)
    } catch (error) {
        alert('获取链接失败')
        router.push('/')
    }
}

const handleLogout = () => {
    store.logout()
    fetchLinks()
}

const handleLogin = () => {
    router.push('/login')
}

const openAddDialog = () => {
    newLink.value = {
        name: '',
        url: '',
        icon: '',
        category: '',
        sortIndex: 0,
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

const onDragEnd = async (event: any, category: string) => {
    console.log('Drag ended:', event)
    // 如果没有实际的拖拽改变，直接返回
    if (!event.added && !event.removed && (event.oldIndex === event.newIndex && event.from === event.to)) {
        return
    }

    // 存储所有需要更新的分类
    const categoriesToUpdate = new Set<string>([category])
    const toCategory = event.to.dataset.category

    // 如果是跨分类拖拽，需要添加目标分类
    if (event.from !== event.to) {
        console.log('toCategory', toCategory)
        categoriesToUpdate.add(toCategory)
    }
    console.log(categoriesToUpdate)

    // 收集所有需要更新的链接
    let updates: SortIndexUpdate[] = []

    // 处理每个受影响的分类
    for (const cat of categoriesToUpdate) {
        const links = groupedLinksData.value[cat] || []
        const dragLink = links[event.newIndex]
        links.forEach((link, index) => {
            const update: SortIndexUpdate = {
                index: link.globalIndex,
                sortIndex: index + 1
            }

            // 如果是跨分类拖拽，为目标分类添加新的category值
            if (cat === toCategory && link.globalIndex === dragLink.globalIndex) {
                update.category = toCategory
            }

            updates.push(update)
        })
    }

    try {
        // 调用更新排序的 API
        await api.updateSortIndices(updates)

        // 更新本地数据
        store.links = store.links.map(link => {
            const update = updates.find(u => u.index === store.links.indexOf(link))
            if (update) {
                return {
                    ...link,
                    sortIndex: update.sortIndex,
                    // 如果有category更新，也要更新
                    ...(update.category ? { category: update.category } : {})
                }
            }
            return link
        })
    } catch (error) {
        alert('更新排序失败')
        console.error('Failed to update sort indices:', error)
        // 发生错误时恢复原始数据
        updateGroupedLinksData()
    }
}

const onCategoryDragEnd = async (event: any) => {
    if (event.oldIndex === event.newIndex) {
        return
    }

    try {
        // 调用更新分类顺序的 API
        await api.updateCategories(categories.value)
        // 更新本地数据
        store.categories = categories.value
    } catch (error) {
        alert('更新分类顺序失败')
        console.error('Failed to update category order:', error)
        // 发生错误时恢复原始顺序
        updateGroupedLinksData()
    }
}

onMounted(() => {
    fetchLinks()
})
</script>

<style scoped>
.category-drag-handle {
    cursor: move;
    touch-action: none;
}
</style>
