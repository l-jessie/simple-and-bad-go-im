import {defineStore} from 'pinia';
import ChatService from '@/services/chatService';

const USER_SESSION_KEY = 'chatUserSession';
const MAX_MESSAGES_ROOM = 100;
const MAX_MESSAGES_PRIVATE = 500;

// Helper to generate a consistent, sorted ID for private chats
const generatePrivateChatId = (userId1, userId2) => {
    return [userId1, userId2].sort().join('-');
};

export const useChatStore = defineStore('chat', {
    state: () => ({
        user: null, // { id, name }
        activeChatTarget: null, // { type: 'room' | 'user', id: '...', name: '...' }
        rooms: [],
        onlineUsers: [],
        usersMap: {}, // Map of userId to {id, name} for quick lookup
        messagesByChat: {}, // e.g., { 'chatId': [messages] }
        chatHistoryTruncated: {}, // e.g., { 'chatId': true }
        isConnected: false,
        unreadFromUsers: new Set(),
        currentRoomDetail: null,
    }),

    getters: {
        activeChatId() {
            if (!this.activeChatTarget || !this.user) return null;
            if (this.activeChatTarget.type === 'room') {
                return this.activeChatTarget.id;
            }
            return generatePrivateChatId(this.user.id, this.activeChatTarget.id);
        },
        activeChatMessages() {
            return this.messagesByChat[this.activeChatId] || [];
        },
        isCurrentChatTruncated() {
            return this.chatHistoryTruncated[this.activeChatId] || false;
        },
        currentChatLimit() {
            if (!this.activeChatTarget) return 0;
            return this.activeChatTarget.type === 'room' ? MAX_MESSAGES_ROOM : MAX_MESSAGES_PRIVATE;
        }
    },

    actions: {
        // --- Message Storage Helper ---
        _addMessage(chatId, message, limit) {
            if (!this.messagesByChat[chatId]) {
                this.messagesByChat[chatId] = [];
            }

            const messageList = this.messagesByChat[chatId];
            messageList.push(message);

            if (messageList.length > limit) {
                messageList.shift(); // Remove the oldest message
                this.chatHistoryTruncated[chatId] = true;
            }
        },
        _addUserToMap(user) {
            if (user && user.id && user.name) {
                this.usersMap[user.id] = {id: user.id, name: user.name};
            }
        },

        // --- Authentication and Connection ---
        async login(username) {
            try {
                const userData = await ChatService.login(username);
                this.user = {id: userData.id, name: userData.username};
                this._addUserToMap(this.user); // Add self to map
                sessionStorage.setItem(USER_SESSION_KEY, JSON.stringify(this.user));
                this.connectWebSocket();
                return true;
            } catch (error) {
                console.error('Login failed:', error);
                this.logout();
                return false;
            }
        },

        logout() {
            ChatService.disconnect();
            this.$reset(); // Use pinia's $reset to go back to initial state
            sessionStorage.removeItem(USER_SESSION_KEY);
        },

        connectWebSocket() {
            if (!this.user?.id || !this.user.name) return;

            ChatService.connect(this.user.id, ChatService.getDeviceId(), this.user.name, {
                onOpen: () => {
                    this.isConnected = true;
                    this.fetchRooms();
                    this.fetchUsers();
                },
                onMessage: (data) => {
                    console.log('Received message:', data);

                    if (data.from === this.user?.id && data.type !== 3) {
                        // Ignore echoes from room chats. We keep private chat echoes
                        // in case the user is chatting with themselves on another device.
                        // But we already optimistically display, so for now, ignore all.
                        return;
                    }
                    if (data.type === 4) { // MessageTypeGlobal
                        if (data.messageEvent) {
                            if (data.messageEvent.type === 0) { // ReloadUsers
                                this.fetchUsers();
                                console.log('Received ReloadUsers event');
                            } else if (data.messageEvent.type === 1) { // ReloadRoomsDetail
                                // Assuming data.messageEvent.data contains the roomId as a string
                                const roomIdToReload = JSON.parse(data.messageEvent.data);
                                if (this.activeChatTarget?.type === 'room' && this.activeChatTarget.id === roomIdToReload) {
                                    this.fetchCurrentRoomDetail(roomIdToReload);
                                    console.log(`Received ReloadRoomsDetail for active room: ${roomIdToReload}`);
                                }
                            } else if (data.messageEvent.type === 2) { // ReloadRooms
                                this.fetchRooms();
                                console.log('Received ReloadRooms event');
                            }
                        }
                        return;
                    }

                    let chatId;
                    let limit;

                    if (data.type === 2) { // Room Message
                        chatId = data.to;
                        limit = MAX_MESSAGES_ROOM;
                    } else if (data.type === 3) { // Private Message
                        const senderId = data.from;
                        chatId = generatePrivateChatId(this.user.id, senderId);
                        limit = MAX_MESSAGES_PRIVATE;

                        if (this.activeChatId !== chatId) {
                            this.unreadFromUsers.add(senderId);
                        }
                    } else {
                        return
                    }

                    const senderName = this.usersMap[data.from]?.name || data.from; // Fallback to ID if name not found
                    const formattedMessage = {
                        id: data.timestamp,
                        content: data.payload?.data || '',
                        user: {id: data.from, name: senderName},
                    };
                    this._addMessage(chatId, formattedMessage, limit);
                },
                onClose: () => {
                    this.isConnected = false;
                },
                onError: (error) => {
                    console.error('WebSocket error:', error);
                },
            });
        },

        checkExistingSession() {
            const savedUser = sessionStorage.getItem(USER_SESSION_KEY);
            if (savedUser) {
                this.user = JSON.parse(savedUser);
                this._addUserToMap(this.user); // Add self to map on session restore
                this.connectWebSocket();
            }
        },

        // --- User Actions ---
        selectChatTarget({type, id, name}) {
            this.activeChatTarget = {type, id, name};
            this.currentRoomDetail = null;

            if (type === 'room') {
                const room = this.rooms.find(r => r.id === id);
                if (room?.hasPassword) {
                    const password = prompt('This room requires a password:', '');
                    if (password !== null) {
                        this.joinRoom(id, password).then(() => this.fetchCurrentRoomDetail(id));
                    } else {
                        this.activeChatTarget = null;
                    }
                } else {
                    this.joinRoom(id).then(() => this.fetchCurrentRoomDetail(id));
                }
            } else if (type === 'user') {
                this.unreadFromUsers.delete(id);
            }
        },

        sendMessage(text) {
            if (!this.activeChatId) return;

            const messagePayload = {
                type: this.activeChatTarget.type === 'room' ? 2 : 3,
                to: this.activeChatTarget.id,
                payload: {type: 0, data: text}
            };
            ChatService.sendMessage(messagePayload);

            const optimisticMessage = {
                id: Date.now(),
                content: text,
                user: {id: this.user.id, name: this.user.name},
            };
            const limit = this.currentChatLimit;
            this._addMessage(this.activeChatId, optimisticMessage, limit);
        },

        // --- API-driven Actions ---
        async fetchRooms() {
            try {
                const rooms = await ChatService.getRooms();
                this.rooms = rooms || [];
                // Add room owner/members to usersMap if available
                // This assumes RoomResponse has userName field for owner
                this.rooms.forEach(room => {
                    if (room.userId && room.userName) {
                        this._addUserToMap({id: room.userId, name: room.userName});
                    }
                });
            } catch (error) {
                console.error('Failed to fetch rooms:', error);
                this.rooms = [];
            }
        },
        async fetchUsers() {
            try {
                const users = await ChatService.getUsers();
                this.onlineUsers = users ? users.filter(u => u.id !== this.user?.id) : [];
                this.onlineUsers.forEach(user => this._addUserToMap(user)); // Add all online users to map
            } catch (error) {
                console.error('Failed to fetch users:', error);
                this.onlineUsers = [];
            }
        },
        async createNewRoom({name, password}) {
            if (!this.user?.id) return;
            try {
                const newRoom = await ChatService.createRoom({name, password, userId: this.user.id});
                this.fetchRooms(); // Refetch all rooms
                // Optionally, if newRoom contains owner details, add to map
                if (newRoom.userId && newRoom.userName) {
                    this._addUserToMap({id: newRoom.userId, name: newRoom.userName});
                }
            } catch (error) {
                console.error('Failed to create new room:', error);
                alert(`Error creating room: ${error.response?.data?.msg || error.message}`);
            }
        },
        async fetchCurrentRoomDetail(roomId) {
            try {
                const roomDetail = await ChatService.getRoomDetail(roomId);
                this.currentRoomDetail = roomDetail;
                // Add room members to usersMap
                roomDetail.users.forEach(user => this._addUserToMap(user));
                // Add room owner to usersMap
                if (roomDetail.userId && roomDetail.userName) {
                    this._addUserToMap({id: roomDetail.userId, name: roomDetail.userName});
                }
            } catch (error) {
                console.error(`Failed to fetch room detail for room ${roomId}:`, error);
                this.currentRoomDetail = null;
                if (this.activeChatTarget?.id === roomId) {
                    this.activeChatTarget = null;
                }
                alert(`Error fetching room details: ${error.response?.data?.msg || error.message}`);
            }
        },
        async joinRoom(roomId, password = '') {
            if (!this.user?.id) return;
            try {
                const deviceId = ChatService.getDeviceId();
                await ChatService.joinRoom({roomId, userId: this.user.id, deviceId, password});
                console.log(`Successfully joined room ${roomId}`);
            } catch (error) {
                console.error(`Failed to join room ${roomId}:`, error);
                if (this.activeChatTarget?.id === roomId) {
                    this.activeChatTarget = null;
                }
                alert(`Error joining room: ${error.response?.data?.msg || error.message}`);
            }
        },
    },
});
