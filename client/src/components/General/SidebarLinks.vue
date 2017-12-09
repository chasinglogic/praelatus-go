<!-- Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
     Use of this source code is governed by the AGPLv3 license that can be found in
     the LICENSE file. -->

<template>
  <div class="sidebar-links">
    <slot></slot>
  </div>
</template>

<script>
 export default {
   name: 'sidebar-links'
 }
</script>

<style lang="scss">
  @import './src/assets/styles/globals';

 .sidebar-links a:hover {
   background-color: $primary-color;
   color: white;
 }

 .sidebar-links a {
   color: $primary-color;
   text-decoration: none;
   text-align: center;
   padding: 0.5rem;
   display: block;
   width: 100%;
 }
</style>
