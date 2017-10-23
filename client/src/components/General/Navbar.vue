<!-- Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights
     reserved. Use of this source code is governed by the AGPLv3 license that
     can be found in the LICENSE file. -->

<template>
  <b-navbar fixed=top toggle-breakpoint="md" toggleable class="bg-praelatus">
    <b-nav-toggle target="nav_collapse"></b-nav-toggle>

    <b-navbar-brand to="/">
      <img src="/assets/img/logo_arrow.svg"
        height="30"
        width="30">
    </b-navbar-brand>

    <b-collapse is-nav id="nav_collapse">

      <b-nav is-nav-bar>
        <b-nav-item to="/projects">Projects</b-nav-item>
        <b-nav-item-dropdown text="Tickets">
          <b-dropdown-item to="/queries/mine">Manage Queries</b-dropdown-item>
          <b-dropdown-item to="/queries">Search</b-dropdown-item>
        </b-nav-item-dropdown>
      </b-nav>
      <b-nav is-nav-bar class="ml-auto">
        <b-nav-item
          class="btn-success create-btn"
          to="/tickets/create">
          Create
        </b-nav-item>
        <b-nav-item-dropdown id="userMenu" right v-if="currentUser">
          <img slot="button-content" width="30" height="30" class="userMenuPic"
            v-bind:src="currentUser.profilePicture" alt="User">
          <b-dropdown-item v-if="currentUser.isAdmin" to="/admin/">
            System Administration
          </b-dropdown-item>
          <b-dropdown-item @click="logout">Log Out</b-dropdown-item>
        </b-nav-item-dropdown>
        <template v-else>
          <b-nav-item to="/login">Sign In</b-nav-item>
          <b-nav-item >/</b-nav-item>
          <b-nav-item to="/register">Sign Up</b-nav-item>
        </template>
      </b-nav>
    </b-collapse>
  </b-navbar>
</template>

<script>
 export default {
   name: 'navbar',
   methods: {
     logout: function (e) {
       e.preventDefault()
       this.$store.commit('logout', {})
       this.$router.push('/')
     }
   },
   computed: {
     currentUser: function () {
       return this.$store.getters.currentUser
     }
   }
 }
</script>


<style lang="scss">
 @import './src/assets/styles/globals.scss';

 .create-btn {
   border-radius: 5px;
   margin-right: 1rem;
 }

 .navbar.bg-praelatus {
   background-color: $primary-color;
   color: #fff;
 }

 .bg-praelatus .nav .nav-item a {
   color: #fff
 }

 .bg-praelatus .nav .nav-item a:hover {
   background-color: $highlight-color;
   color: #fff
 }

 .bg-praelatus .dropdown-item {
   color: #000 !important;
 }

 .bg-praelatus .dropdown-item:hover {
   background-color: $primary-color !important;
   color: #fff !important;
 }

 .bg-praelatus .dropdown-item.active {
   background-color: #fff;
 }

 .b-nav-dropdown .dropdown-item:focus:not(.active), .b-nav-dropdown .dropdown-item:hover:not(.active) {
   box-shadow: none !important;
 }


 #userMenu .dropdown-toggle {
   color: #fff;
   padding: 0;
 }

 #userMenu {
   padding: 5px;
 }
</style>
