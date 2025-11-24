<template>
  <div class="message-display" ref="messagesContainer">
    <div v-if="store.isCurrentChatTruncated" class="history-limit-notice">
      Message history is limited to the last {{ store.currentChatLimit }} messages.
    </div>
    <div 
      v-for="msg in store.activeChatMessages" 
      :key="msg.id" 
      class="message-group" 
      :class="isOwnMessage(msg) ? 'own' : 'other'"
    >
      <div class="avatar">
        {{ msg.user.name?.charAt(0).toUpperCase() }}
      </div>
      <div class="message-content">
        <div class="message-sender" v-if="!isOwnMessage(msg)">{{ msg.user.name }}</div>
        <div class="message-bubble">
          {{ msg.content }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, nextTick } from 'vue';
import { useChatStore } from '@/store';

const store = useChatStore();
const messagesContainer = ref(null);

const isOwnMessage = (msg) => {
    return msg.user.id === store.user?.id;
}

const scrollToBottom = () => {
  nextTick(() => {
    const container = messagesContainer.value;
    if (container) {
      container.scrollTop = container.scrollHeight;
    }
  });
};

watch(() => store.activeChatMessages, scrollToBottom, { deep: true });
watch(() => store.activeChatId, scrollToBottom, { immediate: true });

</script>

<style scoped>
.message-display {
  flex-grow: 1;
  overflow-y: auto;
  padding: 24px;
}

.history-limit-notice {
  text-align: center;
  font-size: 12px;
  color: var(--text-tertiary);
  margin-bottom: 20px;
  padding: 4px;
  background-color: var(--bg-3);
  border-radius: 6px;
}

.message-group {
  display: flex;
  margin-bottom: 20px;
  gap: 12px;
}

.avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: 600;
  flex-shrink: 0;
}

.message-content {
  display: flex;
  flex-direction: column;
  max-width: 75%;
}

.message-sender {
  font-size: 13px;
  font-weight: 600;
  margin-bottom: 6px;
  color: var(--text-secondary);
}

.message-bubble {
  padding: 10px 14px;
  border-radius: 18px;
  line-height: 1.5;
  word-wrap: break-word;
}

/* Other's messages */
.message-group.other .avatar {
  background-color: var(--other-message-bg);
  color: var(--other-message-text);
}
.message-group.other .message-bubble {
  background-color: var(--other-message-bg);
  color: var(--other-message-text);
  border-top-left-radius: 4px;
}

/* Own messages */
.message-group.own {
  flex-direction: row-reverse;
}
.message-group.own .avatar {
  background-color: var(--own-message-bg);
  color: var(--own-message-text);
}
.message-group.own .message-content {
  align-items: flex-end;
}
.message-group.own .message-bubble {
  background-color: var(--own-message-bg);
  color: var(--own-message-text);
  border-top-right-radius: 4px;
}
</style>
