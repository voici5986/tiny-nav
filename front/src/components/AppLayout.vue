<template>
    <div class="min-h-screen flex flex-col">
        <!-- 主要内容 -->
        <main class="flex-grow">
            <slot></slot>
        </main>

        <!-- 页脚 -->
        <Footer />
    </div>
</template>

<script setup lang="ts">
import Footer from './Footer.vue'
</script>
