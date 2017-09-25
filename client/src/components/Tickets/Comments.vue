<template>
  <div class="comment-list">
    <template v-for="comment in comments">
      <comment :comment="comment" />
    </template>
  </div>
</template>

<script>
 import Comment from './Comment'

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
