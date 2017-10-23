<!-- Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights
     reserved. Use of this source code is governed by the AGPLv3 license that
     can be found in the LICENSE file. -->

<template>
  <span>
    <router-link v-bind:to="'/users/' + username">
      <img class="avatar" v-bind:src="avatar" />
      {{ username }}
    </router-link>
  </span>
</template>


<script>
 import Axios from 'axios'

 export default {
   name: 'user-stub',

   data: () => {
     return { avatar: '' }
   },

   props: {
     username: {
       name: 'username',
       default: ''
     }
   },

   methods: {
     getAvatar: function (username) {
       let url = '/api/users/' + this.username + '/avatar'
       let instance = this

       Axios.get(url)
            .then(function (res) {
              instance.avatar = res.data
            })
            .catch(function (err) {
              if (err) {
                console.log(err)
              }

              instance.avatar = 'https://www.gravatar.com/avatar/00000000000000000000000000000000'
            })
     }
   },

   created: function () {
     this.getAvatar(this.username)
   }
 }
</script>
