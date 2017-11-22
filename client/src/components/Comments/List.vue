<!-- Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights
     reserved. Use of this source code is governed by the AGPLv3 license that
     can be found in the LICENSE file. -->

<template>
  <div class="comment-list">
    <template v-for="comment in comments">
      <comment :comment="comment" />
    </template>
  </div>
</template>

<script>
 import Comment from './Show'

 export default {
   name: 'comments',
   components: {
     Comment
   },
   props: {
     comments: {
       name: 'comments',
       default: () => []
     }
   }
 }
</script>
