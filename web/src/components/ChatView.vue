<template>
  <div class="chat-view-container">
    <!-- Placeholder for when no chat is selected -->
    <div v-if="!store.activeChatTarget" class="no-chat-selected">
      <div class="content">
        <h2>Welcome to Go-IM</h2>
        <p>Select a room or a user from the sidebar to start chatting.</p>
      </div>
    </div>

    <!-- Main Chat Interface -->
    <div v-else class="chat-view" :class="{ 'room-mode': isRoomChat }">
      <header class="chat-header">
        <div class="chat-title">
          <span v-if="isRoomChat"># {{ store.activeChatTarget.name }}</span>
          <span v-else>{{ store.activeChatTarget.name }}</span>
        </div>
        <div class="chat-meta">
          <span v-if="isRoomChat && roomDetails">
            {{ roomDetails.count }} members
          </span>
        </div>
      </header>
      
      <MessageDisplay />
      <MessageInput />

      <!-- Right sidebar for in-room users -->
      <aside v-if="isRoomChat" class="room-members-sidebar">
        <div class="sidebar-header">
          <h3>Members</h3>
        </div>
        <div class="member-list">
            <div v-if="!roomDetails" class="loading-members">Loading...</div>
            <ul v-else>
              <li v-for="user in roomDetails.users" :key="user.id" class="member-item">
                <div class="avatar">
                  {{ user.name?.charAt(0).toUpperCase() }}
                </div>
                <span>{{ user.name }}</span>
              </li>
            </ul>
        </div>
      </aside>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { useChatStore } from '@/store';
import MessageDisplay from './MessageDisplay.vue';
import MessageInput from './MessageInput.vue';

const store = useChatStore();

const isRoomChat = computed(() => store.activeChatTarget?.type === 'room');
const roomDetails = computed(() => store.currentRoomDetail);

</script>

<style scoped>
.chat-view-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
}

/* --- No Chat Selected View --- */
.no-chat-selected {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  text-align: center;
  color: var(--text-secondary);
}
.no-chat-selected .content h2 {
  font-size: 24px;
  font-weight: 600;
  color: var(--text-primary);
}
.no-chat-selected .content p {
  font-size: 16px;
  margin-top: 8px;
}


/* --- Main Chat View --- */
.chat-view {
  display: grid;
  height: 100%;
  width: 100%;
  grid-template-columns: 1fr;
  grid-template-rows: 60px 1fr 85px;
  grid-template-areas:
    "header"
    "main"
    "footer";
}

.chat-view.room-mode {
  grid-template-columns: 1fr 240px;
  grid-template-areas:
    "header header"
    "main members"
    "footer members";
}

.chat-header {
  grid-area: header;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  border-bottom: 1px solid var(--border-primary);
  flex-shrink: 0;
}
.chat-title {
  font-size: 18px;
  font-weight: 600;
}
.chat-meta {
  font-size: 14px;
  color: var(--text-secondary);
}

.room-members-sidebar {
  grid-area: members;
  display: flex;
  flex-direction: column;
  border-left: 1px solid var(--border-primary);
  background-color: var(--bg-3);
  height: 100%;
  overflow: hidden;
}
.sidebar-header {
  padding: 16px;
  border-bottom: 1px solid var(--border-primary);
  flex-shrink: 0;
}
.sidebar-header h3 {
  margin: 0;
  font-size: 16px;
}
.member-list {
  padding: 16px;
  overflow-y: auto;
  flex-grow: 1;
}
.member-list ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.member-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 0;
  font-weight: 500;
}
.member-item .avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background-color: var(--other-message-bg);
  color: var(--text-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  font-weight: 600;
  flex-shrink: 0;
}
</style>
