<template>
  <h1>Welcome to Praelatus!</h1>
</template>


<script>
 export default {
   name: 'index',
   created: function () {
     if (this.$store.getters.currentUser !== {}) {
       this.$route.router.go('/dashboard')
     }
   }
 }
</script>
