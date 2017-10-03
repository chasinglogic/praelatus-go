<!-- Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights
     reserved. Use of this source code is governed by the AGPLv3 license that
     can be found in the LICENSE file. -->

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
