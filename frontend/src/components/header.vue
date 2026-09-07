<script setup lang="ts">
import { ref } from 'vue'
import { Menu, Search, ShoppingCart, User } from 'lucide-vue-next'

const isMenuOpen = ref(false)

function toggleMenu() {
  isMenuOpen.value = !isMenuOpen.value
}
</script>

<template>
  <header class="w-full h-[100px] bg-white flex items-center px-10 relative">
    <div class="flex-1 flex items-center">
      <button
        type="button"
        aria-label="Open menu"
        class="w-[60px] h-[60px] flex items-center justify-center rounded-full hover:bg-neutral-100 transition-colors"
        @click="toggleMenu"
      >
        <Menu :size="22" :stroke-width="1.8" color="#111111" />
      </button>
    </div>

    <a href="/" class="font-sans font-extrabold text-3xl tracking-tight text-black shrink-0">
      FLOM
    </a>

    <div class="flex-1 flex items-center justify-end gap-2">
      <button
        type="button"
        aria-label="Search"
        class="w-[50px] h-[50px] flex items-center justify-center rounded-full hover:bg-neutral-100 transition-colors"
      >
        <Search :size="20" :stroke-width="1.8" color="#111111" />
      </button>

      <button
        type="button"
        aria-label="Cart"
        class="w-[50px] h-[50px] flex items-center justify-center rounded-full hover:bg-neutral-100 transition-colors"
      >
        <ShoppingCart :size="20" :stroke-width="1.8" color="#111111" />
      </button>

      <button
        type="button"
        aria-label="Account"
        class="w-[50px] h-[50px] flex items-center justify-center rounded-full hover:bg-neutral-100 transition-colors"
      >
        <User :size="20" :stroke-width="1.8" color="#111111" />
      </button>
    </div>

    <!--
      Placeholder mobile/dropdown menu panel.
      Wire up real nav links once vue-router is installed.
    -->
    <div
      v-if="isMenuOpen"
      class="absolute left-0 top-[100px] w-full bg-white border-t border-neutral-200 shadow-lg z-50 p-8"
    >
      <nav class="flex flex-col gap-4 font-sans text-lg text-black">
        <a href="/" class="hover:text-[#FF591F] transition-colors">Home</a>
        <a href="/shop" class="hover:text-[#FF591F] transition-colors">Shop</a>
        <a href="/about" class="hover:text-[#FF591F] transition-colors">About</a>
        <a href="/contact" class="hover:text-[#FF591F] transition-colors">Contact</a>
      </nav>
    </div>
  </header>
</template>