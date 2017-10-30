<!-- Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights
     reserved. Use of this source code is governed by the AGPLv3 license that
     can be found in the LICENSE file. -->

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
    },

    created: function () {
      this.$store.commit('sidebarShown', this.shown)
    },

    destroyed: function () {
      this.$store.commit('sidebarShown', false)
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
    background-color: #f7f7f9;
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
