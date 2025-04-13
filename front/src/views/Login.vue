<template>
    <AppLayout>
        <div class="bg-white dark:bg-gray-800 p-8 rounded-lg shadow-lg dark:shadow-gray-900 max-w-md mx-auto mt-20">
            <!-- 标题和深色模式按钮 -->
            <div class="flex justify-between items-center mb-6">
                <h1 class="text-4xl font-bold text-blue-900 dark:text-blue-300">欢迎登录</h1>
                <button @click="themeStore.toggleTheme"
                    class="w-10 h-10 flex items-center justify-center rounded-full bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-100 hover:bg-gray-300 dark:hover:bg-gray-600">
                    <div :class="themeStore.isDarkTheme ? 'i-mdi-white-balance-sunny' : 'i-mdi-moon-waxing-crescent'">
                    </div>
                </button>
            </div>

            <form @submit.prevent="handleLogin" class="space-y-6">
                <input v-model="form.username" type="text"
                    class="w-full border-2 border-gray-200 dark:border-gray-600 dark:bg-gray-700 dark:text-gray-200 focus:border-blue-500 dark:focus:border-blue-400 py-3 rounded-lg focus:outline-none transition ease-in-out duration-300"
                    placeholder="用户名" />
                <input v-model="form.password" type="password"
                    class="w-full border-2 border-gray-200 dark:border-gray-600 dark:bg-gray-700 dark:text-gray-200 focus:border-blue-500 dark:focus:border-blue-400 py-3 rounded-lg focus:outline-none transition ease-in-out duration-300"
                    placeholder="密码" />
                <button type="submit" :disabled="loading"
                    class="w-full bg-blue-500 dark:bg-blue-400 text-white dark:text-gray-900 font-semibold py-3 rounded-lg hover:bg-blue-600 dark:hover:bg-blue-500 focus:outline-none focus:ring focus:ring-blue-300">
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
import { useThemeStore } from '@/stores/themeStore'
import AppLayout from '@/components/AppLayout.vue'
import { api } from '@/api'

const themeStore = useThemeStore()
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

// 自动检测无用户密码模式
const detectNoAuthMode = async () => {
    try {
        const token = await api.login({ username: '', password: '' }) // 尝试使用空账号密码登录
        store.setToken(token) // 设置 Token
        store.setNoAuthMode(true) // 设置无用户密码模式为 true
        router.push('/nav') // 跳转到导航页面
    } catch (error) {
        console.log('No-auth mode not enabled or login failed:', error)
        store.setNoAuthMode(false) // 设置无用户密码模式为 false
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
    await detectNoAuthMode() // 尝试无用户密码登录
})
</script>
