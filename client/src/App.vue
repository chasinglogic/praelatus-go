<!-- Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights
     reserved. Use of this source code is governed by the AGPLv3 license that
     can be found in the LICENSE file. -->

<template>
  <div id="app">
    <navbar></navbar>
    <div class="content" v-bind:style="{ marginLeft: showSidebar ? sidebarWidth : '0px' }">
      <router-view></router-view>
    </div>
  </div>
</template>

<script>
 import Navbar from '@/components/General/Navbar'

 import { mapState } from 'vuex'

 export default {
   name: 'app',
   computed: mapState({
     sidebarWidth: (state) => { return state.sidebarWidth },
     showSidebar: (state) => { return state.showSidebar }
   }),

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
   padding-top: 1rem;
 }
</style>
