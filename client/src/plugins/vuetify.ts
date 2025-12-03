import "@mdi/font/css/materialdesignicons.css";
import { createVuetify } from "vuetify";
// import "vuetify/styles"; // Removed as vite-plugin-vuetify handles this
import light from "./theme";

export default createVuetify({
  theme: {
    defaultTheme: "light",
    themes: {
      light: {
        dark: false,
        colors: {
          background: "#FFFFFF",
          surface: "#FFFFFF",
          ...light
        }
      }
    }
  }
});
