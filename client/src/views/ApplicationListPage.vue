<template>
  <div class="application-list">
    <h1>This is an application list page</h1>
    <li><router-link to="/applications/1">id1 </router-link></li>
    <li><router-link to="/applications/2">id2 </router-link></li>
    <li><router-link to="/applications/3">id3 </router-link></li>
    <li><router-link to="/applications/4">id4 </router-link></li>
    <li><router-link to="/applications/5">id5 </router-link></li>
  </div>
</template>
<script>
export default {
  name: "ApplicationList"
};
</script>
