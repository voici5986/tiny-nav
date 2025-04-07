<template>
    <footer class="bg-white py-4 mt-auto border-t border-gray-200">
        <div class="max-w-5xl mx-auto px-4">
            <div class="text-center text-sm text-gray-500">
                <p>
                    Copyright © {{ currentYear }}
                    <a href="https://github.com/hanxi" target="_blank" rel="noopener noreferrer"
                        class="text-blue-500 hover:text-blue-600 hover:underline">
                        @hanxi
                    </a>
                    <span class="mx-2">·</span>
                    <a href="https://github.com/hanxi/tiny-nav" target="_blank" rel="noopener noreferrer"
                        class="text-blue-500 hover:text-blue-600 hover:underline">
                        Source Code
                    </a>
                </p>
            </div>
        </div>
    </footer>
</template>

<script setup lang="ts">
const currentYear = new Date().getFullYear()
</script>

<style scoped>
a {
    transition: all 0.2s ease-in-out;
}
</style>
