<template>
  <div class="sidebar" v-bind:style="{ width: width }" v-if="shown">
    <h1>Sidebar</h1>
    <slot></slot>
  </div>
</template>

<script>
  export default {
    name: 'sidebar',
    props: {
      shown: { type: Boolean, default: true },
      width: { type: String, default: '250px' }
    }
  }
</script>


<style lang="scss">
  @import './src/assets/styles/globals';

  .sidebar {
    padding-left: 0px !important;
    padding-right: 0px !important;
    position: fixed;
    top: 56px;
    bottom: 0;
    left: 0;
    z-index: 1000;
    padding: 20px;
    overflow-x: hidden;
    overflow-y: auto;
    border-right: 1px solid #eee;
  }

  #sidebar-toggle-button {
  padding-left: 0px !important;
    padding-right: 0px !important;
    position: fixed;
    top: 56px;
    bottom: 0;
    left: 0;
    z-index: 1000;

    border: 1px solid #eee;
  }

  .content {
    min-height: 85%;
  }
</style>
