<template>
  <div id="app">
    <el-menu
      id="header"
      default-active="1"
      mode="horizontal"
      @select="handleSelect"
      background-color="#545c64"
      text-color="#fff"
      active-text-color="#ffd04b">
      <el-menu-item id="brand" disabled>dou-ping</el-menu-item>
      <el-menu-item index="1" class="menu-item">PING 主机</el-menu-item>
      <el-menu-item index="2" class="menu-item">PING 网络</el-menu-item>
    </el-menu>

    <ip-ping v-if="curMenuItem == 1"></ip-ping>

    <ipnet-ping v-else></ipnet-ping>
  </div>
</template>

<script>
import "./assets/css/main.css";
import ipPing from "./components/ip-ping"
import ipnetPing from "./components/ipnet-ping"

export default {
  name: "app",
  components: {
    ipPing,
    ipnetPing
  },
  data() {
    return {
      curMenuItem: 1
    }
  },
  methods: {
    handleSelect(index) {
      window.console.log(index, typeof index)
      this.curMenuItem = index
    }
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
