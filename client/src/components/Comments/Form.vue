<template>
  <div>
    <div class="card comment" v-if="currentUser">
      <h1>Add a Comment</h1>
      <editor v-model="comment.body"></editor>
    </div>
    <div class="card comment" v-else>
      <div class="comment card-block">
        You must be <a href="/#/login">logged in</a> to comment.
      </div>
    </div>
  </div>
</template>

<script>
 import Editor from '@/components/General/Editor'

 export default {
   name: 'comment-form',

   components: {
     Editor
   },

   computed: {
     currentUser: function () {
       return this.$store.getters.currentUser
     }
   },

   data: function () {
     return {
       comment: {
         author: this.currentUser ? this.currentUser.username : '',
         body: ''
       }
     }
   },

   props: ['body'],

   mounted: function () {
     if (this.body && this.body !== '') {
       this.comment.body = this.body
     }
   }
 }
</script>
