<template>
  <div>
    <div class="card comment" v-if="currentUser">
      <editor v-model="body" @submit="newComment"></editor>
      <b-button @click="newComment" variant="primary">
        Add Comment
      </b-button>
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
 import Axios from 'axios'

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
     return { body: '' }
   },

   methods: {
     newComment: function () {
       let url = '/api/tickets/' + this.$route.params.key + '/addComment'
       Axios.post(url,
         {
           body: this.body,
           author: this.currentUser.username
         }).then((res) => {
           this.$emit('newComment')
           this.body = ''
         }).catch((err) => {
           console.log('ERROR:', err)
         })
     }
   }
 }
</script>

<style lang="scss">
 .md-editor {
   width: 100%;
 }
</style>
