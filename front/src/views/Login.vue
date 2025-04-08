<template>
    <AppLayout>
        <div class="bg-white p-8 rounded-lg shadow-lg max-w-md mx-auto mt-20">
            <h1 class="text-4xl font-bold text-center mb-6 text-blue-900">欢迎登录</h1>
            <form @submit.prevent="handleLogin" class="space-y-6">
                <input v-model="form.username" type="text" placeholder="用户名"
                    class="w-full border-2 border-gray-200 focus:border-blue-500 py-3 rounded-lg focus:outline-none transition ease-in-out duration-300" />
                <input v-model="form.password" type="password" placeholder="密码"
                    class="w-full border-2 border-gray-200 focus:border-blue-500 py-3 rounded-lg focus:outline-none transition ease-in-out duration-300" />
                <button type="submit" :disabled="loading"
                    class="w-full bg-blue-500 text-white font-semibold py-3 rounded-lg hover:bg-blue-600 focus:outline-none focus:ring focus:ring-blue-300 disabled:opacity-50 disabled:cursor-not-allowed transform transition-transform hover:scale-105">
                    {{ loading ? '登录中...' : '登录' }}
                </button>
            </form>
        </div>
    </AppLayout>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useMainStore } from '@/stores'
import AppLayout from '@/components/AppLayout.vue'
import { api } from '@/api'

const router = useRouter()
const store = useMainStore()
const loading = ref(false)

const form = reactive({
    username: '',
    password: ''
})

const handleLogin = async () => {
    if (loading.value) return

    loading.value = true
    try {
        const token = await api.login(form)
        store.setToken(token)
        router.push('/nav')
    } catch (error) {
        console.error('Login failed:', error)
        alert('登录失败')
    } finally {
        loading.value = false
    }
}

// 页面加载时自动验证 token
onMounted(async () => {
    console.log(store.token)
    if (store.token) {
        loading.value = true
        try {
            const isValid = await store.validateToken()
            if (isValid) {
                router.push('/nav')
            }
        } finally {
            loading.value = false
        }
    }
})
</script>
