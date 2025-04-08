<template>
    <footer class="container mx-auto">
        <div class="flex items-center justify-center text-gray-500">
            <!-- 左侧内容：版权信息和@hanxi链接 -->
            <div class="flex items-center gap-2">
                <span>Copyright © {{ currentYear }}</span>
                <a href="https://github.com/hanxi" target="_blank"
                    class="text-blue-600 hover:underline transition-colors">
                    @hanxi
                </a>
            </div>

            <!-- 分隔符点 -->
            <span class="mx-2 text-gray-500">·</span>

            <!-- 右侧内容：GitHub图标链接 -->
            <div class="flex items-center gap-2">
                <a href="https://github.com/hanxi/tiny-nav" target="_blank"
                    class="i-mdi-github text-gray-600 hover:text-blue-600 transition-colors"></a>
            </div>
        </div>
    </footer>
</template>

<script setup lang="ts">
const currentYear = new Date().getFullYear()
</script>
