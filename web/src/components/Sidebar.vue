<template>
  <aside class="sidebar">
    <!-- User Profile -->
    <div class="user-profile">
      <div class="avatar">
        {{ store.user?.name?.charAt(0).toUpperCase() }}
      </div>
      <div class="user-info">
        <span class="user-name">{{ store.user?.name }}</span>
        <button @click="handleLogout" class="logout-button">Logout</button>
      </div>
    </div>

    <!-- Actions -->
    <div class="actions">
       <button class="button" @click="isCreateRoomModalVisible = true">
         <svg width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M12 4V20M20 12L4 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>
         <span>Create Room</span>
       </button>
    </div>

    <!-- Room List -->
    <div class="chat-list">
      <h3 class="list-header">Rooms ({{ store.rooms.length }})</h3>
      <ul>
        <li
          v-for="room in store.rooms"
          :key="room.id"
          class="list-item"
          :class="{ active: store.activeChatTarget?.type === 'room' && store.activeChatTarget?.id === room.id }"
          @click="selectRoom(room)"
        >
          <span class="room-name"># {{ room.name }}</span>
          <span class="user-count">{{ room.count }}</span>
        </li>
      </ul>
    </div>

    <!-- User List -->
    <div class="chat-list user-list">
      <h3 class="list-header">Direct Messages ({{ store.onlineUsers.length }})</h3>
      <ul>
        <li
          v-for="user in store.onlineUsers"
          :key="user.id"
          class="list-item"
          :class="{ 
            active: store.activeChatTarget?.type === 'user' && store.activeChatTarget?.id === user.id,
            'has-unread': store.unreadFromUsers.has(user.id)
          }"
          @click="selectUser(user)"
        >
          <span class="user-status online"></span>
          <span class="user-name">{{ user.name }}</span>
          <span v-if="store.unreadFromUsers.has(user.id)" class="unread-badge"></span>
        </li>
      </ul>
    </div>

    <!-- Create Room Modal -->
    <div v-if="isCreateRoomModalVisible" class="modal-overlay" @click.self="isCreateRoomModalVisible = false">
      <div class="modal-content">
        <h3>Create a New Room</h3>
        <form @submit.prevent="handleCreateRoom" class="modal-form">
          <input v-model="newRoom.name" placeholder="Room name" required />
          <input v-model="newRoom.password" placeholder="Password (optional)" type="password" />
          <div class="modal-actions">
            <button type="button" class="button secondary" @click="isCreateRoomModalVisible = false">Cancel</button>
            <button type="submit" class="button">Create</button>
          </div>
        </form>
      </div>
    </div>
  </aside>
</template>

<script setup>
import { ref } from 'vue';
import { useChatStore } from '@/store';

const store = useChatStore();
const isCreateRoomModalVisible = ref(false);
const newRoom = ref({ name: '', password: '' });

const handleLogout = () => store.logout();

const selectRoom = (room) => {
  store.selectChatTarget({ type: 'room', id: room.id, name: room.name });
};

const selectUser = (user) => {
  store.selectChatTarget({ type: 'user', id: user.id, name: user.name });
};

const handleCreateRoom = async () => {
  if (!newRoom.value.name.trim()) return;
  await store.createNewRoom(newRoom.value);
  newRoom.value = { name: '', password: '' };
  isCreateRoomModalVisible.value = false;
};
</script>

<style scoped>
.sidebar {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: var(--bg-3);
  border-right: 1px solid var(--border-primary);
}

.user-profile {
  display: flex;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid var(--border-primary);
  flex-shrink: 0;
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background-color: var(--accent-primary);
  color: var(--text-interactive);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  font-weight: 600;
  margin-right: 12px;
}

.user-info {
  flex-grow: 1;
  display: flex;
  flex-direction: column;
}

.user-name {
  font-weight: 600;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.logout-button {
  background: none;
  border: none;
  color: var(--text-tertiary);
  cursor: pointer;
  padding: 0;
  margin-top: 2px;
  font-size: 12px;
  text-align: left;
}
.logout-button:hover {
  color: var(--accent-primary);
}

.actions {
  padding: 16px;
  border-bottom: 1px solid var(--border-primary);
}
.actions .button {
  width: 100%;
  gap: 8px;
}

.chat-list {
  padding: 16px 8px;
  overflow-y: auto;
  flex-grow: 1;
}
.user-list {
  border-top: 1px solid var(--border-primary);
}

.list-header {
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  color: var(--text-secondary);
  margin: 0 8px 8px;
}

.list-item {
  display: flex;
  align-items: center;
  padding: 8px;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 500;
  color: var(--text-secondary);
  margin-bottom: 4px;
}
.list-item:hover {
  background-color: var(--hover-bg);
}
.list-item.active {
  background-color: var(--accent-primary);
  color: var(--text-interactive);
}
.list-item .user-name, .list-item .room-name {
  flex-grow: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-status {
  width: 9px;
  height: 9px;
  border-radius: 50%;
  margin-right: 10px;
  flex-shrink: 0;
}
.user-status.online {
  background-color: #28a745;
}

.user-count {
  font-size: 12px;
  padding: 2px 6px;
  border-radius: 8px;
  background-color: rgba(0, 0, 0, 0.08);
}
.list-item.active .user-count {
  background-color: rgba(255, 255, 255, 0.2);
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}
.modal-content {
  background: var(--bg-2);
  padding: 24px;
  border-radius: 12px;
  width: 90%;
  max-width: 400px;
  box-shadow: 0 5px 15px rgba(0,0,0,0.2);
}
.modal-content h3 {
  margin-top: 0;
  margin-bottom: 24px;
}
.modal-form input {
  width: 100%;
  padding: 10px;
  margin-bottom: 16px;
  border-radius: 6px;
  border: 1px solid var(--border-primary);
  font-size: 16px;
  background-color: var(--bg-1);
}
.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 8px;
}

/* Unread Notification Styles */
.list-item.has-unread .user-name {
  font-weight: 700;
  color: var(--text-primary);
}
.list-item.active.has-unread .user-name {
    color: var(--text-interactive);
}

.unread-badge {
  width: 9px;
  height: 9px;
  border-radius: 50%;
  background-color: var(--accent-unread);
  margin-left: auto;
  flex-shrink: 0;
  animation: pulse 1.5s infinite;
}

@keyframes pulse {
  0% {
    transform: scale(0.95);
    box-shadow: 0 0 0 0 rgba(0, 122, 255, 0.7);
  }
  70% {
    transform: scale(1);
    box-shadow: 0 0 0 5px rgba(0, 122, 255, 0);
  }
  100% {
    transform: scale(0.95);
    box-shadow: 0 0 0 0 rgba(0, 122, 255, 0);
  }
}
</style>
