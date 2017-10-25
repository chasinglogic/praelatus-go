<template>
  <div>
    <div class="card comment" v-if="currentUser">
      <editor text="body" :submit="newComment"></editor>
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

   methods: {
     newComment: function (text) {
       let url = '/api/tickets/' + this.$route.params.key + '/addComment'
       Axios.post(url,
         {
           body: text,
           author: this.comment.author
         },
         {
           headers: {
             Authorization: 'Bearer ' + this.$store.getters.token
           },
           withCredentials: true
         }).then((res) => {
           console.log(res.data)
           if (this.reloadFunc) {
             this.reloadFunc()
           }
         }).catch((err) => {
           console.log(err)
         })
     }
   },

   data: function () {
     return {
       comment: {
         author: this.currentUser ? this.currentUser.username : ''
       }
     }
   },

   props: ['body', 'reloadFunc']
 }
</script>

<style lang="scss">
 .md-editor {
   width: 100%;
 }
</style>
