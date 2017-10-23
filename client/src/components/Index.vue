<!-- Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights
     reserved. Use of this source code is governed by the AGPLv3 license that
     can be found in the LICENSE file. -->

<template>
  <b-container>
    <b-row>
      <b-col>
        <h1>Welcome to Praelatus!</h1>
        <p class="text-big">
          Praelatus is an Open Source bug tracking and ticketing system.
          If you're new to Praelatus and are wondering how to get started
          you can view our <a href="https://doc.praelatus.io/">User Guide</a>.
          Otherwise just sign in or sign up and get cracking!
        </p>
      </b-col>
      <b-col>
        <b-card no-body>
          <b-tabs ref="tabs" card>
            <b-tab title="Sign In" active>
              <login></login>
            </b-tab>
            <b-tab title="Sign Up">
              <login register="true"></login>
            </b-tab>
          </b-tabs>
        </b-card>
      </b-col>
    </b-row>
  </b-container>

</template>


<script>
 import Login from '@/components/Users/Login'

 export default {
   name: 'index',

   components: {
     Login
   },

   created: function () {
     if (this.$store.getters.currentUser !== {} &&
         this.$store.getters.currentUser !== null) {
       this.$router.push('/dashboard')
     }
   }
 }
</script>


<style>
 .text-big {
   font-size: 1.5rem;
 }
</style>
