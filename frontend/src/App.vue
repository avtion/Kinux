<template>
  <router-view></router-view>
</template>

<script lang="ts">
import { defineComponent } from "vue"

export default defineComponent({
  name: "App",
})
</script>

<style>
@import "./assets/css/global.css";

#app {
  height: 100%;
  font-family: "Helvetica Neue", Helvetica, "PingFang SC", "Hiragino Sans GB",
    "Microsoft YaHei", "微软雅黑", Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
</style>
