import "@mdi/font/css/materialdesignicons.css";
import { createVuetify } from "vuetify";
import "vuetify/styles";
import { colors } from "./theme";

export default createVuetify({
  theme: {
    defaultTheme: "light",
    themes: {
      light: {
        dark: false,
        colors: {
          background: "#FFFFFF",
          surface: "#FFFFFF",
          ...colors
        }
      },
      dark: {
        dark: true,
        colors: {
          ...colors
        }
      }
    }
  }
});
