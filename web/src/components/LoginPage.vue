<template>
  <div class="login-page">
    <div class="login-box">
      <h1 class="title">Go-IM</h1>
      <p class="subtitle">A real-time chat application</p>
      <form @submit.prevent="handleLogin" class="login-form">
        <input
          v-model="username"
          type="text"
          placeholder="Enter your name"
          required
        />
        <button type="submit" class="button">Join Chat</button>
        <p v-if="error" class="error-message">{{ error }}</p>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useChatStore } from '@/store';

const username = ref('');
const error = ref('');
const store = useChatStore();

const handleLogin = async () => {
  error.value = '';
  if (!username.value.trim()) {
      error.value = 'Please enter a name.';
      return;
  }
  const success = await store.login(username.value);
  if (!success) {
      error.value = 'Login failed. Please try again.';
  }
};
</script>

<style scoped>
.login-page {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  background-color: var(--bg-1);
}

.login-box {
  width: 100%;
  max-width: 360px;
  padding: 48px;
  text-align: center;
  background-color: var(--bg-2);
  border-radius: 12px;
  box-shadow: 0 8px 24px rgba(0,0,0,0.1);
}

.title {
  font-size: 32px;
  font-weight: 700;
  color: var(--accent-primary);
  margin: 0;
}

.subtitle {
  margin: 8px 0 32px;
  color: var(--text-secondary);
}

.login-form input {
  width: 100%;
  padding: 12px;
  margin-bottom: 16px;
  border-radius: 8px;
  border: 1px solid var(--border-primary);
  font-size: 16px;
  background-color: var(--bg-1);
}
.login-form input:focus {
    outline: none;
    border-color: var(--accent-primary);
    box-shadow: 0 0 0 3px rgba(0, 122, 255, 0.2);
}

.login-form .button {
  width: 100%;
  padding: 12px;
  font-size: 16px;
}

.error-message {
  color: #d93025;
  margin-top: 16px;
  font-size: 14px;
}
</style>