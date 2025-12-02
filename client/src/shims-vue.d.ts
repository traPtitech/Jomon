/* eslint-disable */
import { RouteLocationNormalizedLoaded, Router } from "vue-router";

declare module "@vue/runtime-core" {
  interface ComponentCustomProperties {
    $store: any;
    $router: Router;
    $route: RouteLocationNormalizedLoaded;
    $vuetify: any;
  }
}

export {};
