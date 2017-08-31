<template>
  <b-navbar fixed=top toggle-breakpoint="md" toggleable class="bg-praelatus">
    <b-nav-toggle target="nav_collapse"></b-nav-toggle>

    <b-navbar-brand href="/">
      <img src="/assets/static/img/logo_arrow.svg"
           height="30"
           width="30">
    </b-navbar-brand>

    <b-collapse is-nav id="nav_collapse">

      <b-nav is-nav-bar>
        <b-nav-item href="/projects/search">Projects</b-nav-item>
        <b-nav-item-dropdown text="Tickets">
          <b-dropdown-item to="/queries/mine">Manage Queries</b-dropdown-item>
          <b-dropdown-item to="/queries">Search</b-dropdown-item>
        </b-nav-item-dropdown>
      </b-nav>
      <b-nav is-nav-bar class="ml-auto">
        <b-nav-item class="nav-btn btn-success"
          href="/tickets/create">Create</b-nav-item>
        <b-nav-item-dropdown v-if="currentUser">
          <img width="30" height="30" src="" alt="User">
          <b-dropdown-item to="/admin/">System Administration</b-dropdown-item>
          <b-dropdown-item to="/logout">Log Out</b-dropdown-item>
        </b-nav-item-dropdown>
        <template v-else>
          <b-nav-item href="/login">Sign In</b-nav-item>
          <b-nav-item >/</b-nav-item>
          <b-nav-item href="/signup">Sign Up</b-nav-item>
        </template>
      </b-nav>
    </b-collapse>
  </b-navbar>
</template>

<script>
export default {
  name: 'navbar',
  computed: {
    currentUser: function () {
      return this.$store.getters.currentUser
    }
  }
}
</script>


<style>
.btn-nav {
  margin: 0.1rem;
  border-radius: 5px;
}
</style>
