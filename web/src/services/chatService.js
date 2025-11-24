
import axios from 'axios';
import { v4 as uuidv4 } from 'uuid';

const API_BASE_URL = '/v1/api'; // Use proxy
const WS_BASE_URL = 'ws://localhost:8070/v1/api/ws'; // Connect directly, bypassing Vite proxy for diagnostics

let socket = null;

const DEVICE_ID_KEY = 'chatDeviceId';

const ChatService = {
  // --- Device ID Management ---
  getDeviceId() {
    let deviceId = localStorage.getItem(DEVICE_ID_KEY);
    if (!deviceId) {
      deviceId = uuidv4();
      localStorage.setItem(DEVICE_ID_KEY, deviceId);
    }
    return deviceId;
  },
  
  // --- API Calls ---
  async login(username) {
    const response = await axios.post(`${API_BASE_URL}/login`, { username });
    return response.data.data;
  },
  
  async getUsers() {
    const response = await axios.get(`${API_BASE_URL}/users`);
    return response.data.data;
  },

  async getRooms() {
    const response = await axios.get(`${API_BASE_URL}/rooms`);
    return response.data.data;
  },

  async createRoom({ name, password, userId }) {
    const response = await axios.post(`${API_BASE_URL}/rooms`, { name, password, userId });
    return response.data.data;
  },

  async joinRoom({ roomId, userId, deviceId, password }) {
    const response = await axios.post(`${API_BASE_URL}/rooms/${roomId}/join`, { userId, deviceId: deviceId, roomId, password });
    return response.data;
  },

  async getRoomDetail(roomId) {
    const response = await axios.get(`${API_BASE_URL}/rooms/${roomId}`);
    return response.data.data;
  },

  // --- WebSocket Management ---
  connect(userId, deviceId, username, { onOpen, onMessage, onClose, onError }) {
    if (socket && socket.readyState === WebSocket.OPEN) {
      console.log('WebSocket is already connected.');
      return;
    }

    const url = `${WS_BASE_URL}?token=${userId}&deviceId=${deviceId}&username=${encodeURIComponent(username)}`;
    socket = new WebSocket(url);

    socket.onopen = onOpen;
    socket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        onMessage?.(data);
      } catch (e) {
        console.error('Error parsing WebSocket message:', e);
      }
    };
    socket.onclose = () => {
      onClose?.();
      socket = null;
    };
    socket.onerror = onError;
  },

  sendMessage(payload) {
    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.send(JSON.stringify(payload));
    } else {
      console.error('WebSocket is not connected.');
    }
  },

  disconnect() {
    if (socket) {
      socket.close();
    }
  },
};

export default ChatService;
