<template>
  <div class="card comment">
    <div v-html="markdown(comment.body)" class="card-block">
    </div>

    <div class="card-header comment-header" >
      <p>
        <user-stub :username="comment.author" />
        commented on
        {{ dateFormat(comment.createdDate) }}
      </p>
    </div>
  </div>
</template>

<script>
 import UserStub from '@/components/Users/Stub'

 import dateUtils from '@/lib/dates'
 import Showdown from 'showdown'
 const converter = new Showdown.Converter()

 export default {
   name: 'comment',
   methods: {
     markdown (text) {
       return converter.makeHtml(text)
     },

     dateFormat: dateUtils.dateFormat
   },

   components: {
     'user-stub': UserStub
   },

   props: {
     comment: {
       name: 'comment',
       default: () => { return {} }
     }
   }
 }
</script>

<style>
 .comment {
   text-align: left;
   margin-top: 1rem;
 }

 .comment-header p {
   display: inline-block;
 }

 .comment-author-avatar {
   height: 20px;
   width: 20px;
   display: inline-block;
 }
</style>
