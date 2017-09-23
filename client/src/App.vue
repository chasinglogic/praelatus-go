<template>
  <div id="app">
    <navbar></navbar>
    <div class="content" v-bind:style="{ marginLeft: sidebarWidth }">
      <router-view></router-view>
    </div>
  </div>
</template>

<script>
import Navbar from '@/components/General/Navbar'

export default {
  name: 'app',
  computed: {
    sidebarWidth: function () {
      if (this.$store.getters.showSidebar) {
        return this.$store.getters.sidebarWidth
      }

      return '0px'
    }
  },
  components: {
    Navbar
  }
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
