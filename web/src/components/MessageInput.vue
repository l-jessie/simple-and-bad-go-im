<template>
  <form @submit.prevent="sendMessage" class="message-input-area">
    <div class="input-container">
      <input
        v-model="newMessage"
        type="text"
        :placeholder="inputPlaceholder"
        autocomplete="off"
        :disabled="!store.activeChatTarget"
        @keyup.enter.exact="sendMessage"
      />
      <button type="submit" :disabled="!store.activeChatTarget || !isMessageValid">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M22 2L11 13" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          <path d="M22 2L15 22L11 13L2 9L22 2Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
      </button>
    </div>
  </form>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useChatStore } from '@/store';

const newMessage = ref('');
const store = useChatStore();

const isMessageValid = computed(() => newMessage.value.trim() !== '');

const inputPlaceholder = computed(() => {
  if (!store.activeChatTarget) {
    return 'Select a chat to begin';
  }
  return `Message ${store.activeChatTarget.type === 'room' ? '#' : ''}${store.activeChatTarget.name}`;
});

const sendMessage = () => {
  if (isMessageValid.value && store.activeChatTarget) {
    store.sendMessage(newMessage.value.trim());
    newMessage.value = '';
  }
};
</script>

<style scoped>
.message-input-area {
  padding: 16px 24px;
  background-color: var(--bg-2);
  border-top: 1px solid var(--border-primary);
}

.input-container {
  display: flex;
  align-items: center;
  background-color: var(--other-message-bg);
  border-radius: 12px;
  padding: 4px;
}

input {
  flex-grow: 1;
  border: none;
  background: transparent;
  padding: 10px 12px;
  font-size: 15px;
  outline: none;
  color: var(--text-primary);
}

button {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: none;
  border-radius: 8px;
  background-color: var(--accent-primary);
  color: var(--text-interactive);
  cursor: pointer;
  transition: background-color 0.2s;
}

button:hover {
  background-color: var(--accent-primary-hover);
}

button:disabled {
  background-color: #a8c5ff;
  cursor: not-allowed;
}
</style>
