import App from "@/App.vue";
import vuetify from "@/plugins/vuetify";
import router from "@/router";
import "@/styles/index.scss";
import { createPinia } from "pinia";
import { createApp } from "vue";

const app = createApp(App);

app.use(createPinia());
app.use(router);
app.use(vuetify);

app.mount("#app");
