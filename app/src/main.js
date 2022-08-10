import { createApp } from "vue";
import App from "./App.vue";
import axios from "axios";
import router from "@/router";

axios.defaults.baseURL = "http://localhost:8082";

createApp(App).use(router).mount('#app')
