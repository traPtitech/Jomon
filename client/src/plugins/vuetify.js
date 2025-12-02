import "@mdi/font/css/materialdesignicons.css";
import { createVuetify } from "vuetify";
import * as components from "vuetify/components";
import * as directives from "vuetify/directives";
// import "vuetify/styles"; // Removed as vite-plugin-vuetify handles this
import light from "./theme";

export default createVuetify({
  components,
  directives,
  theme: {
    defaultTheme: "light",
    themes: {
      light: {
        dark: false,
        colors: {
          background: "#FFFFFF",
          surface: "#FFFFFF",
          primary: "#011A27",
          secondary: "#f0daee",
          error: "#f55a4e",
          info: "#00d3ee",
          success: "#5cb860",
          warning: "#ffa21a",
          ...light
        }
      }
    }
  }
});
