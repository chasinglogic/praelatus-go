<template>
  <a v-bind:href="'/#/user/' + username">
    <img class="avatar" v-bind:src="avatar" />
    {{ username }}
  </a>
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
