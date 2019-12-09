<!-- メイン描画画面制御 -->
<template>
  <v-content>
    <nav>
      <router-link to="/">applications </router-link>|
      <router-link to="/applications/:id">applicationsid </router-link>|
      <router-link to="/admin">admin </router-link>|
      <router-link to="/applications/new">applicationnew </router-link>|
    </nav>
    <router-view />
  </v-content>
</template>
