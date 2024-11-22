import { createApp } from "vue";
import App from "./App.vue";
import VueNativeSock from "vue-native-websocket-vue3";

const app = createApp(App);

app.use(VueNativeSock, "ws://localhost:8080/game", {
    connectManually: true,
    reconnection: true,
    reconnectionAttempts: 5,
    reconnectionDelay: 3000,
});

// Монтируем приложение
app.mount("#app");
