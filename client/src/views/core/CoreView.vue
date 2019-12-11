<!-- メイン描画画面制御 -->
<template>
  <v-content>
    <nav>
      <router-link to="/">applications </router-link>|
      <router-link to="/applications/:id">applications_id </router-link>|
      <router-link to="/admin">admin </router-link>|
      <router-link to="/applications/new">new_application </router-link>|
    </nav>
    <router-view />
  </v-content>
</template>
