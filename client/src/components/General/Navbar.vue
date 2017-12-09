<!-- Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights
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
        <b-nav-form>
          <b-form-input size="sm" class="mr-sm-2 navbar-ticket-search"
            type="text" placeholder="Search Tickets"/>
        </b-nav-form>

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

 .bg-praelatus .nav-link {
   color: #fff !important;
 }

 .bg-praelatus .nav-link:hover {
   background-color: lighten($primary-color, 10%);
 }

 .bg-praelatus .btn-success .nav-link:hover, .bg-praelatus .btn-success .nav-link.active:hover {
   background-color: lighten(#28a745, 10%);
   border-radius: 5px;
 }

 .bg-praelatus .btn-success .nav-link.active {
   background-color: #28a745;
   border-radius: 5px;
 }

 .bg-praelatus .nav-link.btn-success:hover {
 }

 .bg-praelatus .dropdown-item {
   color: #000 !important;
 }

 .bg-praelatus .dropdown-item:hover, .bg-praelatus .dropdown-item.active:hover, .bg-praelatus .nav-link.active {
   background-color: $primary-color;
   color: #fff !important;
 }

 .bg-praelatus .dropdown-item.active {
   background-color: white;
   border: 1px solid $primary-color;
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

 .bg-praelatus .nav .form-inline input {
   border: none;
   border-radius: 0px;
   background-color: lighten($highlight-color, 20%);
   color: white;
 }

 .bg-praelatus .nav .form-inline input::-webkit-input-placeholder { /* WebKit, Blink, Edge */
   color:    #fff !important;
 }
 .bg-praelatus .nav .form-inline input:-moz-placeholder { /* Mozilla Firefox 4 to 18 */
   color:    #fff !important;
   opacity:  1;
 }
 .bg-praelatus .nav .form-inline input::-moz-placeholder { /* Mozilla Firefox 19+ */
   color:    #fff !important;
   opacity:  1;
 }
 .bg-praelatus .nav .form-inline input:-ms-input-placeholder { /* Internet Explorer 10-11 */
   color:    #fff !important;
 }
 .bg-praelatus .nav .form-inline input::-ms-input-placeholder { /* Microsoft Edge */
   color:    #fff !important;
 }

 .navbar-ticket-search {
   /* TODO: Make this responsive */
   width: 100px;
    -webkit-transition: width 1s ease-in-out;
    -moz-transition:width 1s ease-in-out;
    -o-transition: width 1s ease-in-out;
    transition: width 1s ease-in-out;
 }

 .navbar-ticket-search:focus {
   width: 300px;
 }

 .navbar-nav .dropdown-menu .dropdown-item {
   border: none !important;
   box-shadow: none !important;
 }
</style>
