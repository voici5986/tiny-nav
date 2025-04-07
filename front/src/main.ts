import { createApp } from 'vue'
import { createPinia } from 'pinia'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'
import App from './App.vue'
import router from './router'
import 'virtual:uno.css';

// 创建应用实例
const app = createApp(App)

// 创建 Pinia 实例
const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)

// 按顺序使用插件
app.use(pinia)  // Pinia 必须在 router 之前安装
app.use(router)

app.mount('#app')
