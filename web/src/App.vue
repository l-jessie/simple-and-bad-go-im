<template>
  <div id="app-container">
    <template v-if="!store.user">
      <LoginPage />
    </template>
    <main v-else class="app-layout">
      <Sidebar />
      <ChatView />
    </main>
  </div>
</template>

<script setup>
import { onMounted } from 'vue';
import { useChatStore } from '@/store';
import LoginPage from '@/components/LoginPage.vue';
import Sidebar from '@/components/Sidebar.vue';
import ChatView from '@/components/ChatView.vue';

const store = useChatStore();

onMounted(() => {
  store.checkExistingSession();
});
</script>

<style scoped>
.app-layout {
  display: grid;
  grid-template-columns: 280px 1fr;
  height: 100vh;
  width: 100vw;
  background-color: var(--bg-2);
}
</style>