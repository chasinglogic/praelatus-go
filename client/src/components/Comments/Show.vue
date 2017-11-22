<!-- Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights
     reserved. Use of this source code is governed by the AGPLv3 license that
     can be found in the LICENSE file. -->

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

 .comment-header {
   border-top: 1px solid rgba(0, 0, 0, 0.125);
   border-bottom: none;
   height: 4rem;
 }
 .comment-author-avatar {
   height: 20px;
   width: 20px;
   display: inline-block;
 }
</style>
