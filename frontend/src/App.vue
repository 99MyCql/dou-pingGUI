<template>
  <div id="app">
  </div>
</template>

<script>
import "./assets/css/main.css";
import ipPing from "./components/ip-ping"

export default {
  name: "app",
  components: {
    ipPing
  }
};
</script>

<style scoped>
#app {
  position: absolute;
  width: 100%;
  min-width: 1024px;
  height: 100%;
  background-color: white;
}

#header {
  padding-left: 10%;
  box-shadow: 0 0 10px 0 rgba(0,0,0,0.2);
}

#brand {
  color: white;
  font-size: 25px;
  margin-right: 100px;
}

.menu-item {
  padding: 0 60px;
}
</style>
